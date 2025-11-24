package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gotomicro/ego/core/econf"
	sdkapi "github.com/shimo-open/sdk-kit-go/model/api"
	"golang.org/x/crypto/bcrypt"

	"sdk-demo-go/pkg/consts"
)

// CustomClaims extends JWT standard claims with additional fields
type CustomClaims struct {
	*jwt.StandardClaims
	Scope   string `json:"scope,omitempty"`
	Version string `json:"version,omitempty"`
}

// HashPassword hashes a user password
func HashPassword(passwrod string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(passwrod), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashed)
}

// Sign issues a Shimo signature
// When strict is true, the signature lasts 4 minutes and includes scope:"license"
// When strict is false, the signature lasts 1 year
func Sign(appID, appSecret string, strict bool) string {
	nowTime := time.Now()

	var addTime time.Duration
	if strict {
		addTime = time.Minute * 4
	} else {
		addTime = time.Hour * 24 * 365
	}
	exp := nowTime.Add(addTime).Unix()

	if strict {
		return SignJWT(appID, appSecret, exp, true)
	} else {
		return SignJWT(appID, appSecret, exp, false)
	}

}

// SignJWT generates a JWT token with the given parameters
func SignJWT(kid, secret string, expires int64, withScope bool) string {
	var token *jwt.Token
	if withScope {
		token = jwt.NewWithClaims(
			jwt.SigningMethodHS256, &CustomClaims{
				StandardClaims: &jwt.StandardClaims{
					ExpiresAt: expires,
				},
				Scope:   "license",
				Version: econf.GetString("shimoSDK.callbackVersion"),
			},
		)
	} else {
		token = jwt.NewWithClaims(
			jwt.SigningMethodHS256, &CustomClaims{
				StandardClaims: &jwt.StandardClaims{
					ExpiresAt: expires,
				},
				Version: econf.GetString("shimoSDK.callbackVersion"),
			},
		)
	}

	token.Header["kid"] = kid

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}

	return tokenStr
}

// UserClaims represents JWT claims for user authentication
type UserClaims struct {
	*jwt.StandardClaims
	UserId int64  `json:"userId"`
	Mode   string `json:"mode"`
}

// SignUserJWT issues a user token
func SignUserJWT(userId int64, expr ...time.Duration) string {
	var expires time.Duration
	if len(expr) > 0 {
		expires = expr[0]
	} else {
		expires = 24 * time.Hour
	}
	secret := econf.GetString("jwt.secret")
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, &UserClaims{
			StandardClaims: &jwt.StandardClaims{
				ExpiresAt: time.Now().Add(expires).Unix(),
			},
			UserId: userId,
		})
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}

	return tokenStr
}

// SignUserJWTWithMode generates a user JWT token with a specific mode
func SignUserJWTWithMode(userId int64, mode string) string {
	secret := econf.GetString("jwt.secret")
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, &UserClaims{
			StandardClaims: &jwt.StandardClaims{
				ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			},
			UserId: userId,
			Mode:   mode,
		})
	token.Header["kid"] = econf.GetString("shimoSDK.appId")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}

	return tokenStr
}

// SDKClaims represents JWT claims for SDK operations
type SDKClaims struct {
	*jwt.StandardClaims
	FileId string `json:"fileId"`
	UserId string `json:"userId"`
}

// GenFileGuid creates a 16-character file GUID
func GenFileGuid() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)[:16]
}

// GenFileName generates a filename with timestamp and file type
func GenFileName(fileType consts.FileType) string {
	return time.Now().Format("2006-01-02 15:04:05") + "-" + fileType.String()
}

// GetAuth generates authentication credentials for a user
func GetAuth(userId int64, isStrict ...bool) (auth sdkapi.Auth) {
	auth.Token = SignUserJWT(userId)
	strict := false
	if len(isStrict) > 0 {
		strict = isStrict[0] // Use the first provided value if present
	}
	auth.Signature = Sign(econf.GetString("shimoSDK.appId"), econf.GetString("shimoSDK.appSecret"), strict)
	if econf.GetString("shimoSDK.callbackVersion") == "v2" {
		auth.UserUuid = GetHashUserUuid(userId)
	}
	return
}

// GetHashUserUuid generates a SHA256 hash of the user ID
func GetHashUserUuid(userId int64) string {
	hasher := sha256.New()
	hasher.Write([]byte(strconv.FormatInt(userId, 10)))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GenerateUserFileUUID generates a unique file UUID based on user ID and file type
func GenerateUserFileUUID(userId string, fileType string) string {
	// Concatenate userId and fileType
	input := userId + "_" + fileType
	// Use MD5 to produce a fixed 16-byte UUID from the combined string
	hash := md5.New()
	hash.Write([]byte(input))
	// Take the first 16 bytes of the hash as the UUID
	return hex.EncodeToString(hash.Sum(nil))[:16]
}

// EncodeCharset follows the JavaScript string order, see
// https://github.com/felipecarrillo100/base62str/blob/40c9acae36/src/index.ts#L45
const EncodeCharset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

const (
	StandardBase int32 = 256
	TargetBase   int32 = 62
)

var GMP = make([]rune, 256, 256)
var lookup = make([]rune, 256, 256)

func init() {
	for idx, s := range EncodeCharset {
		GMP[idx] = s
	}

	var i int32 = 0
	for ; i < 256; i++ {
		lookup[GMP[i]] = i & 0x00FF
	}
}

func getBytes(input string) []rune {
	rs := make([]rune, 0)

	for _, r := range []rune(input) {
		rs = append(rs, r&0x00FF)
	}

	return rs
}

func encode(input []rune) []rune {
	indices := convert(input, StandardBase, TargetBase)
	return translate(indices, GMP)
}

func decode(input []rune) []rune {
	prepared := translate(input, lookup)

	return convert(prepared, TargetBase, StandardBase)
}

// convert see
// https://github.com/felipecarrillo100/base62str/blob/master/src/index.ts#L88
func convert(input []rune, sourceBase int32, targetBase int32) []rune {
	out := make([]int32, 0)
	source := input

	for {
		if len(source) == 0 {
			break
		}

		quotient := make([]int32, 0)
		var remainder int32

		for i, source1 := 0, source; i < len(source1); i++ {
			sourcei := source1[i]
			accumulator := (sourcei & 0x00FF) + remainder*sourceBase
			digit := (accumulator - (accumulator % targetBase)) / targetBase
			remainder = accumulator % targetBase

			if len(quotient) > 0 || digit > 0 {
				quotient = append(quotient, digit)
			}
		}

		out = append(out, remainder)
		source = quotient
	}

	for i := 0; i < len(input)-1 && input[i] == 0; i++ {
		out = append(out, 0)
	}

	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}

	return out
}

func translate(indices []rune, dict []rune) []rune {
	translation := make([]rune, 0)

	for i, indices1 := 0, indices; i < len(indices1); i++ {
		var indicesi = indices1[i]
		translation = append(translation, dict[indicesi])
	}

	return translation
}

func getString(input []rune) string {
	return string(input)
}

// Base62Encode encodes a string to base62 format
func Base62Encode(input string) (output string) {
	if len(input) == 0 {
		return ""
	}

	return getString(encode(getBytes(input)))
}

// Base62Decode decodes a base62 encoded string
func Base62Decode(input string) (output string) {
	if len(input) == 0 {
		return ""
	}

	return getString(decode(getBytes(input)))
}

// ConvertToOSFile converts a multipart.FileHeader to an *os.File by saving it to a temporary location
func ConvertToOSFile(fileHeader *multipart.FileHeader) (*os.File, error) {
	// Open the uploaded file stream
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open multipart file: %w", err)
	}
	defer src.Close()
	tmpDir := os.TempDir()                                    // Determine the temp directory
	tmpFilePath := filepath.Join(tmpDir, fileHeader.Filename) // Destination path

	// Create a file with the fixed name directly
	tmpFile, err := os.Create(tmpFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()
	// Copy the content into the temp file
	buf := make([]byte, 32*1024) // 32KB buffer
	if _, err := io.CopyBuffer(tmpFile, src, buf); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}
	// Reopen the file to obtain an *os.File handle
	osFile, err := os.Open(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to open temp file as os.File: %w", err)
	}
	return osFile, nil
}

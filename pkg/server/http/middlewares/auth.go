package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	sdkapi "github.com/shimo-open/sdk-kit-go/api"

	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/models/db"
	"sdk-demo-go/pkg/utils"
)

func FindAccessToken(c *gin.Context) string {
	str := c.GetHeader("Authorization")
	if strings.HasPrefix(str, "bearer ") {
		return str[7:]
	} else {
		str = c.Query("accessToken")
		return str
	}
}

type StartParam struct {
	Path     string `json:"path"`
	Type     string `json:"type"`
	FileGuid string `json:"fileGuid"`
}

func UserAuthMiddleware(c *gin.Context) {
	c.Set("multipleClientMode", econf.GetBool("shimoSDK.multipleClientMode"))
	fullPath := c.FullPath()
	if fullPath == "/api/users/signin" || fullPath == "/api/users/signup" || fullPath == "/api/users/auth" {
		setAppClient(c, "") // todo
		c.Next()
		return
	}
	// Handle form requests separately
	if fullPath == "/api/files/:fileGuid" {
		mode, _ := c.GetQuery("mode")
		smParam, _ := c.GetQuery("smParams")
		var transSmParam StartParam
		if smParam != "" {
			if err := json.Unmarshal([]byte(utils.Base62Encode(smParam)), &transSmParam); err == nil {
				elog.Info(fmt.Sprintf("smParam: %+v", smParam))
			}
		}
		if (mode == "form_fill") || (strings.Contains(transSmParam.Path, "response-share")) {
			c.Next()
			return
		}
	}

	token := FindAccessToken(c)

	// Skip token requirement for this path
	if strings.HasPrefix(fullPath, "/api/files/importUrl") {
		users, err := db.FindUsersByAppId(invoker.DB, econf.GetString("shimoSDK.appId"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "user not found",
			})
		}
		user := users[0]
		token = utils.SignUserJWT(user.ID)
	}

	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "token is required",
		})
		return
	}

	err := ValidateUserToken(c, token)
	if err != nil {
		return
	}

	_userId, _ := c.Get("userId")
	userId, _ := _userId.(int64)

	user, err := db.FindUserById(invoker.DB, userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "user does not exist",
		})
		return
	}

	setAppClient(c, user.AppID)

	c.Next()
}

func setAppClient(c *gin.Context, appId string) {
	if c.GetBool("multipleClientMode") {
		ac, _ := db.AppClientFindById(invoker.DB, appId)
		c.Set("appClient", ac)
	} else {
		ac := db.AppClient{
			AppID:     econf.GetString("shimoSDK.appId"),
			AppSecret: econf.GetString("shimoSDK.appSecret"),
		}
		c.Set("appClient", ac)
	}
}

// ValidateUserToken verifies the token and stores userId in the context
func ValidateUserToken(c *gin.Context, token string) error {
	if token == sdkapi.AnonymousToken {
		// Anonymous mode: form_fill with userId -1
		c.Set("userId", sdkapi.Anonymous)
		c.Set("mode", "form_fill")
		return nil
	} else {
		decodedToken, err := jwt.ParseWithClaims(token, &utils.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(econf.GetString("jwt.secret")), nil
		})

		if decodedToken == nil || !decodedToken.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
			})
			return err
		}

		claim, ok := decodedToken.Claims.(*utils.UserClaims)
		if !ok {
			panic("parse token error")
		}
		c.Set("userId", claim.UserId)
		c.Set("mode", claim.Mode)
		return err
	}
}

func FrontInspectAuthMiddleware(c *gin.Context) {
	var user *db.User
	token := c.GetHeader("Authorization")
	if token == econf.GetString("frontInspect.internalToken") {
		// Look up the user
		email := econf.GetString("frontInspect.inspectEmail")
		appId := econf.GetString("shimoSDK.appId")
		defaultPassword := econf.GetString("frontInspect.defaultPassword")
		var err error
		user, err = db.FindUserByInstance(invoker.DB, &db.User{
			Email: email,
			AppID: appId,
		})
		if err != nil {
			// Create the configured user if not found
			user = &db.User{
				Email:    email,
				Password: utils.HashPassword(defaultPassword),
				Avatar:   fmt.Sprintf("%sstatic/img/default-avatar-moke.png", econf.GetString("publicPath.publicPath")),
				AppID:    appId,
				Name:     strings.Split(email, "@")[0],
			}
			err = db.CreateUser(invoker.DB, user)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "inspect user create err: " + err.Error(),
				})
				return
			}
		}
	} else {
		// Standard token handling path
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "token is required",
			})
			return
		}
		if strings.HasPrefix(token, "bearer ") {
			token = token[7:]
		}
		err := ValidateUserToken(c, token)
		if err != nil {
			return
		}
		_userId, _ := c.Get("userId")
		userId, _ := _userId.(int64)
		user, err = db.FindUserById(invoker.DB, userId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "user does not exist",
			})
			return
		}
	}
	c.Set("userId", user.ID)
	setAppClient(c, user.AppID)

	c.Next()
}

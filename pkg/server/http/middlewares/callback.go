package middlewares

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	sdkapi "github.com/shimo-open/sdk-kit-go/api"

	"sdk-demo-go/pkg/utils"
)

func findShimoToken(c *gin.Context) string {
	str := c.GetHeader(sdkapi.HeaderShimoToken)
	return str
}

func CallbackAuthMiddleware(c *gin.Context) {
	switch c.GetHeader(sdkapi.HeaderShimoCredentialType) {
	case "0":
		token := findShimoToken(c)
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
		c.Next()
		return
	case "3":
		signature := c.GetHeader(sdkapi.HeaderShimoSignature)
		if signature == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "signature is required",
			})
			return
		}
		checkSignature(c, signature)
		c.Next()
		return
	}
}

// checkSignature verifies the signature and ensures parameters match
func checkSignature(c *gin.Context, signature string) {
	decodedToken, err := jwt.ParseWithClaims(signature, &utils.SDKClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(econf.GetString("shimoSDK.appSecret")), nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "parse signature error:" + err.Error(),
		})
		return
	}

	if !decodedToken.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid signature",
		})
		return
	}

	claims, ok := decodedToken.Claims.(*utils.SDKClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "parse signature error",
		})
		return
	}

	elog.Info("checkSignature", l.A("claims", claims))

	c.Set("claims", claims)
	checkParams(c, claims)
}

// checkParams ensures the request parameters align with the claims
func checkParams(c *gin.Context, claims *utils.SDKClaims) {
	fileGuid := c.Param("fileGuid")

	if fileGuid != "" && fileGuid != claims.FileId {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "fileId does not match",
		})
		return
	}

	userId, exist := c.GetQuery("userId")
	if exist && userId != claims.UserId {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "userId does not match",
		})
		return
	}
}

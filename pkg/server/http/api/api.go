package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sdk-demo-go/pkg/invoker"
)

// SignJWTReq represents the request parameters for JWT signing
type SignJWTReq struct {
	AppID     string `json:"appId" form:"appId" required:"true"`
	AppSecret string `json:"appSecret" form:"appSecret" required:"true"`
	Strict    bool   `json:"strict" form:"strict"`
}

// SignJWT generates a JWT signature for the given credentials
func SignJWT(c *gin.Context) {
	params := SignJWTReq{}
	if err := c.Bind(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	signature := invoker.Services.SignatureService.Sign(params.AppID, params.AppSecret, params.Strict)
	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"signature": signature,
	})
	return
}

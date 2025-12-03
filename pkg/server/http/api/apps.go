package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/econf"
	sdkapi "github.com/shimo-open/sdk-kit-go/api"

	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/utils"
)

// GetAppDetails retrieves application details
func GetAppDetails(c *gin.Context) {
	appId := econf.GetString("shimoSDK.appId")
	auth := utils.GetAuth(getUserIdFromToken(c))
	params := sdkapi.GetAppDetailReq{
		Metadata: auth,
		AppID:    appId,
	}

	details, err := invoker.SdkMgr.GetAppDetail(c.Request.Context(), params)
	if err != nil {
		handleSdkMgrError(c, details.Response().Body(), details.Response().StatusCode())
		return
	}

	c.JSON(200, details)
}

func PutEndpointUrl(c *gin.Context) {
	appId := econf.GetString("shimoSDK.appId")
	url := econf.GetString("shimoSDK.endpoint")
	if url == "" {
		url = "http://svc-sdk2-demo:9001/callback"
	}
	auth := utils.GetAuth(getUserIdFromToken(c))
	params := sdkapi.UpdateCallbackURLReq{
		Metadata:                 auth,
		UpdateCallbackURLReqBody: sdkapi.UpdateCallbackURLReqBody{URL: url},
		AppID:                    appId,
	}

	resp, err := invoker.SdkMgr.UpdateCallbackURL(c.Request.Context(), params)
	if err != nil {
		handleSdkMgrError(c, resp.Response().Body(), resp.Response().StatusCode())
		return
	}

	c.JSON(204, nil)
}

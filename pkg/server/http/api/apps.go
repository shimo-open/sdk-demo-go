package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/econf"
	sdkapi "github.com/shimo-open/sdk-kit-go/model/api"

	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/utils"
)

// GetAppDetails retrieves application details
func GetAppDetails(c *gin.Context) {
	appId := econf.GetString("shimoSDK.appId")
	auth := utils.GetAuth(getUserIdFromToken(c), true)
	params := sdkapi.GetAppDetailParams{
		Auth:  auth,
		AppId: appId,
	}

	details, err := invoker.SdkMgr.GetAppDetail(params)
	if err != nil {
		handleSdkMgrError(c, details.Resp.Body(), details.Resp.StatusCode())
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
	auth := utils.GetAuth(getUserIdFromToken(c), true)
	params := sdkapi.UpdateCallbackUrlParams{
		Auth:                     auth,
		UpdateCallbackUrlReqBody: sdkapi.UpdateCallbackUrlReqBody{Url: url},
		AppId:                    appId,
	}

	resp, err := invoker.SdkMgr.UpdateCallbackUrl(params)

	if err != nil {
		handleSdkMgrError(c, resp.Resp.Body(), resp.Resp.StatusCode())
		return
	}

	c.JSON(204, nil)
}

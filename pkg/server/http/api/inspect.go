package api

import (
	"fmt"
	"net/http"

	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/services/inspect"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/elog"
)

func GetWebInspect(c *gin.Context) {
	htmlContent, err := inspect.GetWebInspect(c, invoker.Services.InspectHttp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("GetWebInspect fail, %w", err).Error()})
		return
	}

	elog.Info(fmt.Sprintf("htmlContent: %s", htmlContent))

	// Return the rendered HTML content
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
	return
}

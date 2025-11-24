package http

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egin"

	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/server/http/middlewares"
	"sdk-demo-go/ui"
)

// ServeHTTP initializes and returns the HTTP server component
func ServeHTTP() *egin.Component {
	r := invoker.Gin
	r.Use(LogRequestURL())

	registerDemoAppAPIs(r)
	registerCallbackAPIs(r)

	r.Use(middlewares.Serve("/", middlewares.EmbedFolder(ui.WebUI, "dist"), false))
	r.Use(middlewares.Serve("/", middlewares.FallbackFileSystem(middlewares.EmbedFolder(ui.WebUI, "dist")), true))

	return r
}

// ParseAppUrlAndSubUrl parses the application URL and extracts the base URL and sub-path
func ParseAppUrlAndSubUrl(appUrl string) (string, string, error) {
	if appUrl == "" {
		appUrl = "http://localhost:9001/"
	}
	if appUrl[len(appUrl)-1] != '/' {
		appUrl += "/"
	}
	// Check whether contains subpaths;
	urlParsed, err := url.Parse(appUrl)
	if err != nil {
		elog.Error("Invalid root_url.", l.S("url", appUrl), l.S("error", err.Error()))
		os.Exit(1)
	}
	appSubUrl := strings.TrimSuffix(urlParsed.Path, "/")
	return appUrl, appSubUrl, nil
}

// LogRequestURL returns a middleware that logs incoming API and callback requests
func LogRequestURL() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log the HTTP method and URL for API and callback requests
		if strings.HasPrefix(c.Request.URL.Path, "/api") || strings.HasPrefix(c.Request.URL.Path, "/callback") {
			elog.Info(fmt.Sprintf("Received %s request for %s\n", c.Request.Method, c.Request.URL.Path))
		}
		// Continue with the remaining handlers
		c.Next()
	}
}

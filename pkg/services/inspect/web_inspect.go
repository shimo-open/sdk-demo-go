package inspect

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/cetus/xpongo"
	"github.com/gotomicro/ego/client/ehttp"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/mitchellh/mapstructure"
)

func GetWebInspect(c *gin.Context, httpClient *ehttp.Component) (htmlContent string, err error) {
	data, err := GetWebInspectData(httpClient)
	if err != nil {
		return "", err
	}

	templatePath, err := GetHtmlTemplate(httpClient, data)
	if err != nil {
		return "", err
	}

	htmlContext, err := GetHtmlContext(data)
	if err != nil {
		return "", err
	}

	elog.Info(fmt.Sprintf("GetWebInspect, templatePath: %s, htmlContext: %v", templatePath, htmlContext))

	// Dynamically load and render the template
	htmlContent, err = renderTemplate(templatePath, htmlContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("template rendering failed, %w", err).Error()})
		return "", err
	}

	return htmlContent, nil
}

// renderTemplate loads the template file dynamically and renders it
func renderTemplate(templatePath string, context map[string]interface{}) (string, error) {
	// Read the template content from the temp file
	tmplContent, err := os.ReadFile(templatePath)
	if err != nil {
		elog.Error("Failed to read template file", l.E(err))
		return "", fmt.Errorf("failed to read template file: %w", err)
	}

	// Use raymond to parse and render the Handlebars template
	rendered, err := raymond.Render(string(tmplContent), context)
	if err != nil {
		elog.Error("Template rendering failed", l.E(err))
		return "", fmt.Errorf("template rendering failed: %w", err)
	}

	return rendered, nil
}

func GetWebInspectData(httpClient *ehttp.Component) (map[string]interface{}, error) {
	baseHost := econf.GetString("frontInspect.http.addr")
	baseUrl := "/__lizard__/deploys/lizard-service-health"
	resp, err := httpClient.R().Get(baseHost + baseUrl)
	if err != nil {
		elog.Error("GetWebInspect failed to get lizard-service-health", l.E(err))
		return nil, err
	} else if resp.StatusCode() != http.StatusOK {
		elog.Error("GetWebInspect request error", l.I("respStatus", resp.StatusCode()), l.S("respRequestID", resp.Header().Get("x-request-id")))
		return nil, errors.New("GetWebInspect request error status not OK")
	}

	result := map[string]interface{}{}
	bodyBytes := resp.Body()
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		elog.Error("failed to unmarshal lizard-service-health", l.S("respBody", string(bodyBytes)), l.E(err))
		return nil, err
	}

	return result, nil
}

func GetHtmlTemplate(httpClient *ehttp.Component, data map[string]interface{}) (string, error) {
	type Hbs struct {
		IndexHbs struct {
			FileUrl string `json:"file"`
		} `json:"index.hbs"`
	}

	// Validate and parse the templates field inside data
	templates, ok := data["templates"]
	if !ok {
		return "", fmt.Errorf("GetFileUrl templates field does not exist")
	}

	jsonBytes, err := json.Marshal(templates)
	if err != nil {
		elog.Error("GetFileUrl JSON serialization error", l.E(err))
		return "", err
	}

	var hbs Hbs
	if err := json.Unmarshal(jsonBytes, &hbs); err != nil {
		elog.Error("GetFileUrl failed to parse JSON", l.S("jsonBytes", string(jsonBytes)), l.E(err))
		return "", errors.New("GetFileUrl failed to unmarshal templates")
	}

	// Retrieve the CDN host and template file URL
	cdnHost := econf.GetString("frontInspect.cdn")
	fileUrl := hbs.IndexHbs.FileUrl
	if fileUrl == "" {
		return "", errors.New("file path is empty")
	}

	resp, err := httpClient.R().
		SetDoNotParseResponse(true).
		Get(cdnHost + fileUrl)
	if err != nil {
		elog.Error("HTTP request failed", l.E(err))
		return "", err
	}
	defer resp.RawResponse.Body.Close() // Ensure the response body is closed

	if resp.StatusCode() != http.StatusOK {
		elog.Error("GetHtmlTemplate HTTP status code error", l.I("respStatus", resp.StatusCode()), l.S("respRequestID", resp.Header().Get("x-request-id")))
		return "", fmt.Errorf("HTTP status code error: %d", resp.StatusCode())
	}

	// Write the response body to a temporary file
	tempFile, err := os.CreateTemp("", "*.html")
	if err != nil {
		elog.Error("Failed to create temp file", l.E(err))
		return "", err
	}
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, resp.RawResponse.Body); err != nil {
		elog.Error("Failed to write to temp file",
			l.E(err),
			l.S("url", cdnHost+fileUrl),
			l.S("resp", fmt.Sprintf("%v", resp)),
			l.S("status", fmt.Sprintf("%v", resp.StatusCode())),
			l.S("headers", fmt.Sprintf("%v", resp.Header())),
		)
		return "", err
	}

	elog.Info("Successfully wrote HTML template to temp file", l.S("filePath", tempFile.Name()))
	return tempFile.Name(), nil
}

type Urls struct {
	Url string `json:"url"`
}

type Script struct {
	Main []Urls `json:"main"`
}

func GetHtmlContext(data map[string]interface{}) (xpongo.Context, error) {
	type Service struct {
		Scripts []string `json:"scripts"`
		Styles  []string `json:"styles"`
	}

	// Validate and parse the service field in data
	serviceList, ok := data["service"]
	if !ok {
		return nil, fmt.Errorf("GetHtmlContext service field does not exist")
	}

	// Convert serviceList into a []Service
	service, ok := serviceList.([]interface{})
	if !ok {
		return nil, fmt.Errorf("service field type error")
	}

	var services []Service
	for _, s := range service {
		svcMap, ok := s.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("service element type error")
		}
		var svc Service
		if err := mapstructure.Decode(svcMap, &svc); err != nil {
			return nil, fmt.Errorf("failed to parse service data: %v", err)
		}
		services = append(services, svc)
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("service array is empty")
	}

	pageContext := make(xpongo.Context)
	pageContext["title"] = "Inspection title"

	var mains []Urls
	for _, script := range AppendCdnHost(services[0].Scripts) {
		mains = append(mains, Urls{Url: script})
	}

	pageContext["scripts"] = Script{
		Main: mains,
	}

	var styles []Urls
	for _, style := range AppendCdnHost(services[0].Styles) {
		styles = append(styles, Urls{Url: style})
	}
	pageContext["styles"] = styles

	return pageContext, nil
}

func AppendCdnHost(data []string) []string {
	cdnHost := econf.GetString("frontInspect.cdn")
	var result []string
	for _, url := range data {
		result = append(result, cdnHost+url)
	}
	return result
}

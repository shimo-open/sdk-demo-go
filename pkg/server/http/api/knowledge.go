package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/econf"
	sdkapi "github.com/shimo-open/sdk-kit-go/api"

	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/models/db"
	"sdk-demo-go/pkg/utils"
)

// AIAssetsResponse models the AI asset response
type AIAssetsResponse struct {
	RuntimeEnv interface{} `json:"runtimeEnv"`
	Modules    []Module    `json:"modules"`
}

// Module represents an asset module
type Module struct {
	Css     []string `json:"css"`
	Js      []string `json:"js"`
	Name    string   `json:"name"`
	Devices []string `json:"devices"`
}

// GetAiAssets fetches AI static assets
func GetAiAssets(c *gin.Context) {
	// Adjust the API path
	url := econf.GetString("shimoSDK.host") + sdkapi.ApiAiAssets
	body, err := NewHttpRequest(c, url, "GET")
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"message": err.Error(),
			"error":   err.Error(),
		})
	}
	// Parse the JSON response
	var response AIAssetsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"message": "parse response json failed",
			"error":   err.Error(),
		})
		return
	}
	// Build the response payload
	result := make(map[string]interface{})

	// Include runtimeEnv
	if response.RuntimeEnv != nil {
		result["runtimeEnv"] = response.RuntimeEnv
	}

	// Find the module whose name is "ai"
	for _, module := range response.Modules {
		if module.Name == "ai" {
			result["aiModule"] = module
			break
		}
	}
	assetsUrl := econf.GetString("shimoSDK.host") + sdkapi.ApiIframeAssets
	assetsBody, err := NewHttpRequest(c, assetsUrl, "GET")
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"message": err.Error(),
			"error":   err.Error(),
		})
	}
	// Parse the JSON response
	var assetsResponse map[string]Module
	if err := json.Unmarshal(assetsBody, &assetsResponse); err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"message": "parse assetsResponse json failed",
			"error":   err.Error(),
		})
		return
	}

	for key, module := range assetsResponse {
		if key == "ai" {
			result["aiAssetsModule"] = module
			break
		}
	}

	// Attach authentication details
	userId := getUserIdFromToken(c)
	user, err := db.FindUserById(invoker.DB, userId)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"message": fmt.Sprintf("find user by id error: %s", err.Error()),
		})
	}

	auth := utils.GetAuth(userId)
	result["token"] = auth.WebofficeToken
	result["signature"] = invoker.SdkMgr.Sign(5*time.Minute, sdkapi.ScopeDefault)
	result["user"] = user

	// Return the filtered result
	c.JSON(200, result)
}

func NewHttpRequest(c *gin.Context, url string, method string) ([]byte, error) {
	// Create the HTTP request
	rq, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %s", err.Error())
	}
	// Send the HTTP request
	rs, err := http.DefaultClient.Do(rq)
	if err != nil {
		return nil, fmt.Errorf("get %s failed: %s", url, err.Error())
	}
	defer rs.Body.Close()

	// Validate the response status
	if rs.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get %s failed: %s", url, rs.Status)
	}
	// Read the response body
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		return nil, fmt.Errorf("read %s response body failed: %s", url, err.Error())
	}
	return body, nil
}

// ListKnowledgeBases returns the list of knowledge bases
func ListKnowledgeBases(c *gin.Context) {
	guids, err := db.FindAllKnowledgeBaseGuids(invoker.DB)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to get knowledge bases", "error": err.Error()})
		return
	}
	// Fetch the file count for each knowledge base
	type KnowledgeBaseInfo struct {
		Guid      string `json:"guid"`
		FileCount int64  `json:"fileCount"`
		CreateAt  int64  `json:"createAt"`
		Name      string `json:"name"`
	}
	var knowledgeBases []KnowledgeBaseInfo
	for _, guid := range guids {
		// Count the files in the knowledge base
		var count int64
		invoker.DB.Model(&db.KnowledgeBase{}).Where("guid = ? and delete_at = 0 and file_guid != ''", guid).Count(&count)

		// Fetch the creation time
		var kb db.KnowledgeBase
		invoker.DB.Where("guid = ? and delete_at = 0", guid).First(&kb)

		knowledgeBases = append(knowledgeBases, KnowledgeBaseInfo{
			Guid:      guid,
			FileCount: count,
			CreateAt:  kb.CreatedAt,
			Name:      kb.Name,
		})
	}
	c.JSON(200, gin.H{"data": knowledgeBases})
}

func GetKnowledgeBase(c *gin.Context) {
	guid := c.Param("knowledgeBaseGuid")
	if guid == "" {
		c.JSON(400, gin.H{"message": "Knowledge base guid is required"})
		return
	}
	base := db.KnowledgeBase{}
	if err := invoker.DB.Where("guid = ?", guid).First(&base).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to get knowledge base",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, base)
}

// CreateKnowledgeBase creates a knowledge base
func CreateKnowledgeBase(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}
	userId := getUserIdFromToken(c)
	// Create an empty knowledge base record
	kb := &db.KnowledgeBase{
		Guid:     utils.GenFileGuid(),
		FileGuid: "", // Empty knowledge base
		Name:     req.Name,
		CreateBy: userId,
	}
	if err := db.CreateKnowledgeBase(invoker.DB, kb); err != nil {
		c.JSON(500, gin.H{"message": "Failed to create knowledge base", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Knowledge base created successfully", "data": kb})
}

// DeleteKnowledgeBase removes a knowledge base
func DeleteKnowledgeBase(c *gin.Context) {
	guid := c.Param("knowledgeBaseGuid")
	if guid == "" {
		c.JSON(400, gin.H{"message": "Knowledge base guid is required"})
		return
	}
	if err := invoker.DB.Where("guid = ?", guid).Delete(&db.KnowledgeBase{}).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to delete knowledge base",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"message": "Knowledge base deleted successfully"})
}

type ImportFileToKnowledgeBaseReq struct {
	KnowledgeBaseGuid string `json:"knowledgeBaseGuid" required:"true"` // The integrator's knowledge base ID
	// Either file or url must be provided
	// Use file to import a Shimo document; the system uses fileGuid to import it
	// Use url to import an external document; the system downloads it via downloadUrl
	ImportType string `json:"importType" required:"true"`
	FileGuid   string `json:"fileGuid"`
	// File type rules:
	// 1. If importType is file, pass a Shimo document type (document, documentPro, spreadsheet, presentation)
	// 2. If importType is url, pass an external type (currently only pdf or rtf)
	FileType string `json:"fileType" required:"true"`
	// Only effective when importType is "file"; the cloud file URL must be reachable to Shimo servers
	// Note: the address must be accessible from Shimo's infrastructure. Example: https://example.com/document.pdf
	DownloadUrl string `json:"downloadUrl"`
}

func ImportFileToKnowledgeBase(c *gin.Context) {
	params := ImportFileToKnowledgeBaseReq{}
	err := c.BindJSON(&params)
	if err != nil {
		c.JSON(400, gin.H{"message": "params invalid"})
		return
	}
	fileGuid := params.FileGuid
	userId := getUserIdFromToken(c)
	// For url imports, create a cloud file locally first
	if params.ImportType == "url" {
		var content []byte
		// ---- Extract the file name from the URL ----
		u, err := url.Parse(params.DownloadUrl)
		var fileName string
		if err == nil {
			segments := strings.Split(u.Path, "/")
			if len(segments) > 0 {
				fileName = segments[len(segments)-1]
				// Strip query parameters (everything after '?')
				if idx := strings.Index(fileName, "?"); idx != -1 {
					fileName = fileName[:idx]
				}
			}
		}
		file := db.File{
			Name:        fileName,
			ShimoType:   "",
			IsShimoFile: 0,
			CreatorId:   getUserIdFromToken(c),
		}
		// Download the content from fileUrl
		rq, _ := http.NewRequest("GET", params.DownloadUrl, nil)
		rs, e := http.DefaultClient.Do(rq)
		if e != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"message": "file download failed",
				"error":   e,
			})
			return
		}
		defer rs.Body.Close()
		// Read the file content
		content, e = io.ReadAll(rs.Body)
		if e != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"message": "file content read failed",
				"error":   e,
			})
			return
		}
		// Use the same MIME detection logic as direct uploads
		mimeType := DetectMimeType(fileName, content)
		// Update the file type in the database
		file.Type = mimeType
		e, _ = db.CreateFile(invoker.DB, &file, getUserIdFromToken(c))
		if e != nil {
			handleDBError(c, e)
			return
		}
		fileGuid = file.Guid
	}
	// Insert a row into the AI knowledge base table
	auth := utils.GetAuth(userId)
	iFile := sdkapi.ImportFileToAiKnowledgeBaseReq{
		Metadata: auth,
		ImportFileToAiKnowledgeBaseReqBody: sdkapi.ImportFileToAiKnowledgeBaseReqBody{
			KnowledgeBaseGuid: params.KnowledgeBaseGuid,
			ImportType:        params.ImportType,
			FileGuid:          fileGuid,
			FileType:          params.FileType,
			DownloadUrl:       params.DownloadUrl,
		},
	}
	res, sdkErr := invoker.SdkMgr.ImportFileToAiKnowledgeBase(c.Request.Context(), iFile)
	if sdkErr != nil {
		c.JSON(res.Response().StatusCode(), gin.H{"error": sdkErr.Error()})
		return
	}
	// Insert a row into the local knowledge base
	kb := &db.KnowledgeBase{
		Guid:     params.KnowledgeBaseGuid,
		FileGuid: fileGuid,
		CreateBy: userId,
	}
	if dbErr := db.CreateKnowledgeBase(invoker.DB, kb); dbErr != nil {
		c.JSON(500, gin.H{"message": "Failed to save to database", "error": dbErr.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "File imported successfully", "data": kb})
}

// GetKnowledgeBaseFiles returns the files inside a knowledge base
func GetKnowledgeBaseFiles(c *gin.Context) {
	guid := c.Param("knowledgeBaseGuid")
	if guid == "" {
		c.JSON(400, gin.H{"message": "Knowledge base guid is required"})
		return
	}

	// Retrieve every file record for the knowledge base
	knowledgeBases, err := db.FindKnowledgeBasesByGuid(invoker.DB, guid)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to get knowledge base files", "error": err.Error()})
		return
	}

	// Build the detailed file list
	type FileInfo struct {
		ID        int64  `json:"id"`
		Guid      string `json:"guid"`
		Name      string `json:"name"`
		Type      string `json:"type"`
		ShimoType string `json:"shimoType"`
		CreateAt  int64  `json:"createAt"`
		CreateBy  int64  `json:"createBy"`
	}

	var files []FileInfo
	for _, kb := range knowledgeBases {
		if kb.FileGuid != "" {
			var file db.File
			if fileErr := invoker.DB.Where("guid = ?", kb.FileGuid).First(&file).Error; fileErr == nil {
				files = append(files, FileInfo{
					ID:        file.ID,
					Guid:      file.Guid,
					Name:      file.Name,
					Type:      file.Type,
					ShimoType: file.ShimoType,
					CreateAt:  file.CreatedAt,
					CreateBy:  file.CreatorId,
				})
			}
		}
	}

	c.JSON(200, gin.H{"data": files})
}

func DeleteFileFromKnowledgeBase(c *gin.Context) {
	guid := c.Param("knowledgeBaseGuid")
	fileGuid := c.Param("fileGuid")
	if guid == "" || fileGuid == "" {
		c.JSON(400, gin.H{"message": "Knowledge base guid and file guid are required"})
		return
	}
	userId := getUserIdFromToken(c)
	auth := utils.GetAuth(userId)
	params := sdkapi.DeleteFileFromAiKnowledgeBaseReq{
		Metadata: auth,
		DeleteFileFromAiKnowledgeBaseReqBody: sdkapi.DeleteFileFromAiKnowledgeBaseReqBody{
			KnowledgeBaseGuid: guid,
			FileGuid:          fileGuid,
		},
	}

	res, err := invoker.SdkMgr.DeleteFileFromAiKnowledgeBase(c.Request.Context(), params)
	if err != nil {
		c.JSON(res.Response().StatusCode(), gin.H{"error": err.Error()})
		return
	}
	// Delete the local data
	if err = invoker.DB.Where("guid = ? AND file_guid = ?", guid, fileGuid).Delete(&db.KnowledgeBase{}).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to delete knowledge base",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"message": "File removed from knowledge base successfully"})
}

func ImportFileToKnowledgeBaseV2(c *gin.Context) {
	params := ImportFileToKnowledgeBaseReq{}
	err := c.BindJSON(&params)
	if err != nil {
		c.JSON(400, gin.H{"message": "params invalid"})
		return
	}
	fileGuid := params.FileGuid
	userId := getUserIdFromToken(c)

	// For url imports, create a cloud file locally first
	if params.ImportType == "url" {
		var content []byte
		// ---- Extract the file name from the URL ----
		u, err := url.Parse(params.DownloadUrl)
		var fileName string
		if err == nil {
			segments := strings.Split(u.Path, "/")
			if len(segments) > 0 {
				fileName = segments[len(segments)-1]
				// Strip query parameters (everything after '?')
				if idx := strings.Index(fileName, "?"); idx != -1 {
					fileName = fileName[:idx]
				}
			}
		}
		file := db.File{
			Name:        fileName,
			ShimoType:   "",
			IsShimoFile: 0,
			CreatorId:   getUserIdFromToken(c),
		}
		// Download the content from fileUrl
		rq, _ := http.NewRequest("GET", params.DownloadUrl, nil)
		rs, e := http.DefaultClient.Do(rq)
		if e != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"message": "file download failed",
				"error":   e,
			})
			return
		}
		defer rs.Body.Close()
		// Read the file content
		content, e = io.ReadAll(rs.Body)
		if e != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"message": "file content read failed",
				"error":   e,
			})
			return
		}
		// Use the same MIME detection logic as direct uploads
		mimeType := DetectMimeType(fileName, content)
		// Update the file type in the database
		file.Type = mimeType
		e, _ = db.CreateFile(invoker.DB, &file, getUserIdFromToken(c))
		if e != nil {
			handleDBError(c, e)
			return
		}
		fileGuid = file.Guid
	}

	// Call the v2 SDK API
	auth := utils.GetAuth(userId)
	iFile := sdkapi.ImportFileToAiKnowledgeBaseV2Req{
		Metadata: auth,
		ImportFileToAiKnowledgeBaseV2ReqBody: sdkapi.ImportFileToAiKnowledgeBaseV2ReqBody{
			KnowledgeBaseGuid: params.KnowledgeBaseGuid,
			ImportType:        params.ImportType,
			FileGuid:          fileGuid,
			FileType:          params.FileType,
			DownloadUrl:       params.DownloadUrl,
		},
	}
	res, sdkErr := invoker.SdkMgr.ImportFileToAiKnowledgeBaseV2(c.Request.Context(), iFile)
	if sdkErr != nil {
		c.JSON(res.Response().StatusCode(), gin.H{"error": sdkErr.Error()})
		return
	}

	// Return the task ID so the front end can poll progress
	c.JSON(200, gin.H{
		"taskId":            res.TaskID,
		"fileGuid":          fileGuid,
		"knowledgeBaseGuid": params.KnowledgeBaseGuid,
	})
}

func ImportFileToKnowledgeBaseV2Progress(c *gin.Context) {
	var req struct {
		TaskId            string `json:"taskId" binding:"required"`
		FileGuid          string `json:"fileGuid" binding:"required"`
		KnowledgeBaseGuid string `json:"knowledgeBaseGuid" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "taskId is required", "error": err.Error()})
		return
	}

	userId := getUserIdFromToken(c)
	auth := utils.GetAuth(userId)

	// Ask the SDK for the current progress
	params := sdkapi.GetImportFileToAiProgressV2Req{
		Metadata: auth,
		GetImportFileToAiProgressV2ReqBody: sdkapi.GetImportFileToAiProgressV2ReqBody{
			TaskID: req.TaskId,
		},
	}

	res, sdkErr := invoker.SdkMgr.GetImportFileToAiProgressV2(c.Request.Context(), params)
	if sdkErr != nil {
		c.JSON(res.Response().StatusCode(), gin.H{"error": sdkErr.Error()})
		return
	}

	// If the task completes, add the file to the local knowledge base
	if res.Status == "completed" && res.Progress == 100 {
		// Fetch the file information as needed and persist it locally
		// The v2 API response might lack file details, so other approaches may be required
		// For now, return the progress; adjust the save logic based on the actual API payload
		// Insert a row into the local knowledge base
		kb := &db.KnowledgeBase{
			Guid:     req.KnowledgeBaseGuid,
			FileGuid: req.FileGuid,
			CreateBy: userId,
		}
		if dbErr := db.CreateKnowledgeBase(invoker.DB, kb); dbErr != nil {
			c.JSON(500, gin.H{"message": "Failed to save to database", "error": dbErr.Error()})
			return
		}
	}

	c.JSON(200, gin.H{
		"taskId":   res.TaskId,
		"status":   res.Status,
		"progress": res.Progress,
		"message":  res.Message,
	})
}

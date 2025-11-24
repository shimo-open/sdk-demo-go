package callback

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	sdkapi "github.com/shimo-open/sdk-kit-go/model/api"
	sdkcommon "github.com/shimo-open/sdk-kit-go/model/common"

	"sdk-demo-go/pkg/consts"
	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/models/db"
	"sdk-demo-go/pkg/server/http/api"
	"sdk-demo-go/pkg/utils"
)

type FileInfo struct {
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	Type        string          `json:"type"`
	CreatorId   string          `json:"creatorId"`
	Views       int             `json:"views"`
	CreatedAt   string          `json:"createdAt"`
	UpdatedAt   string          `json:"updatedAt"`
	TeamGuid    string          `json:"teamGuid"`
	Permissions map[string]bool `json:"permissions"`
}

func GetFileInfo(c *gin.Context) {
	userId := getUserIdFromToken(c)
	if userId == 0 {
		userId = consts.ANONYMOUS
	}
	fileGuid := c.Param("fileGuid")

	file, err := db.FindFileByGuidAndUserId(invoker.DB, userId, fileGuid)
	if err != nil {
		handleDBError(c, err)
		return
	}

	if file.IsShimoFile == 1 {
		sendShimoInfo(c, file)
	} else {
		sendFileInfo(c, file)
	}
}

type CreateFileReq struct {
	CreateCopyInfo *CreateReqCopyInfo `json:"createCopyInfo,omitempty"`
	CreateLinkInfo *CreateReqLinkInfo `json:"createLinkInfo,omitempty"`
	ContentKey     string             `json:"contentKey"`
	Action         string             `json:"action"`
	FileType       string             `json:"fileType"`
}
type CreateReqLinkInfo struct {
	SourceFileID string `json:"sourceFileId"`
	ParentFileID string `json:"parentFileId"`
	NewFileID    string `json:"newFileId,omitempty"`   // SaaS uses this to pass a specific GUID when creating a file
	NewFileName  string `json:"newFileName,omitempty"` // SaaS uses this to pass a specific file name when creating a file
}

type CreateReqCopyInfo struct {
	SourceFileID string `json:"sourceFileId"`
	ParentFileID string `json:"parentFileId"`
	NewFileID    string `json:"newFileId,omitempty"`   // SaaS uses this to pass a specific GUID when creating a file
	NewFileName  string `json:"newFileName,omitempty"` // SaaS uses this to pass a specific file name when creating a file
}

func CreateFiles(c *gin.Context) {
	userId := getUserIdFromToken(c)
	body := CreateFileReq{}
	lang, ok := c.GetQuery("lang")
	if !ok {
		lang = "zh-CN"
	}
	err := c.BindJSON(&body)
	if err != nil {
		return
	}
	name := body.CreateLinkInfo.NewFileName
	if name == "" {
		name = body.CreateCopyInfo.NewFileName
	}
	if name == "" {
		name = api.FormatCurrentTime() + " " + body.FileType
	}
	file := db.File{
		Name:        name,
		ShimoType:   body.FileType,
		CreatorId:   userId,
		IsShimoFile: api.IsShimoType(body.FileType),
	}
	// Create the local metadata record
	err, _ = db.CreateFile(invoker.DB, &file, userId)
	if err != nil {
		handleDBError(c, err)
		return
	}
	if body.ContentKey != "" {
		// Create the SDK file
		auth := utils.GetAuth(userId)
		cFile := sdkapi.CreateFileParams{
			FileType:   sdkcommon.CollabFileType(file.ShimoType),
			Auth:       auth,
			Lang:       sdkcommon.Lang(lang),
			FileId:     file.Guid,
			ContentKey: body.ContentKey,
		}
		_, err = invoker.SdkMgr.CreateFile(cFile)
		if err != nil {
			// Delete the record if creation fails
			_ = db.RemoveFileById(invoker.DB, file.ID)
			handleDBError(c, err)
			return
		}
	}
	fileInfo := loadFileInfo(&file)
	elog.Debug("sendCreateShimoInfo", l.A("fileInfo", fileInfo))
	c.JSON(200, fileInfo)
}

func GetFilesByUser(c *gin.Context) {
	userId := getUserIdFromToken(c)

	_limit, _ := c.GetQuery("limit")
	orderBy, _ := c.GetQuery("orderBy")
	limit, _ := strconv.Atoi(_limit)

	files, err := db.FindFileByUserId(invoker.DB, userId, limit, orderBy)
	if err != nil {
		handleDBError(c, err)
		return
	}

	fileInfos := make([]FileInfo, len(files))
	for i := range files {
		fileInfos[i] = *loadFileInfo(&files[i])
	}

	c.JSON(200, fileInfos)
}

type CollaboratorInfo struct {
	UserInfo
	IsManager bool
}

func GetFileCollaborators(c *gin.Context) {
	fileGuid := c.Param("fileGuid")
	file, err := db.FindFileByGuid(invoker.DB, fileGuid)
	if err != nil {
		handleDBError(c, err)
		return
	}

	fps, err := db.FindFilePermissionsByFileId(invoker.DB, file.ID)
	if err != nil {
		handleDBError(c, err)
		return
	}

	var userIds []int64
	managerMap := make(map[int64]bool)
	for i := range fps {
		userIds = append(userIds, fps[i].UserId)
		if fps[i].Permissions["manageable"] {
			managerMap[fps[i].UserId] = true
		}
	}

	users, err := db.FindUsersByIds(invoker.DB, userIds)

	collInfos := make([]CollaboratorInfo, len(users))
	for i := range users {
		collInfos[i] = CollaboratorInfo{
			UserInfo: UserInfo{
				Id:        strconv.FormatInt(users[i].ID, 10),
				Name:      users[i].Name,
				Avatar:    users[i].Avatar,
				Email:     users[i].Email,
				CanBother: users[i].CanBother,
			},
			IsManager: managerMap[users[i].ID],
		}
	}
	c.JSON(200, collInfos)
}

func GetFileAccessUrl(c *gin.Context) {
	fileGuid := c.Param("fileGuid")

	outerHost := econf.GetString("outerHost")
	if outerHost != "" {
		c.JSON(200, gin.H{
			"url": fmt.Sprintf("%s/collab-files/%d", outerHost, fileGuid),
		})
		return
	}

	host := econf.GetString("host")
	c.JSON(200, gin.H{
		"url": fmt.Sprintf("%s/collab-files/%d", host, fileGuid),
	})
}

func AdminGetFileInfo(c *gin.Context) {
	fileGuid := c.Param("fileGuid")

	file, err := db.FindFileByGuid(invoker.DB, fileGuid)
	if err != nil {
		handleDBError(c, err)
		return
	}

	if econf.GetBool("permissions.setNewFilePermission") {
		file.Permissions = consts.HandleFilePermission(true)
	} else {
		file.Permissions = consts.HandleBasicFilePermission(true)
	}

	if file.IsShimoFile == 1 {
		sendShimoInfo(c, file)
	} else {
		sendFileInfo(c, file)
	}
}

func AdminGetFileInfoByUserId(c *gin.Context) {
	fileGuid := c.Param("fileGuid")

	_userId, ok := c.GetQuery("userId")
	if !ok {
		c.JSON(400, gin.H{
			"message": "userId is required",
		})
		return
	}
	userId, _ := strconv.ParseInt(_userId, 10, 64)

	file, err := db.FindFileByGuidAndUserId(invoker.DB, userId, fileGuid)
	if err != nil {
		handleDBError(c, err)
		return
	}

	if file.IsShimoFile == 1 {
		sendShimoInfo(c, file)
	} else {
		sendFileInfo(c, file)
	}
}

func DownloadFile(c *gin.Context) {
	fileGuid := c.Param("fileGuid")
	file, err := db.FindFileByGuid(invoker.DB, fileGuid)
	if err != nil {
		handleDBError(c, err)
		return
	}

	bytes, err := invoker.Services.AwosService.Get(fileGuid)
	if err != nil {
		c.JSON(500, gin.H{"message": "awos get file error"})
		return
	}
	fileName := url.QueryEscape(file.Name)

	c.Header("Content-Type", file.Type)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Header("Content-Length", strconv.Itoa(len(bytes)))
	c.Data(200, file.Type, bytes)
}

func getSDKClaims(c *gin.Context) *utils.SDKClaims {
	_claims, e := c.Get("claims")
	if !e {
		return nil
	}

	claims, _ := _claims.(*utils.SDKClaims)
	return claims
}

/*
{
    "type": "document",
    "fileUrl": "https://obs.xxxxxx.com/smdev-svc-drive/src?X-AMX-XXXXX",
    "fileName": "亲子旅行策划",
    "sourceFileGuid": "e1AzdpJ6nOUjvoqW"
}
*/

type ImportFileReq struct {
	Type           string `json:"type"`
	FileUrl        string `json:"fileUrl"`
	FileName       string `json:"fileName"`
	SourceFileGuid string `json:"sourceFileGuid"`
}

func ImportFile(c *gin.Context) {
	req := ImportFileReq{}
	err := c.BindJSON(&req)
	if err != nil {
		return
	}

	elog.Debug("ImportFile", l.A("body", req))

	name := req.FileName
	fileUrl := req.FileUrl
	shimoType := req.Type
	var content []byte

	isShimoFile := 1
	if req.Type == "file" {
		isShimoFile = 0
		shimoType = ""
	}

	// Create the source file
	file := db.File{
		Name:        name,
		ShimoType:   shimoType,
		IsShimoFile: isShimoFile,
		CreatorId:   getUserIdFromToken(c),
	}
	if req.Type == "file" {
		// Content is downloaded from fileUrl
		rq, _ := http.NewRequest("GET", fileUrl, nil)
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
		mimeType := api.DetectMimeType(name, content)
		// Update the file type information in the database
		file.Type = mimeType
	}
	err, _ = db.CreateFile(invoker.DB, &file, getUserIdFromToken(c))
	if err != nil {
		handleDBError(c, err)
		return
	}

	switch req.Type {
	case "file":
		// Cloud file
		// Generate a temporary download link
		uploadUrl, err := invoker.Services.AwosService.GetUploadURL(file.Guid, 1800)
		if err != nil {
			// Roll back the created file
			rmErr := db.RemoveFileById(invoker.DB, file.ID)
			if rmErr != nil {
				elog.Warn("rollback file failed", l.E(rmErr))
			}
			c.AbortWithStatusJSON(500, gin.H{
				"message": "GetUploadURL failed",
				"error":   err,
			})
			return
		}
		// Use the download link to persist the file into storage
		uploadReq, _ := http.NewRequest("PUT", uploadUrl, bytes.NewReader(content))
		_, err = http.DefaultClient.Do(uploadReq)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"message": "file save failed",
				"error":   err,
			})
			elog.Error("file save failed", l.E(err))
			return
		}
		c.JSON(200, file)
		return
	default:
		// Import into Shimo

		// Create the SDK file
		auth := utils.GetAuth(getUserIdFromToken(c))
		importBody := sdkapi.ImportFileReqBody{
			FileId:   file.Guid,
			Type:     file.ShimoType,
			FileName: file.Name,
			FileUrl:  fileUrl,
		}
		params := sdkapi.ImportFileParams{
			Auth:              auth,
			ImportFileReqBody: importBody,
		}
		var res sdkapi.ImportFileRespBody
		if econf.GetString("shimoSDK.importByUrlVersion") == "v2" {
			res, err = invoker.SdkMgr.ImportV2File(params)
		} else {
			res, err = invoker.SdkMgr.ImportFile(params)
		}
		if err != nil || res.Status != 0 {
			rmErr := db.RemoveFileByGuid(invoker.DB, file.Guid)
			if rmErr != nil {
				elog.Warn("rollback file failed", l.E(rmErr))
			}
			e := fmt.Errorf("import file failed, got: %d, resp body: %s", res.Resp.StatusCode(), string(res.Resp.Body()))
			c.JSON(res.Resp.StatusCode(), e.Error())
			return
		}

		taskId := res.Data.TaskId
		if taskId == "" {
			c.JSON(500, "taskId not found")
			return
		}
		// Poll the import progress
		var progressResp sdkapi.GetImportProgRespBody
		progressParams := sdkapi.GetImportProgParams{
			Auth:   auth,
			TaskId: taskId,
		}
		timeout := time.After(3 * time.Minute) // Wait at most 2 minutes
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-timeout:
				rmErr := db.RemoveFileByGuid(invoker.DB, file.Guid)
				if rmErr != nil {
					elog.Warn("rollback file failed", l.E(rmErr))
				}
				c.JSON(504, "import progress timeout")
				return
			case <-ticker.C:
				if econf.GetString("shimoSDK.importByUrlVersion") == "v2" {
					progressResp, err = invoker.SdkMgr.GetImportV2Progress(progressParams)
				} else {
					progressResp, err = invoker.SdkMgr.GetImportProgress(progressParams)
				}
				if err != nil {
					rmErr := db.RemoveFileByGuid(invoker.DB, file.Guid)
					if rmErr != nil {
						elog.Warn("rollback file failed", l.E(rmErr))
					}
					c.JSON(500, fmt.Errorf("get import progress failed: %w", err))
					return
				}
				if progressResp.Data.Progress == 100 {
					resp := struct {
						db.File
					}{
						File: file,
					}
					c.JSON(200, resp)
					return
				}
			}
		}
	}
}

func sendShimoInfo(c *gin.Context, file *db.File) {
	mode := getModeFromToken(c)
	userId := getUserIdFromToken(c)

	if file.ShimoType == "form" {
		permissions := map[string]bool{
			"readable":     file.Permissions["readable"],
			"formFillable": true,
			"editable":     file.Permissions["editable"],
			"manageable":   file.Permissions["manageable"],
			"commentable":  file.Permissions["commentable"],
		}

		if userId < 0 {
			permissions = map[string]bool{
				"readable":     false,
				"formFillable": true,
				"commentable":  false,
				"editable":     false,
			}
		}
		file.Permissions = permissions
	}

	fileInfo := loadFileInfo(file)

	allPermissions := []string{
		"readable",
		"formFillable",
		"editable",
		"manageable",
		"commentable",
	}

	// Populate any missing permission flags
	for _, p := range allPermissions {
		if _, ok := fileInfo.Permissions[p]; !ok {
			fileInfo.Permissions[p] = false
		}
	}

	elog.Debug("sendShimoInfo", l.S("mode", mode), l.A("fileInfo", fileInfo))
	c.JSON(200, fileInfo)
}

func sendFileInfo(c *gin.Context, file *db.File) {
	downloadUrl, err := invoker.Services.AwosService.GetDownloadURL(file.Guid, file.Name, 3600)
	if err != nil {
		c.JSON(500, gin.H{"message": "awos get file error"})
		return
	}

	// Replace the download URL prefix
	publicEndpointReplacement := econf.GetString("awos.publicEndpointReplacement")
	awosEndpoint := econf.GetString("awos.endpoint")

	if publicEndpointReplacement != "" && awosEndpoint != "" {
		downloadUrl = strings.Replace(downloadUrl, awosEndpoint, publicEndpointReplacement, 1)
	}

	extension := filepath.Ext(file.Name)
	if len(extension) > 0 {
		extension = strings.ToLower(extension[1:])
	}

	c.JSON(200, gin.H{
		"id":   file.Guid,
		"name": file.Name,
		"type": "file",
		"permissions": gin.H{
			"readable": file.Permissions["readable"],
		},
		"downloadUrl": downloadUrl,
		"ext":         extension,
	})
}

func loadFileInfo(file *db.File) *FileInfo {
	t := time.Unix(file.CreatedAt, 0)
	utcTime := t.UTC()
	ctime := utcTime.Format(time.RFC3339)

	t = time.Unix(file.UpdatedAt, 0)
	utcTime = t.UTC()
	utime := utcTime.Format(time.RFC3339)

	// fullUrl := ""
	// if file.Type == "file" {
	//	fullUrl = "test"
	// }

	var typ string
	if file.IsShimoFile == 1 {
		typ = file.ShimoType
	} else {
		typ = "file"
	}

	return &FileInfo{
		Id:          file.Guid,
		Name:        file.Name,
		Type:        typ,
		Permissions: file.Permissions,
		Views:       0,
		CreatorId:   strconv.FormatInt(file.CreatorId, 10),
		CreatedAt:   ctime,
		UpdatedAt:   utime,
		TeamGuid:    "123",
	}
}

// GetFileUrl returns the file redirect URL
func GetFileUrl(c *gin.Context) {
	userId := getUserIdFromToken(c)
	if userId == 0 {
		userId = consts.ANONYMOUS
	}
	fileGuid := c.Param("fileGuid")
	var u string
	file, err := db.FindFileByGuidAndUserId(invoker.DB, userId, fileGuid)
	if err != nil {
		handleDBError(c, err)
		return
	}
	addr := econf.GetString("host.addr")
	// Standard Shimo file types
	if file.IsShimoFile == 1 {
		// Handle forms separately
		if file.ShimoType == "form" {
			u = fmt.Sprintf("%s/form/%s/fill", addr, file.Guid)
		} else {
			u = fmt.Sprintf("%s/shimo-files/%s", addr, file.Guid)
		}
	} else {
		// Preview non-Shimo documents
		u = fmt.Sprintf("%s/preview/%s", addr, file.Guid)
	}
	c.JSON(200, gin.H{
		"url": u,
	})
}

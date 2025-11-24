package api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/h2non/filetype"
	sdkapi "github.com/shimo-open/sdk-kit-go/model/api"
	sdkcommon "github.com/shimo-open/sdk-kit-go/model/common"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	"sdk-demo-go/pkg/consts"
	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/models/db"
	"sdk-demo-go/pkg/server/http/middlewares"
	"sdk-demo-go/pkg/utils"
)

func GetUserFiles(c *gin.Context) {
	userId := getUserIdFromToken(c)
	files, err := db.FindFileByUserId(invoker.DB, userId, 0, "")
	if err != nil {
		handleDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, files)
}

func GetFileThumbnail(c *gin.Context) {
	// TODO no SDK API is available yet for thumbnails
	c.HTML(501, "thumbnail", gin.H{
		"title":        "缩略图",
		"thumbnailUrl": "https://cdn.shimo.im/thumbnail/1.jpg",
	})
}

func OpenFile(c *gin.Context) {
	fileGuid := c.Param("fileGuid")

	file, err := db.FindFileByGuid(invoker.DB, fileGuid)
	if err != nil {
		handleDBError(c, err)
		return
	}

	appId := econf.GetString("shimoSDK.appId")
	secret := econf.GetString("shimoSDK.appSecret")
	userId := c.GetInt64("userId")

	c.HTML(200, "shimo-file", gin.H{
		"rootCSSClasses": "editor-page",
		"pageTitle":      file.Name,
		"file":           file,
		"config": gin.H{
			"fileId":    file.Guid,
			"signature": invoker.Services.SignatureService.Sign(appId, secret, false),
			"appId":     appId,
			"endpoint":  econf.GetString("shimoSDK.host"),
			"token":     utils.SignUserJWT(userId),
		},
	})
}

func GetPlainText(c *gin.Context) {
	fileGuid := c.Param("fileGuid")

	file, err := db.FindFileByGuid(invoker.DB, fileGuid)
	if err != nil {
		handleDBError(c, err)
		return
	}

	auth := utils.GetAuth(getUserIdFromToken(c))
	params := sdkapi.GetPlainTextParams{
		FileId: fileGuid,
		Auth:   auth,
	}
	resp, err := invoker.SdkMgr.GetPlainText(params)
	if err != nil {
		handleSdkMgrError(c, resp.Resp.Body(), resp.Resp.StatusCode())
		return
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename="+file.Name+".txt")
	c.String(200, resp.Content)
}

type PaginationQuery struct {
	Pagination bool `form:"pagination"`
	Page       int  `form:"page"`
	PageSize   int  `form:"pageSize"`
}

func GetDocSidebar(c *gin.Context) {
	query := PaginationQuery{}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "query params error"})
		return
	}
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}

	fileGuid := c.Param("fileGuid")

	_, err := db.FindFileByGuid(invoker.DB, fileGuid)
	if err != nil {
		handleDBError(c, err)
		return
	}
	auth := utils.GetAuth(getUserIdFromToken(c))
	params := sdkapi.GetHistoryListParams{
		Auth:     auth,
		FileId:   fileGuid,
		PageSize: query.PageSize,
		Count:    (query.Page - 1) * query.PageSize,
	}
	resp, err := invoker.SdkMgr.GetHistoryList(params)
	if err != nil {
		handleSdkMgrError(c, resp.Resp.Body(), resp.Resp.StatusCode())
		return
	}
	c.JSON(200, resp)
}

func GetFileRevision(c *gin.Context) {
	fileGuid := c.Param("fileGuid")

	_, err := db.FindFileByGuid(invoker.DB, fileGuid)
	if err != nil {
		handleDBError(c, err)
		return
	}

	auth := utils.GetAuth(getUserIdFromToken(c))
	params := sdkapi.GetRevisionListParams{
		Auth:   auth,
		FileId: fileGuid,
	}
	resp, res, err := invoker.SdkMgr.GetRevisionList(params)
	if err != nil {
		handleSdkMgrError(c, resp.Resp.Body(), resp.Resp.StatusCode())
		return
	}
	c.JSON(200, res)
}

func CountComments(c *gin.Context) {
	fileGuid := c.Param("fileGuid")
	_, err := db.FindFileByGuid(invoker.DB, fileGuid)
	if err != nil {
		handleDBError(c, err)
		return
	}
	auth := utils.GetAuth(getUserIdFromToken(c))
	params := sdkapi.GetCommentCountParams{
		Auth:   auth,
		FileId: fileGuid,
	}
	resp, err := invoker.SdkMgr.GetCommentCount(params)
	if err != nil {
		handleSdkMgrError(c, resp.Resp.Body(), resp.Resp.StatusCode())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func GetMentionAtList(c *gin.Context) {
	fileGuid := c.Param("fileGuid")

	file, err := db.FindFileByGuid(invoker.DB, fileGuid)
	if err != nil {
		handleDBError(c, err)
		return
	}
	auth := utils.GetAuth(getUserIdFromToken(c))
	params := sdkapi.GetMentionAtParams{
		Auth:   auth,
		FileId: file.Guid,
	}
	resp, err := invoker.SdkMgr.GetMentionAt(params)
	if err != nil {
		handleSdkMgrError(c, resp.Resp.Body(), resp.Resp.StatusCode())
		return
	}

	uids := make([]int64, 0)
	for _, m := range resp.MentionAtList {
		uids = append(uids, cast.ToInt64(m.UserId))
	}

	users, err := db.FindUsersByIds(invoker.DB, uids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "get mention at list db err" + err.Error()})
		return
	}

	finalResp := make([]struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}, 0)
	for _, u := range users {
		finalResp = append(finalResp, struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		}{
			Id:   cast.ToString(u.ID),
			Name: u.Name,
		})
	}
	c.JSON(http.StatusOK, finalResp)
	return
}

type FileBody struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	ShimoType string `json:"shimoType"`
	Url       string `json:"url"`
}

func CreateFile(c *gin.Context) {
	body := FileBody{}
	err := c.BindJSON(&body)
	if err != nil {
		return
	}

	lang, ok := c.GetQuery("lang")
	if !ok {
		lang = "zh-CN"
	}

	name := body.Name
	if name == "" {
		name = FormatCurrentTime() + " " + body.ShimoType
	}

	userId := getUserIdFromToken(c)

	file := db.File{
		Name:        name,
		Type:        body.Type,
		ShimoType:   body.ShimoType,
		CreatorId:   userId,
		IsShimoFile: IsShimoType(body.ShimoType),
	}
	var fileId int64

	err, fileId = db.CreateFile(invoker.DB, &file, userId)
	if err != nil {
		handleDBError(c, err)
		return
	}

	auth := utils.GetAuth(userId)
	cFile := sdkapi.CreateFileParams{
		FileType: sdkcommon.CollabFileType(file.ShimoType),
		Auth:     auth,
		Lang:     sdkcommon.Lang(lang),
		FileId:   file.Guid,
	}
	res, err := invoker.SdkMgr.CreateFile(cFile)
	if err != nil {
		// Delete the record if creation fails
		_ = db.RemoveFileById(invoker.DB, fileId)
		handleSdkMgrError(c, res.Resp.Body(), res.Resp.StatusCode())
		return
	}

	c.JSON(200, file)
}

func UploadFile(c *gin.Context) {
	_file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"message": "file update failed"})
		return
	}

	file, err := _file.Open()
	if err != nil {
		c.JSON(400, gin.H{"message": "file open failed"})
		return

	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		c.JSON(400, gin.H{"message": "file read failed"})
	}

	mimeType := DetectMimeType(_file.Filename, content)
	fileName := _file.Filename
	userId := getUserIdFromToken(c)

	// TODO determine the default MIME type
	f := db.File{
		Name:      fileName,
		Type:      mimeType,
		CreatorId: userId,
	}

	err, _ = db.CreateFile(invoker.DB, &f, userId)
	if err != nil {
		handleDBError(c, err)
		return
	}
	// Generate a temporary download link
	uploadUrl, err := invoker.Services.AwosService.GetUploadURL(f.Guid, 1800)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"message": "GetUploadURL failed",
			"error":   err,
		})
		return
	}

	// Use the download link to store the file in object storage
	req, _ := http.NewRequest("PUT", uploadUrl, bytes.NewReader(content))
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"message": "file save failed",
			"error":   err,
		})
		elog.Error("file save failed", l.E(err))
		return
	}

	c.JSON(200, f)
}

func DetectMimeType(filename string, content []byte) string {
	ext := strings.ToLower(filepath.Ext(filename))
	// Prefer determining the file type by extension
	switch ext {
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".doc":
		return "application/msword"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".pptx":
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case ".ppt":
		return "application/vnd.ms-powerpoint"
	case ".txt":
		return "text/plain"
	case ".md":
		return "text/markdown"
	case ".csv":
		return "text/csv"
	case ".rtf":
		return "text/rtf"
	case ".xmind":
		return "application/vnd.xmind"
	}
	mime := mimetype.Detect(content).String()
	// Use the filetype library as a fallback for additional detection
	if kind, err := filetype.Match(content); err == nil && kind.MIME.Value != "" {
		return kind.MIME.Value
	}
	return mime
}

type Config struct {
	Signature string `json:"signature"`
	AppId     string `json:"appId"`
	Endpoint  string `json:"endpoint"`
	Token     string `json:"token"`
	UserUuid  string `json:"userUuid"`
}

func GetFileInfo(c *gin.Context) {
	fileGuid := c.Param("fileGuid")
	lang := c.Query("lang")
	token := middlewares.FindAccessToken(c)
	var userId int64
	// Treat expired tokens as anonymous users
	if token == "" {
		userId = consts.ANONYMOUS
	} else {
		err := middlewares.ValidateUserToken(c, token)
		if err != nil {
			userId = consts.ANONYMOUS
		}
	}

	if userId != consts.ANONYMOUS {
		userId = getUserIdFromToken(c)
		if userId == 0 {
			userId = consts.ANONYMOUS
		}
	}

	file, err := db.FindFileByGuidAndUserId(invoker.DB, userId, fileGuid)
	if err != nil {
		handleDBError(c, err)
		return
	}

	// TODO revisit the form-handling logic
	if file.ShimoType == "form" {
		formEditable := true
		if userId < 0 {
			formEditable = false
		} else {
			formEditable = file.Permissions["editable"]
		}
		permissions := map[string]bool{
			"formFillable": true,
			"readable":     true,
			"editable":     formEditable,
		}
		file.Permissions = permissions
	}

	mode, ok := c.GetQuery("mode")
	if !ok {
		c.JSON(400, gin.H{
			"message": "mode is required",
		})
	}

	returnConnectConfig := false
	switch mode {
	case "shimo":
		if file.IsShimoFile != 1 {
			c.JSON(http.StatusOK, file)
			return
		}
		returnConnectConfig = true
	case "preview":
		auth := utils.GetAuth(getUserIdFromToken(c))
		params := sdkapi.CreatePreviewParams{
			Auth:   auth,
			FileId: fileGuid,
		}
		// Create the preview
		r, e := invoker.SdkMgr.CreatePreview(params)
		if e != nil {
			handleSdkMgrError(c, r.Resp.Body(), r.Resp.StatusCode())
			return
		}
		// Return the preview URL
		resp := struct {
			db.File
			PreviewUrl string `json:"previewUrl"`
		}{
			File:       *file,
			PreviewUrl: genPreviewUrl(fileGuid, getUserIdFromToken(c), lang),
		}

		c.JSON(http.StatusOK, resp)
	case "form_fill":
		if file.IsShimoFile != 1 {
			c.JSON(http.StatusNotFound, nil)
			return
		}
		if file.ShimoType != "form" {
			c.JSON(http.StatusNotFound, nil)
			return
		}
		file.Permissions["formFillable"] = true
		returnConnectConfig = true
	}

	configToken := utils.SignUserJWTWithMode(userId, mode)
	if userId < 0 {
		configToken = consts.ANONYMOUSTOKEN
	}

	if returnConnectConfig {
		resp := struct {
			db.File
			Config `json:"config"`
		}{
			File: *file,
			Config: Config{
				Signature: invoker.Services.SignatureService.Sign(econf.GetString("shimoSDK.appId"), econf.GetString("shimoSDK.appSecret"), false),
				Endpoint:  econf.GetString("shimoSDK.host"),
				Token:     configToken,
				UserUuid:  utils.GetHashUserUuid(userId),
			},
		}
		c.JSON(200, resp)
	}
}

func ImportFile(c *gin.Context) {
	_file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"message": "file update failed"})
		return
	}
	// File name
	name, _ := c.GetPostForm("name")
	originalName := name
	// File type
	shimoType, _ := c.GetPostForm("shimoType")
	if name == "" {
		name = "untitled"
	} else {
		name = strings.TrimSuffix(name, filepath.Ext(name))
	}
	file := db.File{
		Name:        name,
		ShimoType:   shimoType,
		IsShimoFile: IsShimoType(shimoType),
		CreatorId:   getUserIdFromToken(c),
	}
	err, _ = db.CreateFile(invoker.DB, &file, getUserIdFromToken(c))
	if err != nil {
		handleDBError(c, err)
		return
	}
	auth := utils.GetAuth(getUserIdFromToken(c))
	// Convert multipart.FileHeader to os.File
	f, err := utils.ConvertToOSFile(_file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "file convert failed" + err.Error()})
		return
	}
	body := sdkapi.ImportFileReqBody{
		FileId:   file.Guid,
		Type:     file.ShimoType,
		File:     f,
		FileName: originalName,
	}
	params := sdkapi.ImportFileParams{
		Auth:              auth,
		ImportFileReqBody: body,
	}
	// Upload the file through the SDK
	res, err := invoker.SdkMgr.ImportFile(params)
	if err != nil || res.Status != 0 {
		// Roll back the created file
		rmErr := db.RemoveFileById(invoker.DB, file.ID)
		if rmErr != nil {
			elog.Warn("rollback file failed", l.E(rmErr))
		}
		handleSdkMgrError(c, res.Resp.Body(), res.Resp.StatusCode())
		return
	}
	taskId := res.Data.TaskId
	if taskId == "" {
		c.JSON(500, "taskId not found")
		return
	}
	// Return the result
	resp := struct {
		db.File
		TaskId string `json:"taskId"`
	}{
		File:   file,
		TaskId: taskId,
	}
	c.JSON(200, resp)
}

func ImportFileByUrl(c *gin.Context) {
	name, _ := c.GetPostForm("fileName")
	fileUrl, _ := c.GetPostForm("fileUrl")
	shimoType, _ := c.GetPostForm("shimoType")
	// Create the source file
	file := db.File{
		Name:        name,
		ShimoType:   shimoType,
		IsShimoFile: IsShimoType(shimoType),
		CreatorId:   getUserIdFromToken(c),
	}
	err, _ := db.CreateFile(invoker.DB, &file, getUserIdFromToken(c))
	if err != nil {
		handleDBError(c, err)
		return
	}
	// Create the SDK file
	auth := utils.GetAuth(getUserIdFromToken(c))
	body := sdkapi.ImportFileReqBody{
		FileId:   file.Guid,
		Type:     file.ShimoType,
		FileName: file.Name,
		FileUrl:  fileUrl,
	}
	params := sdkapi.ImportFileParams{
		Auth:              auth,
		ImportFileReqBody: body,
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
		handleSdkMgrError(c, res.Resp.Body(), res.Resp.StatusCode())
		return
	}

	taskId := res.Data.TaskId
	if taskId == "" {
		c.JSON(500, "taskId not found")
		return
	}
	// Return the file information
	resp := struct {
		db.File
		TaskId string `json:"taskId"`
	}{
		File:   file,
		TaskId: taskId,
	}
	c.JSON(200, resp)
}

func CheckImportProgress(c *gin.Context) {
	taskId, ok := c.GetQuery("taskId")
	if !ok {
		c.JSON(400, gin.H{"message": "taskId is required"})
		return
	}
	auth := utils.GetAuth(getUserIdFromToken(c))
	params := sdkapi.GetImportProgParams{
		Auth:   auth,
		TaskId: taskId,
	}
	// Get the upload progress
	resp, err := invoker.SdkMgr.GetImportProgress(params)
	if err != nil {
		fileId, ok := c.GetQuery("fileId")
		if !ok {
			elog.Warn("fileId is nil rollback file failed")
		}
		// Roll back the created file
		rmErr := db.RemoveFileByGuid(invoker.DB, fileId)
		if rmErr != nil {
			elog.Warn("rollback file failed", l.E(rmErr))
		}
		handleSdkMgrError(c, resp.Resp.Body(), resp.Resp.StatusCode())
		return
	}
	c.JSON(200, resp)
}

func CheckImportUrlProgress(c *gin.Context) {
	taskId, ok := c.GetQuery("taskId")
	if !ok {
		c.JSON(400, gin.H{"message": "taskId is required"})
		return
	}
	auth := utils.GetAuth(getUserIdFromToken(c))
	params := sdkapi.GetImportProgParams{
		Auth:   auth,
		TaskId: taskId,
	}
	var resp sdkapi.GetImportProgRespBody
	var err error
	if econf.GetString("shimoSDK.importByUrlVersion") == "v2" {
		resp, err = invoker.SdkMgr.GetImportV2Progress(params)
	} else {
		resp, err = invoker.SdkMgr.GetImportProgress(params)
	}
	if err != nil {
		fileId, ok := c.GetQuery("fileId")
		if !ok {
			elog.Warn("fileId is nil rollback file failed")
		}
		// Roll back the created file
		rmErr := db.RemoveFileByGuid(invoker.DB, fileId)
		if rmErr != nil {
			elog.Warn("rollback file failed", l.E(rmErr))
		}
		handleSdkMgrError(c, resp.Resp.Body(), resp.Resp.StatusCode())
		return
	}
	c.JSON(200, resp)
}

func ExportFile(c *gin.Context) {
	fileGuid := c.Param("fileGuid")

	_, err := db.FindFileByGuid(invoker.DB, fileGuid)
	if err != nil {
		handleDBError(c, err)
		return
	}

	exportType, ok := c.GetQuery("type")
	if !ok {
		c.JSON(400, gin.H{"message": "type is required"})
		return
	}
	auth := utils.GetAuth(getUserIdFromToken(c))
	params := sdkapi.ExportFileParams{
		Auth:   auth,
		FileId: fileGuid,
		Type:   exportType,
	}
	res, err := invoker.SdkMgr.ExportFile(params)
	if err != nil {
		handleSdkMgrError(c, res.Resp.Body(), res.Resp.StatusCode())
		return
	}
	c.JSON(200, res)
}

func DuplicateFile(c *gin.Context) {
	fileGuid := c.Param("fileGuid")
	userId := getUserIdFromToken(c)

	file, err := db.FindFileByGuid(invoker.DB, fileGuid)
	if err != nil {
		handleDBError(c, err)
		return
	}
	newFile := db.File{
		Name:        file.Name + " copy",
		Type:        file.Type,
		ShimoType:   file.ShimoType,
		IsShimoFile: file.IsShimoFile,
		CreatorId:   userId,
		Role:        "owner",
	}
	err, _ = db.CreateFile(invoker.DB, &newFile, userId)
	if err != nil {
		handleDBError(c, err)
		return
	}
	auth := utils.GetAuth(getUserIdFromToken(c))
	params := sdkapi.CreateFileCopyParams{
		Auth:         auth,
		OriginFileId: fileGuid,
		TargetFileId: newFile.Guid,
	}
	res, err := invoker.SdkMgr.CreateFileCopy(params)
	if err != nil {
		rmErr := db.RemoveFileById(invoker.DB, newFile.ID)
		if rmErr != nil {
			elog.Warn("rollback file failed", l.E(rmErr))
		}
		handleSdkMgrError(c, res.Resp.Body(), res.Resp.StatusCode())
		return
	}
	c.JSON(200, newFile)
}

func CheckExportFileProgress(c *gin.Context) {
	taskId, ok := c.GetQuery("taskId")
	if !ok {
		c.JSON(400, gin.H{"message": "taskId is required"})
		return
	}
	auth := utils.GetAuth(getUserIdFromToken(c))
	params := sdkapi.GetExportProgParams{
		Auth:   auth,
		TaskId: taskId,
	}
	res, err := invoker.SdkMgr.GetExportProgress(params)
	if err != nil {
		handleSdkMgrError(c, res.Resp.Body(), res.Resp.StatusCode())
		return
	}
	c.JSON(200, res)
}

func DeleteFile(c *gin.Context) {
	fileGuid := c.Param("fileGuid")
	file, err := db.FindFileByGuid(invoker.DB, fileGuid)
	if err != nil {
		handleDBError(c, err)
		return
	}
	if file.IsShimoFile == 1 {
		auth := utils.GetAuth(getUserIdFromToken(c))
		params := sdkapi.DeleteFileParams{
			Auth:   auth,
			FileId: fileGuid,
		}
		_, rErr := invoker.SdkMgr.DeleteFile(params)
		if rErr != nil {
			elog.Warn("file remove failed", l.E(rErr))
		}
	} else {
		rErr := invoker.Services.AwosService.Remove(file.Guid)
		if rErr != nil {
			elog.Warn("file remove failed", l.E(rErr))
		}
	}
	err = db.RemoveFileByGuid(invoker.DB, fileGuid)
	if err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(204, nil)
}

func BatchDeleteFile(c *gin.Context) {
	fileGuids := make([]string, 0)
	if err := c.BindJSON(&fileGuids); err != nil {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
			return
		}
	}
	err := db.RemoveFileByGuids(invoker.DB, fileGuids)
	if err != nil {
		handleDBError(c, err)
		return
	}
	// TODO batch deletion does not call the SDK delete API yet
	c.JSON(204, nil)
}

func RenameFile(c *gin.Context) {
	fileGuid := c.Param("fileGuid")

	body := FileBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"message": "request body error"})
		return
	}

	err := db.RenameFile(invoker.DB, fileGuid, body.Name)
	if err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(204, nil)
}

type Collaborators struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Permissions []map[string]bool `json:"permissions"`
	IsCreator   bool              `json:"isCreator"`
}

func GetCollaborators(c *gin.Context) {
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
	mp := map[int64][]map[string]bool{}
	for _, fp := range fps {
		mp[fp.UserId] = append(mp[fp.UserId], fp.Permissions)
	}

	userId := getUserIdFromToken(c)
	me, err := db.FindUserById(invoker.DB, userId)
	if err != nil {
		handleDBError(c, err)
		return
	}

	users, err := db.FindUsersByAppId(invoker.DB, me.AppID)
	if err != nil {
		handleDBError(c, err)
		return
	}

	returnAll, _ := c.GetQuery("all")
	res := make([]Collaborators, 0, len(users))

	for _, user := range users {
		if user.ID == me.ID {
			continue
		}
		permissions := mp[user.ID]
		if len(permissions) > 0 {
			isCreator := false
			if user.ID == file.CreatorId {
				isCreator = true
			}
			res = append(res, Collaborators{
				Id:          strconv.FormatInt(user.ID, 10),
				Name:        user.Name,
				Permissions: permissions,
				IsCreator:   isCreator,
			})
		} else if returnAll == "true" || returnAll == "1" {
			res = append(res, Collaborators{
				Id:          strconv.FormatInt(user.ID, 10),
				Name:        user.Name,
				Permissions: []map[string]bool{},
				IsCreator:   false,
			})
		}
	}
	c.JSON(200, res)
}

func UpdateCollaborators(c *gin.Context) {
	fileGuid := c.Param("fileGuid")

	myId := getUserIdFromToken(c)
	file, err := db.FindFileByGuidAndUserId(invoker.DB, myId, fileGuid)
	if err != nil {
		handleDBError(c, err)
		return

	}

	mp := map[int64]map[string]bool{}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		fps, _ := db.FindFilePermissionsByFileId(invoker.DB, file.ID)
		for _, fp := range fps {
			mp[fp.UserId] = fp.Permissions
		}
		wg.Done()
	}()

	var body map[int64]map[string]bool
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if file.Permissions["manageable"] == false {
		c.JSON(403, gin.H{"message": "requires manageable permission"})
		return
	}

	wg.Wait()
	for userId, p := range body {
		checkPermission(p)

		curUserP, ok := mp[userId]
		if ok {
			if curUserP["manageable"] && userId == file.CreatorId {
				continue
			}
			mergePermission(curUserP, p)

			mp[userId] = curUserP
		} else {
			mergeNewPermission := map[string]bool{}
			mergePermission(mergeNewPermission, p)
			mp[userId] = mergeNewPermission
		}

	}

	var row []db.FilePermissions
	for userId, p := range mp {
		r := "collaborator"
		if p["manageable"] {
			r = "owner"
		}
		row = append(row, db.FilePermissions{
			FileId:      file.ID,
			UserId:      userId,
			Role:        r,
			Permissions: p,
		})
		if r == "collaborator" {
			// Check whether all permissions have been removed
			allCanceled := true
			for _, v := range p {
				if v {
					allCanceled = false
					break
				}
			}
			if allCanceled {
				// Remove the user's permission record
				_ = db.RemovePermissionsByFileId(invoker.DB, file.ID, userId)
				continue
			}
		}
	}

	err = db.BatchSavePermissions(invoker.DB, row)
	if err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(204, nil)
}

func ExportTableSheets(c *gin.Context) {
	fileGuid := c.Param("fileGuid")

	auth := utils.GetAuth(getUserIdFromToken(c))
	params := sdkapi.ExportTableSheetsParams{
		Auth:   auth,
		FileId: fileGuid,
	}
	res, err := invoker.SdkMgr.ExportTableSheets(params)
	if err != nil {
		handleSdkMgrError(c, res.Resp.Body(), res.Resp.StatusCode())
		return
	}
	c.JSON(200, gin.H{
		"downloadUrl": res.DownloadUrl,
	})
}

func mergePermission(base, src map[string]bool) {
	if src == nil {
		return
	}
	for _, initK := range consts.InitFilePermission() {
		if src[string(initK)] || src["manageable"] {
			base[string(initK)] = true
		} else {
			base[string(initK)] = false
		}
	}
	for _, newK := range consts.NewFilePermission() {
		if econf.GetBool("permissions.setNewFilePermission") {
			if src[string(newK)] || src["manageable"] {
				base[string(newK)] = true
			} else {
				base[string(newK)] = false
			}
		} else {
			delete(base, string(newK))
		}
	}
}

func checkPermission(p map[string]bool) {
	if p["manageable"] || p["exportable"] || p["copyable"] || p["editable"] || p["commentable"] || p["lockable"] || p["unlockable"] || p["attachmentCopyable"] || p["attachmentPreviewable"] || p["attachmentDownloadable"] || p["copyablePasteClipboard"] || p["cutable"] || p["imageDownloadable"] {
		p["readable"] = true
	}
}

func IsShimoType(typ string) int {
	if typ == "" {
		return 0
	}
	return 1
}

func genPreviewUrl(fileGuid string, userId int64, lang string) string {
	previewUrl := econf.GetString("shimoSDK.host") + fmt.Sprintf("/api/cloud-files/%s/page", fileGuid)
	appId := econf.GetString("shimoSDK.appId")
	parseUrl, err := url.Parse(previewUrl)
	if err != nil {
		return ""
	}

	queryParams := parseUrl.Query()
	queryParams.Add("lang", lang)
	queryParams.Add("appId", appId)
	queryParams.Add("token", utils.SignUserJWT(userId))
	queryParams.Add("signature", utils.Sign(appId, econf.GetString("shimoSDK.appSecret"), false))

	parseUrl.RawQuery = queryParams.Encode()
	return parseUrl.String()
}

func genInspectPreviewUrl(fileGuid string, userId int64, lang string) string {
	previewUrl := econf.GetString("shimoSDK.host") + fmt.Sprintf("/api/cloud-files/%s/page", fileGuid)
	appId := econf.GetString("shimoSDK.appId")
	parseUrl, err := url.Parse(previewUrl)
	if err != nil {
		return ""
	}

	queryParams := parseUrl.Query()
	queryParams.Add("lang", lang)
	queryParams.Add("appId", appId)
	queryParams.Add("token", utils.SignUserJWT(userId, 365*24*time.Hour))
	queryParams.Add("signature", utils.Sign(appId, econf.GetString("shimoSDK.appSecret"), false))

	parseUrl.RawQuery = queryParams.Encode()
	return parseUrl.String()
}

func FormatCurrentTime() string {
	now := time.Now()
	return now.Format("2006-01-02 15:04:05")
}

func GetImportUrl0(c *gin.Context) {
	isDev := econf.GetBool("import.isDev")
	if !isDev {
		return
	}
	userId := getUserIdFromToken(c)
	body := struct {
		DownloadUrl string `json:"downloadUrl"`
	}{}
	_ = c.BindJSON(&body)
	fileURL := body.DownloadUrl
	parsedURL, err := url.Parse(fileURL)
	filePath := parsedURL.Path
	fileName := path.Base(filePath)

	// Persist to the database
	file := db.File{
		Name:        fileName,
		ShimoType:   "spreadsheet",
		IsShimoFile: IsShimoType("spreadsheet"),
		CreatorId:   getUserIdFromToken(c),
	}

	err, _ = db.CreateFile(invoker.DB, &file, getUserIdFromToken(c))
	if err != nil {
		handleDBError(c, err)
		return
	}

	// Import via URL
	auth := utils.GetAuth(getUserIdFromToken(c))
	reqBody := sdkapi.ImportFileReqBody{
		FileId:   file.Guid,
		Type:     "spreadsheet",
		FileName: fileName,
		FileUrl:  fileURL,
	}
	params := sdkapi.ImportFileParams{
		Auth:              auth,
		ImportFileReqBody: reqBody,
	}
	ImportResp, err := invoker.SdkMgr.ImportFile(params)
	if err != nil {
		handleSdkMgrError(c, ImportResp.Resp.Body(), ImportResp.Resp.StatusCode())
		return
	}

	// Fetch the upload progress
	var progressResp sdkapi.GetImportProgRespBody
	success := false
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	timeout := time.After(20 * time.Second)
	for {
		select {
		case <-ticker.C:
			progressParams := sdkapi.GetImportProgParams{
				Auth:   auth,
				TaskId: ImportResp.Data.TaskId,
			}
			// Fetch the upload progress
			if progressResp, _ = invoker.SdkMgr.GetImportProgress(progressParams); progressResp.Status == 0 {
				success = true
				break
			}
		case <-timeout:
			handleDBError(c, errors.New("import file timeout"))
			return
		}
		if success {
			break
		}
	}

	// After a successful upload, create a preview and fetch its URL
	createParams := sdkapi.CreatePreviewParams{
		Auth:   auth,
		FileId: file.Guid,
	}
	// Create the preview
	previewResp, err := invoker.SdkMgr.CreatePreview(createParams)
	if err != nil || previewResp.Code != "" {
		rmErr := db.RemoveFileByGuid(invoker.DB, file.Guid)
		if rmErr != nil {
			elog.Warn("rollback file failed", l.E(rmErr))
		}
		handleSdkMgrError(c, previewResp.Resp.Body(), previewResp.Resp.StatusCode())
		return
	}

	confEnv := os.Getenv("EGO_CONFIG_PATH")
	// Get the path part of the request (excluding query parameters)
	fullPath := c.Request.URL.Path
	var redirectURL string
	// For other environments
	if confEnv != "" {
		redirectURL = econf.GetString("host.addr") + fullPath + "/redirect"
	} else {
		// Build the base URL without query parameters
		baseURL := "http" + "://" + c.Request.Host + fullPath
		// Append the "/redirect" path to the URL
		redirectURL = baseURL + "/redirect"
	}

	// Build the query parameters
	queryParams := url.Values{}
	queryParams.Add("fileId", file.Guid)
	queryParams.Add("userId", strconv.FormatInt(userId, 10))
	// Attach the query parameters to the URL
	redirectURLWithParams := redirectURL + "?" + queryParams.Encode()
	// Respond with the constructed address
	c.JSON(http.StatusOK, gin.H{
		"url": redirectURLWithParams,
	})
}

func GetImportRedirectUrl1(c *gin.Context) {
	// Retrieve the query parameters
	fileId := c.Query("fileId")
	userId := c.Query("userId")
	uId, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
	}
	previewUrl := genPreviewUrl(fileId, uId, "")

	c.Redirect(http.StatusFound, previewUrl)
}

func GetImportUrl(c *gin.Context) {
	isDev := econf.GetBool("import.isDev")
	if !isDev {
		return
	}
	userId := getUserIdFromToken(c)
	body := struct {
		DownloadUrl string `json:"downloadUrl"`
	}{}
	_ = c.BindJSON(&body)
	fileURL := body.DownloadUrl
	parsedURL, err := url.Parse(fileURL)
	filePath := parsedURL.Path
	fileName := path.Base(filePath)

	// Persist to the database
	file := db.File{
		Name:        fileName,
		ShimoType:   "spreadsheet",
		IsShimoFile: IsShimoType("spreadsheet"),
		CreatorId:   getUserIdFromToken(c),
	}

	err, _ = db.CreateFile(invoker.DB, &file, getUserIdFromToken(c))
	if err != nil {
		handleDBError(c, err)
		return
	}

	auth := utils.GetAuth(getUserIdFromToken(c))
	// Import via URL
	reqBody := sdkapi.ImportFileReqBody{
		FileId:   file.Guid,
		Type:     "spreadsheet",
		FileName: fileName,
		FileUrl:  fileURL,
	}
	params := sdkapi.ImportFileParams{
		Auth:              auth,
		ImportFileReqBody: reqBody,
	}
	ImportResp, err := invoker.SdkMgr.ImportFile(params)
	if err != nil {
		handleSdkMgrError(c, ImportResp.Resp.Body(), ImportResp.Resp.StatusCode())
		return
	}

	// Fetch the upload progress
	var progressResp sdkapi.GetImportProgRespBody
	success := false
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	timeout := time.After(20 * time.Second)
	for {
		select {
		case <-ticker.C:
			progressParams := sdkapi.GetImportProgParams{
				Auth:   auth,
				TaskId: ImportResp.Data.TaskId,
			}
			// Fetch the upload progress
			if progressResp, _ = invoker.SdkMgr.GetImportProgress(progressParams); progressResp.Status == 0 {
				success = true
				break
			}
		case <-timeout:
			handleDBError(c, errors.New("import file timeout"))
			return
		}
		if success {
			break
		}
	}

	confEnv := os.Getenv("EGO_CONFIG_PATH")
	// Get the path part of the request (excluding query parameters)
	fullPath := c.Request.URL.Path
	var redirectURL string
	// For other environments
	if confEnv != "" {
		redirectURL = econf.GetString("host.addr") + fullPath + "/redirect"
	} else {
		// Build the base URL without query parameters
		baseURL := "http" + "://" + c.Request.Host + fullPath
		// Append the "/redirect" path to the URL
		redirectURL = baseURL + "/redirect"
	}

	// Build the query parameters
	queryParams := url.Values{}
	queryParams.Add("fileId", file.Guid)
	queryParams.Add("userId", strconv.FormatInt(userId, 10))
	// Attach the query parameters to the URL
	redirectURLWithParams := redirectURL + "?" + queryParams.Encode()
	// Respond with the constructed address
	c.JSON(http.StatusOK, gin.H{
		"url": redirectURLWithParams,
	})
}

func GetImportRedirectUrl(c *gin.Context) {
	// Retrieve the query parameters
	fileId := c.Query("fileId")
	userId := c.Query("userId")
	uId, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
	}

	token := utils.SignUserJWT(uId)

	var fullPath string
	queryParams := url.Values{}
	queryParams.Add("accessToken", token)
	confEnv := os.Getenv("EGO_CONFIG_PATH")
	if confEnv != "" {
		fullPath = econf.GetString("host.addr") + "/shimo-files/" + fileId + "?" + queryParams.Encode()
	} else {
		// Build the base URL without query parameters
		fullPath = "http" + "://" + "localhost:8000" + "/sdk/demo/shimo-files/" + fileId + "?" + queryParams.Encode()
	}

	c.Redirect(http.StatusFound, fullPath)
}

type TempFile struct {
	Name string `json:"name"`
}
type FileUrl struct {
	Preview     string `json:"preview"`
	Collaborate string `json:"collaborate"`
}

func FrontInspectCreate(c *gin.Context) {
	var body map[string]TempFile
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	userId := getUserIdFromToken(c)

	resultMap := make(map[consts.FileType]FileUrl)

	for key, value := range body {
		fileType := consts.GetFileType(key)
		if fileType == consts.FileTypeInvalid {
			continue
		}
		// Generate the file URL
		fileUrl, err := handleFileCreation(userId, value, fileType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		resultMap[fileType] = fileUrl
	}

	c.JSON(http.StatusOK, resultMap)
}

// Extract the file-creation logic into a separate function
func handleFileCreation(userId int64, value TempFile, fileType consts.FileType) (FileUrl, error) {
	guid := utils.GenerateUserFileUUID(strconv.FormatInt(userId, 10), string(fileType))
	// Check whether the file exists
	if f, err := db.FindFileByGuid(invoker.DB, guid); err == nil && f.Guid != "" {
		return FileUrl{
			Preview:     genInspectPreviewUrl(f.Guid, userId, ""),
			Collaborate: fmt.Sprintf("%s/shimo-files/%s?accessToken=%s", econf.GetString("host.addr"), f.Guid, utils.SignUserJWT(userId, 365*24*time.Hour)),
		}, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return FileUrl{}, fmt.Errorf("%s failed to find file: %v", fileType, err)
	}
	// Create a new file
	file := db.File{
		Guid:        guid,
		Name:        value.Name,
		ShimoType:   string(fileType),
		CreatorId:   userId,
		IsShimoFile: IsShimoType(string(fileType)),
	}

	err, fileId := db.CreateFile(invoker.DB, &file, userId)
	if err != nil {
		return FileUrl{}, fmt.Errorf("%s failed to create file: %v", fileType, err)
	}
	// Call the SDK to create the file
	auth := utils.GetAuth(userId)
	cFile := sdkapi.CreateFileParams{
		FileType: sdkcommon.CollabFileType(file.ShimoType),
		Auth:     auth,
		FileId:   file.Guid,
	}

	if _, err := invoker.SdkMgr.CreateFile(cFile); err != nil {
		_ = db.RemoveFileById(invoker.DB, fileId)
		return FileUrl{}, fmt.Errorf("SDK failed to create file: %v", err)
	}

	return FileUrl{
		Preview:     genInspectPreviewUrl(file.Guid, userId, ""),
		Collaborate: fmt.Sprintf("%s/shimo-files/%s?accessToken=%s", econf.GetString("host.addr"), file.Guid, utils.SignUserJWT(userId, 365*24*time.Hour)),
	}, nil
}

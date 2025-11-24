package sdkctl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/pkg/errors"
	sdkapi "github.com/shimo-open/sdk-kit-go/model/api"
	sdkcommon "github.com/shimo-open/sdk-kit-go/model/common"

	"sdk-demo-go/pkg/consts"
	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/models/db"
	"sdk-demo-go/pkg/utils"
)

// TestCreatePreview creates a file preview
func TestCreatePreview(ctx context.Context, fileId string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiCreatePreview
	printStartTestMsg(apiName.Name())
	auth := utils.GetAuth(userId)
	params := sdkapi.CreatePreviewParams{
		Auth:   auth,
		FileId: fileId,
	}
	// Create the preview
	resp, err := invoker.SdkMgr.CreatePreview(params)
	if err != nil {
		mgrErrHandler(err, string(resp.Resp.Body()))
	}
	testRes = testMgrResHandler(apiName, resp, resp.RawResponse, err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestGetPreview TODO open the file preview (consider chromeless)
func TestGetPreview(ctx context.Context, fileId string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiGetPreview
	printStartTestMsg(apiName.Name())
	params := sdkapi.CreatePreviewParams{
		Auth:   utils.GetAuth(userId),
		FileId: fileId,
	}
	// Create the preview
	resp, err := invoker.SdkMgr.CreatePreview(params)
	if err != nil {
		errHandler(err)
	}

	testRes = testMgrResHandler(apiName, resp, resp.RawResponse, err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

func isShimoType(typ string) int {
	if typ == "" {
		return 0
	}
	return 1
}

// TestCreateFile creates a collaborative document
func TestCreateFile(ctx context.Context, fileName string, fileType string, shimoType string, lang string, acLang string) (fileGuid string, testRes consts.SingleApiTestRes, err error) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiCreateFile
	printStartTestMsg(apiName.Name())
	if fileType == "" {
		fileType = "document"
	}
	if fileName == "" {
		fileName = fmt.Sprintf(`%s - %s`, time.Now().Format("2006-01-02 15:04:05"), fileType)
	}

	shimoType = fileType

	userId = getUserId()
	file := db.File{
		Name:        fileName,
		Type:        fileType,
		ShimoType:   shimoType,
		CreatorId:   userId,
		IsShimoFile: isShimoType(shimoType),
	}
	err, _ = db.CreateFile(invoker.DB, &file, userId)
	if err != nil {
		errHandler(err)
		return
	}

	fileGuid = file.Guid
	auth := utils.GetAuth(userId)
	cFile := sdkapi.CreateFileParams{
		Auth:     auth,
		FileType: sdkcommon.CollabFileType(file.ShimoType),
		Lang:     sdkcommon.Lang(lang),
		FileId:   file.Guid,
	}
	resp, err := invoker.SdkMgr.CreateFile(cFile)
	if err != nil {
		// Remove the record if creation fails
		_ = db.RemoveFileById(invoker.DB, file.ID)
		mgrErrHandler(err, string(resp.Resp.Body()))
	}
	testRes = testMgrResHandler(apiName, resp, resp, err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestCreateFileCopy creates a copy of a collaborative document
func TestCreateFileCopy(ctx context.Context, fileId string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiCreateFileCopy
	printStartTestMsg(apiName.Name())
	userId := getUserId()

	// Find the file to copy
	file, err := db.FindFileByGuid(invoker.DB, fileId)
	if err != nil {
		errHandler(err)
		return
	}

	// Create a new file
	newFile := db.File{
		Name:        file.Name + " copy",
		Type:        file.Type,
		ShimoType:   file.ShimoType,
		CreatorId:   userId,
		IsShimoFile: file.IsShimoFile,
	}

	// Persist the new file to the database
	err, _ = db.CreateFile(invoker.DB, &newFile, userId)
	if err != nil {
		errHandler(err)
		return
	}
	auth := utils.GetAuth(userId)
	params := sdkapi.CreateFileCopyParams{
		Auth:         auth,
		OriginFileId: fileId,
		TargetFileId: newFile.Guid,
	}
	resp, err := invoker.SdkMgr.CreateFileCopy(params)

	// Call the API and capture the result
	mgrErrHandler(err, string(resp.Resp.Body()))

	testRes = testMgrResHandler(apiName, resp, resp, err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestDeleteFile deletes a collaborative document
func TestDeleteFile(ctx context.Context, fileId string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiDeleteFile
	printStartTestMsg(apiName.Name())
	file, err := db.FindFileByGuid(invoker.DB, fileId)
	if err != nil {
		errHandler(err)
		return
	}

	var resp sdkapi.RawResponse
	if file.IsShimoFile == 1 {
		auth := utils.GetAuth(userId)
		params := sdkapi.DeleteFileParams{
			Auth:   auth,
			FileId: fileId,
		}
		resp, err = invoker.SdkMgr.DeleteFile(params)
		mgrErrHandler(err, string(resp.Resp.Body()))
	} else {
		err = invoker.Services.AwosService.Remove(file.Guid)
		if err != nil {
			errHandler(err)
			fmt.Println(fmt.Sprintf("awos remove file failed, guid: %s, err:%e", file.Guid, err))
		}
	}

	_ = db.RemoveFileByGuid(invoker.DB, fileId)

	testRes = testMgrResHandler(apiName, resp, resp, err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestImportFile imports a file
func TestImportFile(ctx context.Context, filePath string, fileName string, shimoType string) (taskId, fileGuid string, statusCode int, err error) {
	apiName := consts.ShimoSdkApiImportFile
	printStartTestMsg(apiName.Name())

	file := db.File{
		Name:        fileName,
		ShimoType:   shimoType,
		IsShimoFile: isShimoType(shimoType),
		CreatorId:   getUserId(),
	}

	err, _ = db.CreateFile(invoker.DB, &file, getUserId())
	if err != nil {
		errHandler(err)
		return
	}
	// Open the file
	f, err := os.Open(filePath)
	if err != nil {
		// Roll back the newly created file
		_ = db.RemoveFileById(invoker.DB, file.ID)
		errHandler(errors.Wrap(err, "open file failed"))
		return
	}
	body := sdkapi.ImportFileReqBody{
		FileId:   file.Guid,
		Type:     file.ShimoType,
		File:     f,
		FileName: fileName,
	}
	params := sdkapi.ImportFileParams{
		Auth:              utils.GetAuth(userId),
		ImportFileReqBody: body,
	}
	// Upload the file via the SDK
	resp, err := invoker.SdkMgr.ImportFile(params)
	if err != nil {
		fmt.Println("invoker.SdkMgr.ImportFile err:" + err.Error())
		// Roll back the newly created file
		_ = db.RemoveFileById(invoker.DB, file.ID)
		mgrErrHandler(err, string(resp.Resp.Body()))
	}
	return resp.Data.TaskId, file.Guid, resp.Resp.StatusCode(), err
}

// TestImportFileByUrl imports a file via URL
func TestImportFileByUrl(ctx context.Context, url string, fileName string, fileType string, shimoType string) {
	printStartTestMsg("url导入")
	file := db.File{
		Name:        fileName,
		Type:        fileType,
		ShimoType:   shimoType,
		IsShimoFile: isShimoType(shimoType),
		CreatorId:   getUserId(),
	}

	err, _ := db.CreateFile(invoker.DB, &file, getUserId())
	if err != nil {
		errHandler(err)
		return
	}

	body := sdkapi.ImportFileReqBody{
		FileId:   file.Guid,
		Type:     file.ShimoType,
		FileUrl:  url,
		FileName: file.Name,
	}
	params := sdkapi.ImportFileParams{
		Auth:              utils.GetAuth(userId),
		ImportFileReqBody: body,
	}
	// Upload the file via the SDK
	resp, err := invoker.SdkMgr.ImportFile(params)
	if err != nil {
		_ = db.RemoveFileByGuid(invoker.DB, file.Guid)
		mgrErrHandler(err, string(resp.Resp.Body()))
	}
	return
}

// TestImportFileProgress checks import progress
func TestImportFileProgress(ctx context.Context, taskId string) (importSuccess bool, statusCode int, err error) {
	apiName := consts.ShimoSdkApiImportFileProgress
	printStartTestMsg(apiName.Name())
	auth := utils.GetAuth(userId)
	params := sdkapi.GetImportProgParams{
		Auth:   auth,
		TaskId: taskId,
	}
	// Fetch the upload progress
	resp, err := invoker.SdkMgr.GetImportProgress(params)
	mgrErrHandler(err, string(resp.Resp.Body()))
	return resp.Data.Progress == 100, resp.Resp.StatusCode(), err
}

// TestExportFile exports a file
func TestExportFile(ctx context.Context, fileId string, exportType string) (taskId string) {
	printStartTestMsg("导出文件")
	auth := utils.GetAuth(userId)
	params := sdkapi.ExportFileParams{
		Auth:   auth,
		FileId: fileId,
		Type:   exportType,
	}
	resp, err := invoker.SdkMgr.ExportFile(params)
	mgrErrHandler(err, string(resp.Resp.Body()))

	return resp.Data.TaskId
}

// TestExportFileProgress checks export progress
func TestExportFileProgress(ctx context.Context, taskId string) (exportSuccess bool, statusCode int, err error) {
	printStartTestMsg("导出进度")
	auth := utils.GetAuth(userId)
	params := sdkapi.GetExportProgParams{
		Auth:   auth,
		TaskId: taskId,
	}
	resp, err := invoker.SdkMgr.GetExportProgress(params)
	mgrErrHandler(err, string(resp.Resp.Body()))
	fmt.Print(fmt.Sprintf(" 进度：%d ", resp.Data.Progress))
	return resp.Data.Progress == 100, resp.Resp.StatusCode(), err
}

// TestExportTableAsSheets exports application tables to Excel
func TestExportTableAsSheets(ctx context.Context, fileId string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiExportTableAsExcel
	printStartTestMsg(apiName.Name())
	params := sdkapi.ExportTableSheetsParams{
		Auth:   utils.GetAuth(userId),
		FileId: fileId,
	}
	resp, err := invoker.SdkMgr.ExportTableSheets(params)
	mgrErrHandler(err, string(resp.Resp.Body()))
	testRes = testMgrResHandler(apiName, resp, resp.RawResponse, err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestGetFilePlainText retrieves the plain-text content of a document
func TestGetFilePlainText(ctx context.Context, fileId string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiGetFilePlainText
	printStartTestMsg(apiName.Name())
	// Define the response validator
	validators := []utils.Validator{
		{
			Check:   "content",
			Assert:  "not_equal",
			Expect:  "",
			Message: "content 字段应存在且不为空",
		},
	}
	// Invoke the API and collect the result
	params := sdkapi.GetPlainTextParams{
		FileId: fileId,
		Auth:   utils.GetAuth(userId),
	}
	resp, err := invoker.SdkMgr.GetPlainText(params)
	mgrErrHandler(err, string(resp.Resp.Body()))
	testRes = testMgrResHandler(apiName, resp, resp.RawResponse, err, time.Now().Sub(startTime).String(), startTime.Unix(), validators)
	return
}

// TestGetFilePlainTextWC retrieves the plain-text word count
func TestGetFilePlainTextWC(ctx context.Context, fileId string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiGetFilePlainTextWordCount
	printStartTestMsg(apiName.Name())
	validators := []utils.Validator{
		{
			Check:   "wordCount",
			Assert:  "not_equal",
			Expect:  "",
			Message: "wordCount 字段应存在且不为空",
		},
	}
	params := sdkapi.GetPlainTextWCParams{
		Auth:   utils.GetAuth(userId),
		FileId: fileId,
	}
	resp, err := invoker.SdkMgr.GetPlainTextWC(params)
	mgrErrHandler(err, string(resp.Resp.Body()))
	testRes = testMgrResHandler(apiName, resp, resp.RawResponse, err, time.Now().Sub(startTime).String(), startTime.Unix(), validators)
	return
}

// TestGetDocSidebarInfo fetches the history list
func TestGetDocSidebarInfo(ctx context.Context, fileId string, page int, size int) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiGetDocSidebarInfo
	printStartTestMsg(apiName.Name())
	validators := []utils.Validator{
		{
			Check:   "length(histories)",
			Assert:  "equal",
			Expect:  1,
			Message: "histories 字段记录为一条",
		},
		{
			Check:   "isLastPage",
			Assert:  "equal",
			Expect:  true,
			Message: "isLastPage 字段应为 true",
		},
		{
			Check:   "users",
			Assert:  "not_equal",
			Expect:  "",
			Message: "users 字段不能为空",
		},
	}
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	params := sdkapi.GetHistoryListParams{
		Auth:     utils.GetAuth(userId),
		FileId:   fileId,
		PageSize: size,
		Count:    (page - 1) * size,
	}
	resp, err := invoker.SdkMgr.GetHistoryList(params)
	mgrErrHandler(err, string(resp.Resp.Body()))
	testRes = testMgrResHandler(apiName, resp, resp.RawResponse, err, time.Now().Sub(startTime).String(), startTime.Unix(), validators)
	return
}

// TestGetRevision retrieves the revision list
func TestGetRevision(ctx context.Context, fileId string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiGetRevision
	printStartTestMsg(apiName.Name())
	params := sdkapi.GetRevisionListParams{
		Auth:   utils.GetAuth(userId),
		FileId: fileId,
	}
	resp, _, err := invoker.SdkMgr.GetRevisionList(params)
	mgrErrHandler(err, string(resp.Resp.Body()))
	testRes = testMgrResHandler(apiName, resp, resp, err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

func TestGetExcelContent(ctx context.Context, fileId string, rg string, validators []utils.Validator) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiGetExcelContent
	printStartTestMsg(apiName.Name())
	params := sdkapi.GetTableContentParams{
		Auth:   utils.GetAuth(userId),
		FileId: fileId,
		Rg:     rg,
	}
	resp, err := invoker.SdkMgr.GetTableContent(params)
	mgrErrHandler(err, string(resp.Resp.Body()))
	testRes = testMgrResHandler(apiName, resp, resp.RawResponse, err, time.Now().Sub(startTime).String(), startTime.Unix(), validators)
	return
}

// TestUpdateExcelContent updates spreadsheet content
func TestUpdateExcelContent(ctx context.Context, fileId string, rg string, values []string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiUpdateExcelContent
	printStartTestMsg(apiName.Name())
	realValues, err := utils.ParseStringToObject(values)
	if err != nil {
		errHandler(err)
		return
	}
	params := sdkapi.UpdateTableContentParams{
		Auth:   utils.GetAuth(userId),
		FileId: fileId,
		UpdateTableContentRequestBody: sdkapi.UpdateTableContentRequestBody{
			Rg: rg,
			Resource: struct {
				Values [][]interface{} `json:"values"`
			}{Values: realValues},
		},
	}
	resp, err := invoker.SdkMgr.UpdateTableContent(params)
	mgrErrHandler(err, string(resp.Resp.Body()))
	testRes = testMgrResHandler(apiName, resp, resp, err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestAppendExcelContent appends spreadsheet content
func TestAppendExcelContent(ctx context.Context, fileId string, rg string, values []string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiAppendExcelContent
	printStartTestMsg(apiName.Name())
	realValues, err := utils.ParseStringToObject(values)
	if err != nil {
		errHandler(err)
		return
	}
	params := sdkapi.AppendTableContentParams{
		Auth:   utils.GetAuth(userId),
		FileId: fileId,
		AppendTableContentReqBody: sdkapi.AppendTableContentReqBody{
			Rg: rg,
			Resource: sdkapi.Resource{
				Values: realValues,
			},
		},
	}
	resp, err := invoker.SdkMgr.AppendTableContent(params)
	mgrErrHandler(err, string(resp.Resp.Body()))
	testRes = testMgrResHandler(apiName, resp, resp, err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestDeleteExcelRows deletes spreadsheet rows
func TestDeleteExcelRows(ctx context.Context, fileId string, sheetName string, index int) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiDeleteExcelRows
	printStartTestMsg(apiName.Name())
	params := sdkapi.DeleteTableRowParams{
		Auth:      utils.GetAuth(userId),
		FileId:    fileId,
		SheetName: sheetName,
		Index:     index,
		Count:     1,
	}
	resp, err := invoker.SdkMgr.DeleteTableRow(params)
	mgrErrHandler(err, string(resp.Resp.Body()))
	testRes = testMgrResHandler(apiName, resp, resp, err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestCreateExcelSheet creates a new spreadsheet sheet
func TestCreateExcelSheet(ctx context.Context, fileId string, sheetName string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiCreateExcelSheet
	printStartTestMsg(apiName.Name())
	params := sdkapi.AddTableSheetParams{
		Auth:   utils.GetAuth(userId),
		FileId: fileId,
		AddTableSheetReqBody: sdkapi.AddTableSheetReqBody{
			Name: sheetName,
		},
	}
	resp, err := invoker.SdkMgr.AddTableSheet(params)
	mgrErrHandler(err, string(resp.Resp.Body()))
	testRes = testMgrResHandler(apiName, resp, resp, err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestGetDocProBookmark reads a bookmark in a document pro file
func TestGetDocProBookmark(ctx context.Context, fileId string, bookmark []string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiGetDocProBookmark
	printStartTestMsg(apiName.Name())
	params := sdkapi.ReadBookmarkContentParams{
		Auth:      utils.GetAuth(userId),
		FileId:    fileId,
		Bookmarks: bookmark,
	}
	resp, err := invoker.SdkMgr.ReadBookmarkContent(params)
	mgrErrHandler(err, string(resp.Resp.Body()))
	testRes = testMgrResHandler(apiName, resp, resp.RawResponse, err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestReplaceDocProBookmark replaces a bookmark in a document pro file
func TestReplaceDocProBookmark(ctx context.Context, fileId string, req sdkapi.RepBookmarkContentReqBody) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiReplaceDocProBookmark
	printStartTestMsg(apiName.Name())
	params := sdkapi.RepBookmarkContentParams{
		Auth:                      utils.GetAuth(userId),
		FileId:                    fileId,
		RepBookmarkContentReqBody: req,
	}
	resp, err := invoker.SdkMgr.ReplaceBookmarkContent(params)
	mgrErrHandler(err, string(resp.Resp.Body()))
	testRes = testMgrResHandler(apiName, resp, resp, err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestGetAppDetails fetches app details
func TestGetAppDetails(ctx context.Context, appId string) (resp sdkapi.GetAppDetailRespBody, testRes consts.SingleApiTestRes, err error) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiGetAppDetails
	printStartTestMsg(apiName.Name())
	params := sdkapi.GetAppDetailParams{
		Auth:  utils.GetAuth(userId, true),
		AppId: appId,
	}
	resp, err = invoker.SdkMgr.GetAppDetail(params)
	mgrErrHandler(err, string(resp.Resp.Body()))
	testRes = testMgrResHandler(apiName, resp, resp.RawResponse, err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestUpdateAppEndpoint updates the callback endpoint
func TestUpdateAppEndpoint(ctx context.Context, appId string, url string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiUpdateAppEndpoint
	printStartTestMsg(apiName.Name())
	params := sdkapi.UpdateCallbackUrlParams{
		Auth:  utils.GetAuth(userId, true),
		AppId: appId,
		UpdateCallbackUrlReqBody: sdkapi.UpdateCallbackUrlReqBody{
			Url: url,
		},
	}
	resp, err := invoker.SdkMgr.UpdateCallbackUrl(params)
	mgrErrHandler(err, string(resp.Resp.Body()))
	testRes = testMgrResHandler(apiName, resp, resp, err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestGetUsersWithStatus fetches users and their seat status
func TestGetUsersWithStatus(ctx context.Context, page int, size int) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiGetUsersWithStatus
	printStartTestMsg(apiName.Name())
	params := sdkapi.GetUserAndStatusParams{
		Auth: utils.GetAuth(userId, true),
		Page: page,
		Size: size,
	}
	resp, _, err := invoker.SdkMgr.GetUserAndStatus(params)
	mgrErrHandler(err, string(resp.Resp.Body()))
	testRes = testMgrResHandler(apiName, resp, resp, err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestActivateUsers activates user seats in bulk
func TestActivateUsers(ctx context.Context, userIds []string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiActivateUsers
	printStartTestMsg(apiName.Name())
	params := sdkapi.ActivateUserSeatParams{
		Auth: utils.GetAuth(userId, true),
		ActivateUserSeatReqBody: sdkapi.ActivateUserSeatReqBody{
			UserIds: userIds,
		},
	}
	resp, err := invoker.SdkMgr.ActivateUserSeat(params)
	mgrErrHandler(err, string(resp.Resp.Body()))
	testRes = testMgrResHandler(apiName, resp, resp, err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestDeactivateUsers deactivates user seats in bulk
func TestDeactivateUsers(ctx context.Context, userIds []string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiDeactivateUsers
	printStartTestMsg(apiName.Name())
	params := sdkapi.CancelUserSeatParams{
		Auth: utils.GetAuth(userId, true),
		CancelUserSeatReqBody: sdkapi.CancelUserSeatReqBody{
			UserIds: userIds,
		},
	}
	resp, err := invoker.SdkMgr.CancelUserSeat(params)
	mgrErrHandler(err, string(resp.Resp.Body()))
	testRes = testMgrResHandler(apiName, resp, resp, err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestBatchSetUsersStatus updates seat status in bulk
func TestBatchSetUsersStatus(ctx context.Context, userIds []string, status int) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiBatchSetUsersStatus
	printStartTestMsg(apiName.Name())
	params := sdkapi.BatchSetUserSeatParams{
		Auth: utils.GetAuth(userId, true),
		BatchSetUserSeatReqBody: sdkapi.BatchSetUserSeatReqBody{
			UserIds: userIds,
			Status:  status,
		},
		Status: status,
	}
	resp, err := invoker.SdkMgr.BatchSetUserSeat(params)
	mgrErrHandler(err, string(resp.Resp.Body()))
	testRes = testMgrResHandler(apiName, resp, resp, err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

func TestGetSystemMessages(ctx context.Context, from, to, appIdQuery string) (testRes consts.SingleApiTestRes) {
	// startTime := time.Now()
	apiName := consts.ShimoSdkApiErrorCallback
	printStartTestMsg(apiName.Name())
	// resp, httpCode, err, pathStr, query := invoker.Shimo.GetSystemMessages(ctx, getToken(), from, to, appIdQuery)
	// if err != nil {
	// 	errHandler(err)
	// }
	// testRes = testResHandler(apiName, resp, httpCode, err, pathStr, query, nil, nil, time.Now().Sub(startTime).String(), startTime.Unix())
	return
}

func TestMentionAtList(ctx context.Context, fileGuid string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiGetMentionAt
	printStartTestMsg(apiName.Name())
	params := sdkapi.GetMentionAtParams{
		Auth:   utils.GetAuth(userId),
		FileId: fileGuid,
	}
	resp, err := invoker.SdkMgr.GetMentionAt(params)
	mgrErrHandler(err, string(resp.Resp.Body()))
	testRes = testMgrResHandler(apiName, resp, resp.RawResponse, err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

func TestCommentCount(ctx context.Context, fileGuid string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiGetCommentCount
	printStartTestMsg(apiName.Name())
	params := sdkapi.GetCommentCountParams{
		Auth:   utils.GetAuth(userId),
		FileId: fileGuid,
	}
	resp, err := invoker.SdkMgr.GetCommentCount(params)
	mgrErrHandler(err, string(resp.Resp.Body()))
	testRes = testMgrResHandler(apiName, resp, resp.RawResponse, err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TODO consider generating tokens from uid later
func getToken() string {
	return utils.SignUserJWT(userId)
}

func getUserId() int64 {
	email := econf.GetString("sdkctl.email")
	appId := econf.GetString("shimoSDK.appId")
	defaultPassword := econf.GetString("sdkctl.defaultPassword")
	user := db.User{
		Email: email,
		AppID: appId,
	}
	ins, err := db.FindUserByInstance(invoker.DB, &user)
	if err != nil {
		// Create the configured user if it does not exist
		ins = &db.User{
			Email:    email,
			Password: utils.HashPassword(defaultPassword),
			Avatar:   fmt.Sprintf("%sstatic/img/default-avatar-moke.png", econf.GetString("publicPath.publicPath")),
			AppID:    appId,
			Name:     strings.Split(email, "@")[0],
		}
		err = db.CreateUser(invoker.DB, ins)
		if err != nil {
			errHandler(err)
		}
	}
	return ins.ID
}

func errHandler(err error) {
	if err != nil {
		elog.Error("Request error: " + err.Error())
	}
}

func mgrErrHandler(err error, message string) {
	if err != nil {
		elog.Error("Request error: " + message)
	}
}

func testResHandler(apiName consts.ShimoSdkApi, rawResp interface{}, httpCode int, err error, pathStr string, query string, bodyReq interface{}, formData interface{}, timeConsuming string, startTime int64) consts.SingleApiTestRes {
	data, _ := json.Marshal(rawResp)
	body, _ := json.Marshal(bodyReq)
	form, _ := json.Marshal(formData)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return consts.SingleApiTestRes{
		ApiName:       apiName.Name(),
		Success:       err == nil,
		HttpCode:      httpCode,
		HttpResp:      handleNullData(string(data)),
		ErrMsg:        errMsg,
		PathStr:       pathStr,
		BodyReq:       handleNullData(string(body)),
		FormData:      handleNullData(string(form)),
		Query:         query,
		TimeConsuming: timeConsuming,
		StartTime:     startTime,
	}
}

func testMgrResHandler(apiName consts.ShimoSdkApi, res interface{}, rawResp sdkapi.RawResponse, err error, timeConsuming string, startTime int64, validators []utils.Validator) consts.SingleApiTestRes {
	data, _ := json.Marshal(res)
	body, _ := json.Marshal(rawResp.Resp.Request.Body)
	form, _ := json.Marshal(rawResp.Resp.Request.FormData)
	errMsg := ""
	if err != nil {
		errMsg = string(rawResp.Resp.Body())
	}
	success := err == nil
	// Perform assertions
	if len(validators) > 0 && err == nil {
		var validationResults []utils.ValidationResult
		for _, validator := range validators {
			// Retrieve the check field value
			checkValue, err := utils.GetCheckValue(data, validator.Check)
			if err != nil {
				success = false
				// Log an error if the check value cannot be obtained
				validationResults = append(validationResults, utils.ValidationResult{
					Validator:   validator,
					CheckResult: fmt.Sprintf("获取 check 值失败: %v", err),
				})
				// Abort further assertions when retrieval fails
				break
			}
			// Execute the assertion
			_, err = utils.Validate(checkValue, validator.Assert, validator.Expect)
			if err != nil {
				success = false
				// Log an error when the assertion fails
				validationResults = append(validationResults, utils.ValidationResult{
					Validator:   validator,
					CheckValue:  checkValue,
					CheckResult: fmt.Sprintf("断言失败: %v", err),
				})
				// Stop additional assertions if a failure occurs
				break
			} else {
				// Assertion succeeded
				validationResults = append(validationResults, utils.ValidationResult{
					Validator:   validator,
					CheckValue:  checkValue,
					CheckResult: "断言成功",
				})
			}
		}
		if !success {
			errMsg = utils.GetInterfaceToString(validationResults)
		}
	}
	return consts.SingleApiTestRes{
		ApiName:       apiName.Name(),
		Success:       success,
		HttpCode:      rawResp.Resp.StatusCode(),
		HttpResp:      handleNullData(string(data)),
		ErrMsg:        errMsg,
		PathStr:       rawResp.Resp.Request.URL,
		BodyReq:       handleNullData(string(body)),
		FormData:      handleNullData(string(form)),
		Query:         rawResp.Resp.RawResponse.Request.URL.RawQuery,
		TimeConsuming: timeConsuming,
		StartTime:     startTime,
	}
}

// generateTestResult produces the test results shown by the frontend
func generateTestResult(apiName consts.ShimoSdkApi, res *consts.SingleApiTestRes, err error) consts.SingleApiTestRes {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return consts.SingleApiTestRes{
		ApiName:       apiName.Name(),
		Success:       err == nil,
		HttpCode:      res.HttpCode,
		HttpResp:      handleNullData(res.HttpResp),
		ErrMsg:        errMsg,
		PathStr:       res.PathStr,
		BodyReq:       handleNullData(res.BodyReq),
		FormData:      handleNullData(res.FormData),
		Query:         res.Query,
		TimeConsuming: res.TimeConsuming,
		StartTime:     res.StartTime,
	}
}

func handleNullData(data string) (string string) {
	if data == "null" {
		return ""
	} else {
		return data
	}
}

func prettyJsonPrint(jsonStr []byte) {
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, jsonStr, "", "  ")
	if error != nil {
		log.Println("prettyJsonPrint error: ", error)
		log.Println("origin str: " + string(jsonStr))
		return
	}

	fmt.Println(string(prettyJSON.Bytes()))
}

// getFile return a FileHeader of local file according to its path
func getFile(path string) (*multipart.FileHeader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// create a buffer to hold the file in memory
	var buff bytes.Buffer
	buffWriter := io.Writer(&buff)

	// create a new form and create a new file field
	formWriter := multipart.NewWriter(buffWriter)
	formPart, err := formWriter.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return nil, err
	}

	// copy the content of the file to the form's file field
	if _, err := io.Copy(formPart, file); err != nil {
		return nil, err
	}

	// close the form writer after the copying process is finished
	// I don't use defer in here to avoid unexpected EOF error
	formWriter.Close()

	// transform the bytes buffer into a form reader
	buffReader := bytes.NewReader(buff.Bytes())
	formReader := multipart.NewReader(buffReader, formWriter.Boundary())

	// read the form components with max stored memory of 1MB
	multipartForm, err := formReader.ReadForm(1 << 20)
	if err != nil {
		return nil, err
	}

	// return the multipart file header
	files, exists := multipartForm.File["file"]
	if !exists || len(files) == 0 {
		return nil, err
	}

	return files[0], nil
}

func printStartTestMsg(name string) {
	fmt.Print("\n--------\n开始请求 " + name)
}

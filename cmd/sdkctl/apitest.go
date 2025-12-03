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

	"github.com/go-resty/resty/v2"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/pkg/errors"
	sdkapi "github.com/shimo-open/sdk-kit-go/api"

	"sdk-demo-go/pkg/consts"
	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/models/db"
	"sdk-demo-go/pkg/utils"
)

// TestCreatePreview creates a file preview
func TestCreatePreview(ctx context.Context, fileID string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiCreatePreview
	printStartTestMsg(apiName.Name())
	auth := utils.GetAuth(userId)
	params := sdkapi.CreatePreviewReq{
		Metadata: auth,
		FileID:   fileID,
	}
	// Create the preview
	res, err := invoker.SdkMgr.CreatePreview(ctx, params)
	if err != nil {
		mgrErrHandler(err, string(res.Response().Body()))
	}
	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestGetPreview TODO open the file preview (consider chromeless)
func TestGetPreview(ctx context.Context, fileID string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiGetPreview
	printStartTestMsg(apiName.Name())
	params := sdkapi.CreatePreviewReq{
		Metadata: utils.GetAuth(userId),
		FileID:   fileID,
	}
	// Create the preview
	res, err := invoker.SdkMgr.CreatePreview(ctx, params)
	if err != nil {
		errHandler(err)
	}

	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
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
	cFile := sdkapi.CreateFileReq{
		Metadata: auth,
		FileType: sdkapi.CollabFileType(file.ShimoType),
		Lang:     sdkapi.Lang(lang),
		FileID:   file.Guid,
	}
	res, err := invoker.SdkMgr.CreateFile(ctx, cFile)
	if err != nil {
		// Remove the record if creation fails
		_ = db.RemoveFileById(invoker.DB, file.ID)
		mgrErrHandler(err, string(res.Response().Body()))
	}
	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestCreateFileCopy creates a copy of a collaborative document
func TestCreateFileCopy(ctx context.Context, fileID string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiCreateFileCopy
	printStartTestMsg(apiName.Name())
	userId := getUserId()

	// Find the file to copy
	file, err := db.FindFileByGuid(invoker.DB, fileID)
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
	params := sdkapi.CreateFileCopyReq{
		Metadata:     auth,
		OriginFileID: fileID,
		TargetFileID: newFile.Guid,
	}
	res, err := invoker.SdkMgr.CreateFileCopy(ctx, params)

	// Call the API and capture the result
	mgrErrHandler(err, string(res.Response().Body()))

	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestDeleteFile deletes a collaborative document
func TestDeleteFile(ctx context.Context, fileID string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiDeleteFile
	printStartTestMsg(apiName.Name())
	file, err := db.FindFileByGuid(invoker.DB, fileID)
	if err != nil {
		errHandler(err)
		return
	}

	var res sdkapi.DeleteFileRes
	if file.IsShimoFile == 1 {
		auth := utils.GetAuth(userId)
		params := sdkapi.DeleteFileReq{
			Metadata: auth,
			FileID:   fileID,
		}
		res, err = invoker.SdkMgr.DeleteFile(ctx, params)
		mgrErrHandler(err, string(res.Response().Body()))
	} else {
		err = invoker.Services.AwosService.Remove(file.Guid)
		if err != nil {
			errHandler(err)
			fmt.Println(fmt.Sprintf("awos remove file failed, guid: %s, err:%e", file.Guid, err))
		}
	}

	_ = db.RemoveFileByGuid(invoker.DB, fileID)

	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
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
		FileID:   file.Guid,
		Type:     file.ShimoType,
		File:     f,
		FileName: fileName,
	}
	params := sdkapi.ImportFileReq{
		Metadata:          utils.GetAuth(userId),
		ImportFileReqBody: body,
	}
	// Upload the file via the SDK
	res, err := invoker.SdkMgr.ImportFile(ctx, params)
	if err != nil {
		fmt.Println("invoker.SdkMgr.ImportFile err:" + err.Error())
		// Roll back the newly created file
		_ = db.RemoveFileById(invoker.DB, file.ID)
		mgrErrHandler(err, string(res.Response().Body()))
	}
	return res.Data.TaskID, file.Guid, res.Response().StatusCode(), err
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
		FileID:   file.Guid,
		Type:     file.ShimoType,
		FileUrl:  url,
		FileName: file.Name,
	}
	params := sdkapi.ImportFileReq{
		Metadata:          utils.GetAuth(userId),
		ImportFileReqBody: body,
	}
	// Upload the file via the SDK
	res, err := invoker.SdkMgr.ImportFile(ctx, params)
	if err != nil {
		_ = db.RemoveFileByGuid(invoker.DB, file.Guid)
		mgrErrHandler(err, string(res.Response().Body()))
	}
	return
}

// TestImportFileProgress checks import progress
func TestImportFileProgress(ctx context.Context, taskId string) (importSuccess bool, statusCode int, err error) {
	apiName := consts.ShimoSdkApiImportFileProgress
	printStartTestMsg(apiName.Name())
	auth := utils.GetAuth(userId)
	params := sdkapi.GetImportProgReq{
		Metadata: auth,
		TaskId:   taskId,
	}
	// Fetch the upload progress
	res, err := invoker.SdkMgr.GetImportProgress(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))
	return res.Data.Progress == 100, res.Response().StatusCode(), err
}

// TestExportFile exports a file
func TestExportFile(ctx context.Context, fileID string, exportType string) (taskId string) {
	printStartTestMsg("导出文件")
	auth := utils.GetAuth(userId)
	params := sdkapi.ExportFileReq{
		Metadata: auth,
		FileID:   fileID,
		Type:     exportType,
	}
	res, err := invoker.SdkMgr.ExportFile(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))

	return res.Data.TaskID
}

// TestExportFileProgress checks export progress
func TestExportFileProgress(ctx context.Context, taskId string) (exportSuccess bool, statusCode int, err error) {
	printStartTestMsg("导出进度")
	auth := utils.GetAuth(userId)
	params := sdkapi.GetExportProgReq{
		Metadata: auth,
		TaskId:   taskId,
	}
	res, err := invoker.SdkMgr.GetExportProgress(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))
	fmt.Print(fmt.Sprintf(" 进度：%d ", res.Data.Progress))
	return res.Data.Progress == 100, res.Response().StatusCode(), err
}

// TestExportTableAsSheets exports application tables to Excel
func TestExportTableAsSheets(ctx context.Context, fileID string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiExportTableAsExcel
	printStartTestMsg(apiName.Name())
	params := sdkapi.ExportTableSheetsReq{
		Metadata: utils.GetAuth(userId),
		FileID:   fileID,
	}
	res, err := invoker.SdkMgr.ExportTableSheets(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))
	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestGetFilePlainText retrieves the plain-text content of a document
func TestGetFilePlainText(ctx context.Context, fileID string) (testRes consts.SingleApiTestRes) {
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
	params := sdkapi.GetPlainTextReq{
		FileID:   fileID,
		Metadata: utils.GetAuth(userId),
	}
	res, err := invoker.SdkMgr.GetPlainText(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))
	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), validators)
	return
}

// TestGetFilePlainTextWC retrieves the plain-text word count
func TestGetFilePlainTextWC(ctx context.Context, fileID string) (testRes consts.SingleApiTestRes) {
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
	params := sdkapi.GetPlainTextWCReq{
		Metadata: utils.GetAuth(userId),
		FileID:   fileID,
	}
	res, err := invoker.SdkMgr.GetPlainTextWC(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))
	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), validators)
	return
}

// TestGetDocSidebarInfo fetches the history list
func TestGetDocSidebarInfo(ctx context.Context, fileID string, page int, size int) (testRes consts.SingleApiTestRes) {
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
	params := sdkapi.GetHistoryListReq{
		Metadata: utils.GetAuth(userId),
		FileID:   fileID,
		PageSize: size,
		Count:    (page - 1) * size,
	}
	res, err := invoker.SdkMgr.GetHistoryList(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))
	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), validators)
	return
}

// TestGetRevision retrieves the revision list
func TestGetRevision(ctx context.Context, fileID string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiGetRevision
	printStartTestMsg(apiName.Name())
	params := sdkapi.GetRevisionListReq{
		Metadata: utils.GetAuth(userId),
		FileID:   fileID,
	}
	res, err := invoker.SdkMgr.GetRevisionList(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))
	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

func TestGetExcelContent(ctx context.Context, fileID string, rg string, validators []utils.Validator) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiGetExcelContent
	printStartTestMsg(apiName.Name())
	params := sdkapi.GetTableContentReq{
		Metadata: utils.GetAuth(userId),
		FileID:   fileID,
		Rg:       rg,
	}
	res, err := invoker.SdkMgr.GetTableContent(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))
	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), validators)
	return
}

// TestUpdateExcelContent updates spreadsheet content
func TestUpdateExcelContent(ctx context.Context, fileID string, rg string, values []string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiUpdateExcelContent
	printStartTestMsg(apiName.Name())
	realValues, err := utils.ParseStringToObject(values)
	if err != nil {
		errHandler(err)
		return
	}
	params := sdkapi.UpdateTableContentReq{
		Metadata: utils.GetAuth(userId),
		FileID:   fileID,
		UpdateTableContentRequestBody: sdkapi.UpdateTableContentRequestBody{
			Rg: rg,
			Resource: struct {
				Values [][]interface{} `json:"values"`
			}{Values: realValues},
		},
	}
	res, err := invoker.SdkMgr.UpdateTableContent(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))
	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestAppendExcelContent appends spreadsheet content
func TestAppendExcelContent(ctx context.Context, fileID string, rg string, values []string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiAppendExcelContent
	printStartTestMsg(apiName.Name())
	realValues, err := utils.ParseStringToObject(values)
	if err != nil {
		errHandler(err)
		return
	}
	params := sdkapi.AppendTableContentReq{
		Metadata: utils.GetAuth(userId),
		FileID:   fileID,
		AppendTableContentReqBody: sdkapi.AppendTableContentReqBody{
			Rg: rg,
			Resource: sdkapi.Resource{
				Values: realValues,
			},
		},
	}
	res, err := invoker.SdkMgr.AppendTableContent(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))
	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestDeleteExcelRows deletes spreadsheet rows
func TestDeleteExcelRows(ctx context.Context, fileID string, sheetName string, index int) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiDeleteExcelRows
	printStartTestMsg(apiName.Name())
	params := sdkapi.DeleteTableRowReq{
		Metadata:  utils.GetAuth(userId),
		FileID:    fileID,
		SheetName: sheetName,
		Index:     index,
		Count:     1,
	}
	res, err := invoker.SdkMgr.DeleteTableRow(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))
	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestCreateExcelSheet creates a new spreadsheet sheet
func TestCreateExcelSheet(ctx context.Context, fileID string, sheetName string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiCreateExcelSheet
	printStartTestMsg(apiName.Name())
	params := sdkapi.AddTableSheetReq{
		Metadata: utils.GetAuth(userId),
		FileID:   fileID,
		AddTableSheetReqBody: sdkapi.AddTableSheetReqBody{
			Name: sheetName,
		},
	}
	res, err := invoker.SdkMgr.AddTableSheet(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))
	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestGetDocProBookmark reads a bookmark in a document pro file
func TestGetDocProBookmark(ctx context.Context, fileID string, bookmark []string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiGetDocProBookmark
	printStartTestMsg(apiName.Name())
	params := sdkapi.ReadBookmarkContentReq{
		Metadata:  utils.GetAuth(userId),
		FileID:    fileID,
		Bookmarks: bookmark,
	}
	res, err := invoker.SdkMgr.ReadBookmarkContent(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))
	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestReplaceDocProBookmark replaces a bookmark in a document pro file
func TestReplaceDocProBookmark(ctx context.Context, fileID string, req sdkapi.RepBookmarkContentReqBody) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiReplaceDocProBookmark
	printStartTestMsg(apiName.Name())
	params := sdkapi.RepBookmarkContentReq{
		Metadata:                  utils.GetAuth(userId),
		FileID:                    fileID,
		RepBookmarkContentReqBody: req,
	}
	res, err := invoker.SdkMgr.ReplaceBookmarkContent(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))
	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestGetAppDetails fetches app details
func TestGetAppDetails(ctx context.Context, appId string) (res sdkapi.GetAppDetailRes, testRes consts.SingleApiTestRes, err error) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiGetAppDetails
	printStartTestMsg(apiName.Name())
	params := sdkapi.GetAppDetailReq{
		Metadata: utils.GetAuth(userId),
		AppID:    appId,
	}
	res, err = invoker.SdkMgr.GetAppDetail(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))
	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestUpdateAppEndpoint updates the callback endpoint
func TestUpdateAppEndpoint(ctx context.Context, appId string, url string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiUpdateAppEndpoint
	printStartTestMsg(apiName.Name())
	params := sdkapi.UpdateCallbackURLReq{
		Metadata: utils.GetAuth(userId),
		AppID:    appId,
		UpdateCallbackURLReqBody: sdkapi.UpdateCallbackURLReqBody{
			URL: url,
		},
	}
	res, err := invoker.SdkMgr.UpdateCallbackURL(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))
	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestGetUsersWithStatus fetches users and their seat status
func TestGetUsersWithStatus(ctx context.Context, page int, size int) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiGetUsersWithStatus
	printStartTestMsg(apiName.Name())
	params := sdkapi.GetUserAndStatusReq{
		Metadata: utils.GetAuth(userId),
		Page:     page,
		Size:     size,
	}
	res, err := invoker.SdkMgr.GetUserAndStatus(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))
	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestActivateUsers activates user seats in bulk
func TestActivateUsers(ctx context.Context, userIds []string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiActivateUsers
	printStartTestMsg(apiName.Name())
	params := sdkapi.ActivateUserSeatReq{
		Metadata: utils.GetAuth(userId),
		ActivateUserSeatReqBody: sdkapi.ActivateUserSeatReqBody{
			UserIds: userIds,
		},
	}
	res, err := invoker.SdkMgr.ActivateUserSeat(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))
	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestDeactivateUsers deactivates user seats in bulk
func TestDeactivateUsers(ctx context.Context, userIds []string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiDeactivateUsers
	printStartTestMsg(apiName.Name())
	params := sdkapi.CancelUserSeatReq{
		Metadata: utils.GetAuth(userId),
		CancelUserSeatReqBody: sdkapi.CancelUserSeatReqBody{
			UserIds: userIds,
		},
	}
	res, err := invoker.SdkMgr.CancelUserSeat(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))
	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

// TestBatchSetUsersStatus updates seat status in bulk
func TestBatchSetUsersStatus(ctx context.Context, userIds []string, status int) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiBatchSetUsersStatus
	printStartTestMsg(apiName.Name())
	params := sdkapi.BatchSetUserSeatReq{
		Metadata: utils.GetAuth(userId),
		BatchSetUserSeatReqBody: sdkapi.BatchSetUserSeatReqBody{
			UserIds: userIds,
			Status:  status,
		},
		Status: status,
	}
	res, err := invoker.SdkMgr.BatchSetUserSeat(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))
	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

func TestGetSystemMessages(ctx context.Context, from, to, appIdQuery string) (testRes consts.SingleApiTestRes) {
	// startTime := time.Now()
	apiName := consts.ShimoSdkApiErrorCallback
	printStartTestMsg(apiName.Name())
	// res, httpCode, err, pathStr, query := invoker.Shimo.GetSystemMessages(ctx, getToken(), from, to, appIdQuery)
	// if err != nil {
	// 	errHandler(err)
	// }
	// testRes = testResHandler(apiName, res, httpCode, err, pathStr, query, nil, nil, time.Now().Sub(startTime).String(), startTime.Unix())
	return
}

func TestMentionAtList(ctx context.Context, fileGuid string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiGetMentionAt
	printStartTestMsg(apiName.Name())
	params := sdkapi.GetMentionAtReq{
		Metadata: utils.GetAuth(userId),
		FileID:   fileGuid,
	}
	res, err := invoker.SdkMgr.GetMentionAt(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))
	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
	return
}

func TestCommentCount(ctx context.Context, fileGuid string) (testRes consts.SingleApiTestRes) {
	startTime := time.Now()
	apiName := consts.ShimoSdkApiGetCommentCount
	printStartTestMsg(apiName.Name())
	params := sdkapi.GetCommentCountReq{
		Metadata: utils.GetAuth(userId),
		FileID:   fileGuid,
	}
	res, err := invoker.SdkMgr.GetCommentCount(ctx, params)
	mgrErrHandler(err, string(res.Response().Body()))
	testRes = testMgrResHandler(apiName, res, res.Response(), err, time.Now().Sub(startTime).String(), startTime.Unix(), nil)
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

func testMgrResHandler(apiName consts.ShimoSdkApi, res interface{}, rawResp *resty.Response, err error, timeConsuming string, startTime int64, validators []utils.Validator) consts.SingleApiTestRes {
	data, _ := json.Marshal(res)
	body, _ := json.Marshal(rawResp.Request.Body)
	form, _ := json.Marshal(rawResp.Request.FormData)
	errMsg := ""
	if err != nil {
		errMsg = string(rawResp.Body())
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
		HttpCode:      rawResp.StatusCode(),
		HttpResp:      handleNullData(string(data)),
		ErrMsg:        errMsg,
		PathStr:       rawResp.Request.URL,
		BodyReq:       handleNullData(string(body)),
		FormData:      handleNullData(string(form)),
		Query:         rawResp.RawResponse.Request.URL.RawQuery,
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

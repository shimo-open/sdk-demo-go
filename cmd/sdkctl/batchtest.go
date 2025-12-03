package sdkctl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	fp "path/filepath"
	"strings"
	"time"

	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	sdk "github.com/shimo-open/sdk-kit-go"
	sdkapi "github.com/shimo-open/sdk-kit-go/api"

	"sdk-demo-go/pkg/consts"
	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/utils"
)

// TestAll performs batch tests covering every public SDK API
// Some APIs span multiple suites, so we test comprehensive coverage
// TestAll exercises all functionality
func TestAll(ctx context.Context, fileTypeStr string) (testAllRes consts.AllApiTestRes) {
	fileTypes := fileTypesByFileTypeStr(fileTypeStr)
	TestProgress = 0
	testAllRes = consts.AllApiTestRes{
		BaseTestResMap: TestBase(ctx, fileTypes),
		FileIOResMap:   TestFileIO(ctx, fileTypes),
		SpreadsheetRes: TestSpreadsheet(ctx),
		DocumentProRes: TestDocumentPro(ctx),
		TableRes:       TestTable(ctx),
		SystemRes:      TestSystem(ctx),
	}
	resStr, err := json.Marshal(&testAllRes)
	if err != nil {
		fmt.Println("TestAll json.Marshal error:", err)
		return
	}
	fmt.Println(string(resStr))

	return
}

// TestCommon runs the shared scenarios
// Includes base functionality and import/export flows
// When fileTypeStr is empty/invalid/all, every suite is tested
func TestCommon(ctx context.Context, fileTypeStr string) {
	if fileTypeStr == "" || sdk.FileType(fileTypeStr) == sdk.FileTypeInvalid {
		fileTypeStr = "all"
	}

	fileTypes := fileTypesByFileTypeStr(fileTypeStr)

	TestBase(ctx, fileTypes)
	TestFileIO(ctx, fileTypes)
}

func TestFileIO(ctx context.Context, fileTypes []sdk.FileType) (resMap map[sdk.FileType]consts.FileIORes) {
	resMap = make(map[sdk.FileType]consts.FileIORes)

	for _, ft := range fileTypes {
		// Application tables lack the shared import/export APIs
		if ft == sdk.FileTypeTable {
			continue
		}

		resMap[ft] = consts.FileIORes{
			ImportFileRes: TestImport(ctx, ft),
			ExportFileRes: TestExport(ctx, ft),
		}
	}
	TestProgress++
	return
}

// TestBase validates core functionality for every suite:
// create doc, create copy, create preview, open preview,
// fetch history list, fetch revisions, fetch plain text,
// fetch plain-text word count, fetch mention list,
// fetch comment count, and delete the document at the end
func TestBase(ctx context.Context, fileTypes []sdk.FileType) (testResMap map[sdk.FileType]consts.BaseTestRes) {
	testResMap = make(map[sdk.FileType]consts.BaseTestRes)

	for _, ft := range fileTypes {
		tmpBaseRes := consts.BaseTestRes{}
		fileId, testCreateRes, err := TestCreateFile(ctx, utils.GenFileName(ft), ft.String(), ft.String(), "", "")
		tmpBaseRes.CreateFileRes = testCreateRes
		tmpBaseRes.CreateCopyRes = TestCreateFileCopy(ctx, fileId)
		// Create a placeholder file for previewing
		err = invoker.Services.AwosService.Save(fileId, []byte{})
		if err != nil {
			elog.Error("TestBase TestCreatePreview", l.S("error: ", err.Error()))
			tmpBaseRes.CreatePreviewRes = consts.SingleApiTestRes{
				ApiName: string(consts.ShimoSdkApiCreatePreview),
				Success: false,
				ErrMsg:  fmt.Sprintf("Create cloud file error: %s", err.Error()),
			}
			TestProgress++
			return
		}
		tmpBaseRes.CreatePreviewRes = TestCreatePreview(ctx, fileId)
		tmpBaseRes.GetPreviewRes = TestGetPreview(ctx, fileId)
		tmpBaseRes.GetHistoryListRes = TestGetDocSidebarInfo(ctx, fileId, 0, 0)
		tmpBaseRes.GetRevisionListRes = TestGetRevision(ctx, fileId)

		if ft != sdk.FileTypeTable {
			fileExts := sdk.ImportTypeMap[ft]
			filePath := "resources/import/test." + fileExts[0]
			tmpBaseRes.CreateFileRes = testCreateRes
			succ, fileGuid, _, _ := doImportOnce(ctx, filePath, ft, 20*time.Second)
			if !succ {
				TestProgress++
				fmt.Println("Import failed; cannot test the plain-text APIs")
				return
			}
			// Application tables do not expose these APIs
			tmpBaseRes.GetPlainTextRes = TestGetFilePlainText(ctx, fileGuid)
			tmpBaseRes.GetPlainTextWordCountRes = TestGetFilePlainTextWC(ctx, fileGuid)
			TestDeleteFile(ctx, fileGuid)
		}

		if ft != sdk.FileTypeSlide && ft != sdk.FileTypeTable {
			// Slides and application tables lack the mention API
			tmpBaseRes.GetMentionAtListRes = TestMentionAtList(ctx, fileId)
		}

		tmpBaseRes.DeleteFileRes = TestDeleteFile(ctx, fileId)

		testResMap[ft] = tmpBaseRes
	}
	TestProgress++
	return
}

func TestImport(ctx context.Context, fileType sdk.FileType) (resMap map[sdk.FileType]map[string]consts.ImportFileRes) {
	importTypes := make([]sdk.FileType, 0)

	// If "all" is specified, test every file type
	if fileType == sdk.FileTypeAll {
		for t, _ := range sdk.ImportTypeMap {
			importTypes = append(importTypes, t)
		}
	} else {
		importTypes = append(importTypes, fileType)
	}

	resMap = make(map[sdk.FileType]map[string]consts.ImportFileRes)
	for _, t := range importTypes {
		tmpResMap := TestImportOnce(ctx, t, 20)
		resMap[t] = tmpResMap
	}
	return
}

func TestExport(ctx context.Context, fileType sdk.FileType) (resMap map[sdk.FileType]map[string]consts.ExportFileRes) {
	if fileType == sdk.FileTypeInvalid {
		return
	}

	exportTypes := make([]sdk.FileType, 0)

	if fileType == sdk.FileTypeAll {
		for t, _ := range sdk.ExportTypeMap {
			exportTypes = append(exportTypes, t)
		}
	} else {
		exportTypes = append(exportTypes, fileType)
	}

	resMap = make(map[sdk.FileType]map[string]consts.ExportFileRes)
	for _, t := range exportTypes {
		fileExts := sdk.ImportTypeMap[t]
		filePath := "resources/import/test."
		if len(fileExts) > 0 {
			filePath += fileExts[0]
		} else {
			elog.Error(fmt.Sprintf("Export failed, no corresponding import type, fileType:%s, filePath:%s", t.String(), filePath))
			continue
		}

		// Exporting requires a prior import
		// Import the first available format by default
		importSuccess, fileGuid, statusCode, errMsg := doImportOnce(ctx, filePath, t, 20*time.Second)
		if !importSuccess {
			resMap[t] = map[string]consts.ExportFileRes{
				"": {
					ErrMsg:   fmt.Sprintf("前置导入失败: %s, 文件地址: %s", errMsg, filePath),
					HttpCode: statusCode,
					Success:  false,
				},
			}
			elog.Error("Import failed, unable to export")
			continue
		}
		tmpResMap := TestExportOnce(ctx, t, fileGuid, 20)
		resMap[t] = tmpResMap
	}
	return
}

// TestImportOnce runs a full import pass for each supported extension of a file type
// Includes a timeout
func TestImportOnce(ctx context.Context, fileType sdk.FileType, timeoutSec int64) (resMap map[string]consts.ImportFileRes) {
	if fileType == sdk.FileTypeInvalid {
		return
	}

	resMap = make(map[string]consts.ImportFileRes)
	filePaths := getImportFilePathByType(fileType)
	for _, path := range filePaths {
		startTime := time.Now()
		succ, fileGuid, statusCode, errMsg := doImportOnce(ctx, path, fileType, time.Duration(timeoutSec)*time.Second)
		formData, _ := json.Marshal(consts.FileFormData{FileId: fileGuid, FilePath: path, FileType: fileType})
		resMap[fp.Ext(path)] = consts.ImportFileRes{
			Success:       succ,
			PathStr:       fmt.Sprintf("%s,%s", consts.ShimoSdkPathImportFile, consts.ShimoSdkPathImportFileProgress),
			FormData:      string(formData),
			StartTime:     startTime.Unix(),
			TimeConsuming: time.Now().Sub(startTime).String(),
			HttpCode:      statusCode,
			ErrMsg:        errMsg,
		}
	}

	return
}

func doImportOnce(ctx context.Context, filePath string, fileType sdk.FileType, timeoutSec time.Duration) (success bool, fileGuid string, statusCode int, errMsg string) {
	fmt.Println(fmt.Sprintf("Import test: type=%s, path=%s", fileType.String(), filePath))
	taskId, fileGuid, statusCode, err := TestImportFile(ctx, filePath, fileType.String(), fileType.String())
	if err != nil {
		elog.Error("import failed", l.S("error: ", err.Error()))
		return
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeout := time.After(timeoutSec)
	startTime := time.Now()

	for {
		select {
		case <-ticker.C:
			if s, c, e := TestImportFileProgress(ctx, taskId); s {
				fmt.Println("Import succeeded, duration:", time.Now().Sub(startTime).String())
				statusCode = c
				success = true
				return
			} else if e != nil {
				statusCode = c
				success = false
				errMsg = e.Error()
				return
			}
		case <-timeout:
			fmt.Println("导入超时，耗时：", time.Now().Sub(startTime).String())
			statusCode = http.StatusRequestTimeout
			success = false
			errMsg = fmt.Sprintf("导入超时, 耗时: %s", time.Now().Sub(startTime).String())
			return
		}
	}
}

type ExportFileReq struct {
	Type string `json:"type"`
}

func TestExportOnce(ctx context.Context, fileType sdk.FileType, fileId string, timeoutSec int64) (resMap map[string]consts.ExportFileRes) {
	exportExts := sdk.ExportTypeMap[fileType]
	resMap = make(map[string]consts.ExportFileRes)

	for _, ext := range exportExts {
		startTime := time.Now()
		succ, statusCode, errMsg := doExportOnce(ctx, fileType, ext, fileId, time.Duration(timeoutSec)*time.Second)
		body, _ := json.Marshal(consts.ExportBodyReq{
			Type: fileType,
		})
		resMap[ext] = consts.ExportFileRes{
			Success:       succ,
			PathStr:       fmt.Sprintf("%s%s,%s", consts.ShimoSdkPathExportFile, fileId, consts.ShimoSdkPathExportFileProgress),
			BodyReq:       string(body),
			StartTime:     startTime.Unix(),
			TimeConsuming: time.Now().Sub(startTime).String(),
			HttpCode:      statusCode,
			ErrMsg:        errMsg,
		}
	}
	return
}

func doExportOnce(ctx context.Context, fileType sdk.FileType, exportExt, fileGuid string, timeoutSec time.Duration) (success bool, statusCode int, errMsg string) {
	fmt.Print(fmt.Sprintf("\n导出测试，文件id：%s 文件类型：%s，导出类型：%s", fileGuid, fileType, exportExt))
	taskId := TestExportFile(ctx, fileGuid, exportExt)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timeout := time.After(timeoutSec)
	startTime := time.Now()

	for {
		select {
		case <-ticker.C:
			if s, c, e := TestExportFileProgress(ctx, taskId); s {
				fmt.Println("导出成功，耗时：", time.Now().Sub(startTime).String())
				statusCode = c
				success = true
				return
			} else if e != nil {
				fmt.Println("导出失败，耗时：", time.Now().Sub(startTime).String())
				statusCode = c
				success = false
				return
			}
		case <-timeout:
			fmt.Println("导出超时，耗时：", time.Now().Sub(startTime).String())
			statusCode = http.StatusRequestTimeout
			success = false
			return
		}
	}
}

// Mapping of Shimo file types to importable types
var importTypeMap = map[sdk.FileType][]string{
	sdk.FileTypeDocument:    {"docx", "doc", "md", "txt"},
	sdk.FileTypeSpreadsheet: {"xlsx", "xls", "csv", "xlsm"},
	sdk.FileTypeDocPro:      {"docx", "doc", "wps"},
	sdk.FileTypeSlide:       {"pptx", "ppt"},
}

func getImportFilePathByType(ft sdk.FileType) []string {
	folderPath := "resources/import/"
	res := []string{}
	for _, ext := range importTypeMap[ft] {
		res = append(res, folderPath+"test."+ext)
	}

	return res
}

// TestSpecial covers suite-specific functionality
// Spreadsheet suite
// Document Pro suite
func TestSpecial(ctx context.Context) {
	TestSpreadsheet(ctx)
	TestDocumentPro(ctx)
}

// TestSpreadsheet focuses on spreadsheet capabilities:
// fetch content, update content, append content, delete rows,
// add sheets, and validate using a populated test sheet
// // Entry point for spreadsheet tests
func TestSpreadsheet(ctx context.Context) (testRes consts.SpreadsheetRes) {
	testRes = consts.SpreadsheetRes{}
	// Import a spreadsheet file
	succ, fileGuid, _, errMsg := doImportOnce(ctx, "resources/import/test.xlsx", sdk.FileTypeSpreadsheet, 20*time.Second)
	if !succ {
		errResp := consts.SingleApiTestRes{Success: false, ErrMsg: fmt.Sprintf("导入失败，无法测试表格功能, err: %s", errMsg)}
		elog.Error("Import failed, unable to test spreadsheet functionality")
		testRes = consts.SpreadsheetRes{
			GetTableContentRes:    errResp,
			UpdateTableContentRes: errResp,
			AppendTableContentRes: errResp,
			DeleteTableRowRes:     errResp,
			AddTableSheetRes:      errResp,
			GetCommentCountRes:    errResp,
		}
		TestProgress++
		return
	}

	// Define test parameters
	rg := "Sheet1!A1:C5"

	initialValues := []string{
		`["姓名", "年龄", "性别"]`,
		`["小红", "25", "女"]`,
		`["小强", "30", "男"]`,
	}
	appendValue := []string{
		`["追加姓名", "35", "男"]`,
	}

	newSheetName := "新工作表"
	rowToDelete := 1 // Delete the second row (0-based)

	// Step 1: Fetch spreadsheet content
	testRes.GetTableContentRes = TestGetExcelContent(ctx, fileGuid, rg, []utils.Validator{
		{
			Check:   "values[1][0]",
			Assert:  "equal",
			Expect:  "saf",
			Message: "校验获取表格内容",
		},
	})
	// Step 2: Update spreadsheet content
	res := TestUpdateExcelContent(ctx, fileGuid, rg, initialValues)
	// Step 3: Validate the updated content
	check := TestGetExcelContent(ctx, fileGuid, rg, []utils.Validator{
		{
			Check:   "values[1][0]",
			Assert:  "equal",
			Expect:  "小红",
			Message: "校验更新表格内容",
		},
		{
			Check:   "values[2][0]",
			Assert:  "equal",
			Expect:  "小强",
			Message: "校验更新表格内容",
		},
	})
	if res.Success && !check.Success {
		res.Success = false
		res.ErrMsg = check.ErrMsg
		elog.Error("Spreadsheet update did not take effect")
	}
	testRes.UpdateTableContentRes = res

	// Step 4: Append spreadsheet content
	res = TestAppendExcelContent(ctx, fileGuid, rg, appendValue)

	// Step 5: Validate the appended content
	check = TestGetExcelContent(ctx, fileGuid, rg, []utils.Validator{
		{
			Check:   "values[4][0]",
			Assert:  "equal",
			Expect:  "追加姓名",
			Message: "校验追加后表格内容",
		},
	})
	if res.Success && !check.Success {
		res.Success = false
		res.ErrMsg = check.ErrMsg
		elog.Error("Append spreadsheet content did not take effect")
	}
	testRes.AppendTableContentRes = res

	// Step 6: Delete spreadsheet rows
	res = TestDeleteExcelRows(ctx, fileGuid, "Sheet1", rowToDelete)

	// Step 7: Validate content after row deletion
	check = TestGetExcelContent(ctx, fileGuid, rg, []utils.Validator{
		{
			Check:   "values[1][0]",
			Assert:  "not_contains",
			Expect:  "小红",
			Message: "验证删除行",
		},
	})
	if res.Success && !check.Success {
		res.Success = false
		res.ErrMsg = check.ErrMsg
		elog.Error("Row deletion did not take effect")
	}
	testRes.DeleteTableRowRes = res

	// Step 8: Add a sheet
	testRes.AddTableSheetRes = TestCreateExcelSheet(ctx, fileGuid, newSheetName)

	// Step 9: Fetch the spreadsheet comment count
	testRes.GetCommentCountRes = TestCommentCount(ctx, fileGuid)

	TestProgress++
	return
}

// TestDocumentPro covers Document Pro functionality:
// read bookmarks, replace bookmarks, and test with a document containing bookmarks
func TestDocumentPro(ctx context.Context) (testRes consts.DocumentProRes) {
	succ, fileGuid, _, errMsg := doImportOnce(ctx, "resources/import/test.docx", sdk.FileTypeDocPro, 20*time.Second)
	if !succ {
		errResp := consts.SingleApiTestRes{Success: false, ErrMsg: fmt.Sprintf("导入失败，无法测试传统文档功能, err: %s", errMsg)}
		testRes = consts.DocumentProRes{
			ReadBookmarkContentRes:    errResp,
			ReplaceBookmarkContentRes: errResp,
		}
		elog.Error("Import failed, unable to test document pro functionality")
		TestProgress++
		return
	}

	bookmarks := []string{"我是书签"}
	testRes = consts.DocumentProRes{
		ReadBookmarkContentRes: TestGetDocProBookmark(ctx, fileGuid, bookmarks),
		ReplaceBookmarkContentRes: TestReplaceDocProBookmark(ctx, fileGuid, sdkapi.RepBookmarkContentReqBody{
			Replacements: []sdkapi.Replacement{
				{
					Bookmark: "我是书签",
					Type:     "text",
					Value:    "我是替换后的内容",
				},
			},
		})}
	TestProgress++
	return
}

// TestTable handles application-table features
// Export an application table as a professional spreadsheet
func TestTable(ctx context.Context) (testRes consts.TableRes) {
	fileGuid, _, err := TestCreateFile(ctx, "应用表格转表格", sdk.FileTypeTable.String(), sdk.FileTypeTable.String(), "", "")
	if err != nil {
		errResp := consts.SingleApiTestRes{Success: false, ErrMsg: fmt.Sprintf("创建文件失败，无法测试应用表格功能, err: %s", err.Error())}
		testRes = consts.TableRes{
			SheetExportToExcelRes: errResp,
		}
		elog.Error("TestTable TestCreateFile", l.S("error: ", err.Error()))
		TestProgress++
		return
	}

	testRes = consts.TableRes{
		SheetExportToExcelRes: TestExportTableAsSheets(ctx, fileGuid),
	}
	TestProgress++
	return
}

// -------- System functionality is temporarily excluded from full tests

// TestSystem covers system-level APIs:
// App: get details, update callback URL (currently skipped)
// Users: fetch seat status, activate seats, deactivate seats, batch update seats
func TestSystem(ctx context.Context) (testRes consts.SystemRes) {
	details, getAppDetailRes, err := TestGetAppDetails(ctx, econf.GetString("shimoSDK.appId"))
	if err != nil {
		testRes = buildSystemErrorRes(fmt.Sprintf("获取app详情失败，无法测试系统功能, err: %s", err.Error()))
		TestProgress++
		elog.Error("TestGetAppDetails", l.S("error: ", err.Error()))
		return
	}
	if details.EndpointUrl == "" {
		testRes = buildSystemErrorRes(fmt.Sprintf("回调地址为空，无法测试系统功能, err: %s", err.Error()))
		TestProgress++
		elog.Error("TestGetAppDetails", l.S("error: ", "endpoint url is empty"))
		return
	}

	// Wait 2s between calls to avoid identical signatures
	time.Sleep(time.Second * 2)

	// updateAppEndpointRes := TestUpdateAppEndpoint(ctx, econf.GetString("shimoSDK.appId"), details.EndpointUrl+"/test")
	// time.Sleep(time.Second * 2)

	// Restore the original endpoint after testing
	TestUpdateAppEndpoint(ctx, econf.GetString("shimoSDK.appId"), details.EndpointUrl)
	time.Sleep(time.Second * 2)

	testUids := []string{
		"ensure_auto_increment",
	}

	getUsersWithStatusRes := TestGetUsersWithStatus(ctx, 0, 0)
	time.Sleep(time.Second * 2)

	// First deactivate seats
	deactivateUsersRes := TestDeactivateUsers(ctx, testUids)
	time.Sleep(time.Second * 2)

	// Then activate them
	activateUsersRes := TestActivateUsers(ctx, testUids)
	time.Sleep(time.Second * 2)

	// Next, batch deactivate
	TestBatchSetUsersStatus(ctx, testUids, -1)
	time.Sleep(time.Second * 2)

	// Finally, batch activate
	batchSetUsersStatusRes := TestBatchSetUsersStatus(ctx, testUids, 1)
	time.Sleep(time.Second * 2)

	testRes = consts.SystemRes{
		GetAppDetailRes: getAppDetailRes,
		// UpdateCallbackUrlRes:        updateAppEndpointRes,
		GetUserListAndSeatStatusRes: getUsersWithStatusRes,
		ActivateUserSeatRes:         activateUsersRes,
		CancelUserSeatRes:           deactivateUsersRes,
		BatchSetUserSeatRes:         batchSetUsersStatusRes,
	}
	TestProgress++
	return
}

func buildSystemErrorRes(errMsg string) consts.SystemRes {
	errResp := consts.SingleApiTestRes{
		Success: false,
		ErrMsg:  errMsg,
	}
	return consts.SystemRes{
		GetAppDetailRes:             errResp,
		GetUserListAndSeatStatusRes: errResp,
		ActivateUserSeatRes:         errResp,
		CancelUserSeatRes:           errResp,
		BatchSetUserSeatRes:         errResp,
	}
}

func fileTypesByFileTypeStr(fileTypeStr string) (fts []sdk.FileType) {
	fts = make([]sdk.FileType, 0)
	if fileTypeStr == "all" {
		fts = sdk.AllFileTypes()
	} else {
		if strings.Contains(fileTypeStr, ",") {
			s := strings.Split(fileTypeStr, ",")
			for _, v := range s {
				fts = append(fts, sdk.GetFileType(v))
			}
		} else {
			fts = append(fts, sdk.GetFileType(fileTypeStr))
		}
	}

	return fts
}

package consts

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ego-component/excelplus"
	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/ego/core/elog"
	"github.com/shimo-open/sdk-kit-go"
)

type fileStatus struct {
	ApiName    string `excel:"接口名字" excel_width:"15"`
	OriginType string ` excel:"上传文件后缀" excel_width:"25"`
	FileType   string ` excel:"石墨文件类型" excel_width:"25"`
	FileStatus bool   `excel:"请求结果" excel_width:"15"`
}

func TestSaveExcel(t *testing.T) {
	testRes := AllApiTestRes{
		BaseTestResMap: map[sdk.FileType]BaseTestRes{
			sdk.FileTypeSpreadsheet: BaseTestRes{
				CreateFileRes:            SingleApiTestRes{ApiName: "CreateFileRes", Success: true, HttpCode: 200, HttpResp: "File created", ErrMsg: ""},
				CreateCopyRes:            SingleApiTestRes{ApiName: "CreateCopyRes", Success: true, HttpCode: 200, HttpResp: "Copy created", ErrMsg: ""},
				DeleteFileRes:            SingleApiTestRes{ApiName: "DeleteFileRes", Success: true, HttpCode: 204, HttpResp: "", ErrMsg: ""},
				CreatePreviewRes:         SingleApiTestRes{ApiName: "CreatePreviewRes", Success: true, HttpCode: 200, HttpResp: "Preview created", ErrMsg: ""},
				GetPreviewRes:            SingleApiTestRes{ApiName: "GetPreviewRes", Success: true, HttpCode: 200, HttpResp: "Preview retrieved", ErrMsg: ""},
				GetHistoryListRes:        SingleApiTestRes{ApiName: "GetHistoryListRes", Success: true, HttpCode: 200, HttpResp: "History list retrieved", ErrMsg: ""},
				GetRevisionListRes:       SingleApiTestRes{ApiName: "GetRevisionListRes", Success: true, HttpCode: 200, HttpResp: "Revision list retrieved", ErrMsg: ""},
				GetPlainTextRes:          SingleApiTestRes{ApiName: "GetPlainTextRes", Success: true, HttpCode: 200, HttpResp: "Plain text retrieved", ErrMsg: ""},
				GetPlainTextWordCountRes: SingleApiTestRes{ApiName: "GetPlainTextWordCountRes", Success: true, HttpCode: 200, HttpResp: "Word count retrieved", ErrMsg: ""},
				GetMentionAtListRes:      SingleApiTestRes{ApiName: "GetMentionAtListRes", Success: true, HttpCode: 200, HttpResp: "Mention list retrieved", ErrMsg: ""},
			},
		},
		FileIOResMap: map[sdk.FileType]FileIORes{
			sdk.FileTypeSpreadsheet: {
				ImportFileRes: map[sdk.FileType]map[string]ImportFileRes{
					sdk.FileTypeDocument: {
						"docx": {},
						"doc":  {Success: false, PathStr: "", FormData: ""},
					},
				},
				ExportFileRes: map[sdk.FileType]map[string]ExportFileRes{
					sdk.FileTypeDocument: {
						"docx": {Success: true, PathStr: "", BodyReq: ""},
					},
					sdk.FileTypeDocPro: {
						"wps": {Success: true, PathStr: "", BodyReq: ""},
					},
				},
			},
		},
		// Spreadsheet-specific functionality
		SpreadsheetRes: SpreadsheetRes{
			GetTableContentRes: SingleApiTestRes{ApiName: "GetTableContentRes", Success: true, HttpCode: 200, HttpResp: "File created", ErrMsg: ""},
			// Update spreadsheet content
			UpdateTableContentRes: SingleApiTestRes{ApiName: "UpdateTableContentRes", Success: true, HttpCode: 200, HttpResp: "File created", ErrMsg: ""},
		},
		// Document Pro-specific functionality
		DocumentProRes: DocumentProRes{
			ReadBookmarkContentRes: SingleApiTestRes{ApiName: "ReadBookmarkContentRes", Success: true, HttpCode: 200, HttpResp: "File created", ErrMsg: ""},
			// Replace bookmark content
			ReplaceBookmarkContentRes: SingleApiTestRes{ApiName: "ReplaceBookmarkContentRes", Success: true, HttpCode: 200, HttpResp: "File created", ErrMsg: ""},
		},
		// Application table-specific functionality
		TableRes: TableRes{
			SheetExportToExcelRes: SingleApiTestRes{ApiName: "SheetExportToExcelRes", Success: true, HttpCode: 200, HttpResp: "File created", ErrMsg: ""},
		},
		// System functionality
		SystemRes: SystemRes{
			GetAppDetailRes: SingleApiTestRes{ApiName: "GetAppDetailRes", Success: true, HttpCode: 200, HttpResp: "File created", ErrMsg: ""},
			// // Update callback URL
			// UpdateCallbackUrlRes: SingleApiTestRes{ApiName: "UpdateCallbackUrlRes", Success: true, HttpCode: 200, HttpResp: "File created", ErrMsg: ""},
		},
	}

	columnArray := [6]string{"BaseTestResMap", "FileIOResMap", "SpreadsheetRes", "DocumentProRes", "TableRes", "SystemRes"}
	exFile := excelplus.Load().Build(excelplus.WithDefaultSheetName("BaseTestResMap"))
	for _, value := range columnArray {
		sheetName := value
		switch value {
		case "BaseTestResMap":
			// Create a new excelplus sheet
			BaseTestSheet, err := exFile.NewSheet(sheetName, SingleApiTestRes{})
			BaseTestBody := make([]SingleApiTestRes, 0)
			if err != nil {
				elog.Panic("BaseTestResMap sheet error", l.E(err))
			}
			BaseTestBody = append(BaseTestBody, testRes.BaseTestResMap[sdk.FileTypeSpreadsheet].CreateFileRes)
			BaseTestBody = append(BaseTestBody, testRes.BaseTestResMap[sdk.FileTypeSpreadsheet].CreateCopyRes)
			BaseTestBody = append(BaseTestBody, testRes.BaseTestResMap[sdk.FileTypeSpreadsheet].DeleteFileRes)
			BaseTestBody = append(BaseTestBody, testRes.BaseTestResMap[sdk.FileTypeSpreadsheet].CreatePreviewRes)
			BaseTestBody = append(BaseTestBody, testRes.BaseTestResMap[sdk.FileTypeSpreadsheet].GetPreviewRes)
			BaseTestBody = append(BaseTestBody, testRes.BaseTestResMap[sdk.FileTypeSpreadsheet].GetHistoryListRes)
			BaseTestBody = append(BaseTestBody, testRes.BaseTestResMap[sdk.FileTypeSpreadsheet].GetRevisionListRes)
			BaseTestBody = append(BaseTestBody, testRes.BaseTestResMap[sdk.FileTypeSpreadsheet].GetPlainTextRes)
			BaseTestBody = append(BaseTestBody, testRes.BaseTestResMap[sdk.FileTypeSpreadsheet].GetPlainTextWordCountRes)
			BaseTestBody = append(BaseTestBody, testRes.BaseTestResMap[sdk.FileTypeSpreadsheet].GetMentionAtListRes)
			for _, value := range BaseTestBody {
				err := BaseTestSheet.SetRow(value)
				if err != nil {
					elog.Panic("set sheet body row error", l.E(err), l.S("sheetName", sheetName))
				}
			}

		case "FileIOResMap":
			// Create a new excelplus sheet
			FileIOResSheet, err := exFile.NewSheet(sheetName, fileStatus{})
			FileIOResBody := make([]fileStatus, 0)
			if err != nil {
				elog.Panic("BaseTestResMap sheet error", l.E(err))
			}

			for _, statusMap := range testRes.FileIOResMap {
				for ext, status := range statusMap.ImportFileRes {
					for key := range status {
						FileIOResBody = append(FileIOResBody, fileStatus{ApiName: "ImportFileRes", OriginType: key, FileType: string(ext), FileStatus: status[key].Success})
					}
				}
				for extE, statusE := range statusMap.ExportFileRes {
					for keyE := range statusE {
						FileIOResBody = append(FileIOResBody, fileStatus{ApiName: "ExportFileRes", OriginType: keyE, FileType: string(extE), FileStatus: statusE[keyE].Success})
					}
				}
			}
			for _, value := range FileIOResBody {
				err := FileIOResSheet.SetRow(value)
				if err != nil {
					elog.Panic("set sheet body row error", l.E(err), l.S("sheetName", sheetName))
				}
			}
		case "SpreadsheetRes":
			// Create a new excelplus sheet
			SpreadsheetSheet, err := exFile.NewSheet(sheetName, SingleApiTestRes{})
			SpreadsheetBody := make([]SingleApiTestRes, 0)
			if err != nil {
				elog.Panic("BaseTestResMap sheet error", l.E(err))
			}

			SpreadsheetBody = append(SpreadsheetBody, testRes.SpreadsheetRes.GetTableContentRes)
			SpreadsheetBody = append(SpreadsheetBody, testRes.SpreadsheetRes.UpdateTableContentRes)
			for _, value := range SpreadsheetBody {
				err := SpreadsheetSheet.SetRow(value)
				if err != nil {
					elog.Panic("set sheet body row error", l.E(err), l.S("sheetName", sheetName))
				}
			}
		case "DocumentProRes":
			// Create a new excelplus sheet
			DocumentProSheet, err := exFile.NewSheet(sheetName, SingleApiTestRes{})
			DocumentProBody := make([]SingleApiTestRes, 0)
			if err != nil {
				elog.Panic("BaseTestResMap sheet error", l.E(err))
			}

			DocumentProBody = append(DocumentProBody, testRes.DocumentProRes.ReadBookmarkContentRes)
			DocumentProBody = append(DocumentProBody, testRes.DocumentProRes.ReplaceBookmarkContentRes)
			for _, value := range DocumentProBody {
				err := DocumentProSheet.SetRow(value)
				if err != nil {
					elog.Panic("set sheet body row error", l.E(err), l.S("sheetName", sheetName))
				}
			}
		case "TableRes":
			// Create a new excelplus sheet
			TableSheet, err := exFile.NewSheet(sheetName, SingleApiTestRes{})
			TableBody := make([]SingleApiTestRes, 0)
			if err != nil {
				elog.Panic("BaseTestResMap sheet error", l.E(err))
			}

			TableBody = append(TableBody, testRes.TableRes.SheetExportToExcelRes)
			for _, value := range TableBody {
				err := TableSheet.SetRow(value)
				if err != nil {
					elog.Panic("set sheet body row error", l.E(err), l.S("sheetName", sheetName))
				}
			}
		case "SystemRes":
			// Create a new excelplus sheet
			SystemSheet, err := exFile.NewSheet(sheetName, SingleApiTestRes{})
			SystemBody := make([]SingleApiTestRes, 0)
			if err != nil {
				elog.Panic("BaseTestResMap sheet error", l.E(err))
			}

			SystemBody = append(SystemBody, testRes.SystemRes.GetAppDetailRes)
			// SystemBody = append(SystemBody, testRes.SystemRes.UpdateCallbackUrlRes)
			for _, value := range SystemBody {
				err := SystemSheet.SetRow(value)
				if err != nil {
					elog.Panic("set sheet body row error", l.E(err), l.S("sheetName", sheetName))
				}
			}
		}
	}

	fileName := fmt.Sprintf("/Users/xuyixian/Downloads/test-result/files/%s/sdk_test_all_%s.xlsx", time.Now().Format("2006_01_02"), time.Now().Format("2006_01_02_15_04_05"))
	if err := exFile.SaveAs(context.Background(), fileName); err != nil {
		elog.Panic("save excel error", l.E(err))
	}
}

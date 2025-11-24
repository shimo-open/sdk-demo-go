package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/ego-component/excelplus"
	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/ego/core/elog"

	"sdk-demo-go/pkg/consts"
)

type fileStatus struct {
	ApiName    string `excel:"接口名字" excel_width:"25"`
	OriginType string `excel:"上传文件后缀" excel_width:"20"`
	FileStatus bool   `excel:"请求结果" excel_width:"15"`
}

// TestSaveExcel generates an Excel report from API test results
func TestSaveExcel(testRes consts.AllApiTestRes) {
	exFile := excelplus.Load().Build(excelplus.WithDefaultSheetName(consts.SheetNameMap["BaseTestResMap"]))
	for k, value := range consts.SheetNameMap {
		sheetName := value
		switch k {
		case "BaseTestResMap":
			// Create a new excelplus sheet
			BaseTestSheet, err := exFile.NewSheet(sheetName, consts.SingleApiTestRes{})
			BaseTestBody := make([]consts.SingleApiTestRes, 0)
			if err != nil {
				elog.Panic("BaseTestResMap sheet error", l.E(err))
			}
			for key, _ := range testRes.BaseTestResMap {
				if key != consts.FileTypeTable {
					// Application tables do not expose these APIs
					BaseTestBody = append(BaseTestBody, testRes.BaseTestResMap[key].GetPlainTextRes)
					BaseTestBody = append(BaseTestBody, testRes.BaseTestResMap[key].GetPlainTextWordCountRes)
				}
				if key != consts.FileTypeSlide {
					// Slides lack the mention API
					BaseTestBody = append(BaseTestBody, testRes.BaseTestResMap[key].GetMentionAtListRes)
				}
				BaseTestBody = append(BaseTestBody, testRes.BaseTestResMap[key].CreateFileRes)
				BaseTestBody = append(BaseTestBody, testRes.BaseTestResMap[key].CreateCopyRes)
				BaseTestBody = append(BaseTestBody, testRes.BaseTestResMap[key].DeleteFileRes)
				BaseTestBody = append(BaseTestBody, testRes.BaseTestResMap[key].CreatePreviewRes)
				BaseTestBody = append(BaseTestBody, testRes.BaseTestResMap[key].GetPreviewRes)
				BaseTestBody = append(BaseTestBody, testRes.BaseTestResMap[key].GetHistoryListRes)
				BaseTestBody = append(BaseTestBody, testRes.BaseTestResMap[key].GetRevisionListRes)

			}
			for _, value := range BaseTestBody {
				err := BaseTestSheet.SetRow(value)
				if err != nil {
					elog.Panic("set BaseTestSheet body row error", l.E(err), l.S("sheetName", sheetName))
				}
			}

		case "FileIOResMap":
			// Create a new excelplus sheet
			FileIOResSheet, err := exFile.NewSheet(sheetName, fileStatus{})
			FileIOResBody := make([]fileStatus, 0)
			if err != nil {
				elog.Panic("FileIOResMap sheet error", l.E(err))
			}

			for _, statusMap := range testRes.FileIOResMap {
				for ext, status := range statusMap.ImportFileRes {
					for key := range status {
						FileIOResBody = append(FileIOResBody, fileStatus{ApiName: "导入文件" + "-" + string(ext), OriginType: key, FileStatus: status[key].Success})
					}
				}
				for extE, statusE := range statusMap.ExportFileRes {
					for keyE := range statusE {
						FileIOResBody = append(FileIOResBody, fileStatus{ApiName: "导出文件" + "-" + string(extE), OriginType: keyE, FileStatus: statusE[keyE].Success})
					}
				}
			}
			for _, value := range FileIOResBody {
				err := FileIOResSheet.SetRow(value)
				if err != nil {
					elog.Panic("set FileIOResSheet body row error", l.E(err), l.S("sheetName", sheetName))
				}
			}
		case "SpreadsheetRes":
			// Create a new excelplus sheet
			SpreadsheetSheet, err := exFile.NewSheet(sheetName, consts.SingleApiTestRes{})
			SpreadsheetBody := make([]consts.SingleApiTestRes, 0)
			if err != nil {
				elog.Panic("SpreadsheetRes sheet error", l.E(err))
			}

			SpreadsheetBody = append(SpreadsheetBody, testRes.SpreadsheetRes.GetTableContentRes)
			SpreadsheetBody = append(SpreadsheetBody, testRes.SpreadsheetRes.UpdateTableContentRes)
			SpreadsheetBody = append(SpreadsheetBody, testRes.SpreadsheetRes.AppendTableContentRes)
			SpreadsheetBody = append(SpreadsheetBody, testRes.SpreadsheetRes.AddTableSheetRes)
			SpreadsheetBody = append(SpreadsheetBody, testRes.SpreadsheetRes.DeleteTableRowRes)
			SpreadsheetBody = append(SpreadsheetBody, testRes.SpreadsheetRes.GetCommentCountRes)
			for _, value := range SpreadsheetBody {
				err := SpreadsheetSheet.SetRow(value)
				if err != nil {
					elog.Panic("set SpreadsheetBody body row error", l.E(err), l.S("sheetName", sheetName))
				}
			}
		case "DocumentProRes":
			// Create a new excelplus sheet
			DocumentProSheet, err := exFile.NewSheet(sheetName, consts.SingleApiTestRes{})
			DocumentProBody := make([]consts.SingleApiTestRes, 0)
			if err != nil {
				elog.Panic("DocumentProRes sheet error", l.E(err))
			}

			DocumentProBody = append(DocumentProBody, testRes.DocumentProRes.ReadBookmarkContentRes)
			DocumentProBody = append(DocumentProBody, testRes.DocumentProRes.ReplaceBookmarkContentRes)
			for _, value := range DocumentProBody {
				err := DocumentProSheet.SetRow(value)
				if err != nil {
					elog.Panic("set DocumentProSheet body row error", l.E(err), l.S("sheetName", sheetName))
				}
			}
		case "TableRes":
			// Create a new excelplus sheet
			TableSheet, err := exFile.NewSheet(sheetName, consts.SingleApiTestRes{})
			TableBody := make([]consts.SingleApiTestRes, 0)
			if err != nil {
				elog.Panic("TableRes sheet error", l.E(err))
			}

			TableBody = append(TableBody, testRes.TableRes.SheetExportToExcelRes)
			for _, value := range TableBody {
				err := TableSheet.SetRow(value)
				if err != nil {
					elog.Panic("set TableSheet body row error", l.E(err), l.S("sheetName", sheetName))
				}
			}
		case "SystemRes":
			// Create a new excelplus sheet
			SystemSheet, err := exFile.NewSheet(sheetName, consts.SingleApiTestRes{})
			SystemBody := make([]consts.SingleApiTestRes, 0)
			if err != nil {
				elog.Panic("SystemRes sheet error", l.E(err))
			}

			SystemBody = append(SystemBody, testRes.SystemRes.GetAppDetailRes)
			// SystemBody = append(SystemBody, testRes.SystemRes.UpdateCallbackUrlRes)
			SystemBody = append(SystemBody, testRes.SystemRes.ActivateUserSeatRes)
			SystemBody = append(SystemBody, testRes.SystemRes.CancelUserSeatRes)
			SystemBody = append(SystemBody, testRes.SystemRes.GetUserListAndSeatStatusRes)
			SystemBody = append(SystemBody, testRes.SystemRes.BatchSetUserSeatRes)
			for _, value := range SystemBody {
				err := SystemSheet.SetRow(value)
				if err != nil {
					elog.Panic("set SystemSheet body row error", l.E(err), l.S("sheetName", sheetName))
				}
			}
		}
	}

	fileName := fmt.Sprintf("test-result/files/%s/sdk_test_all_%s.xlsx", time.Now().Format("2006_01_02"), time.Now().Format("2006_01_02_15_04_05"))
	if err := exFile.SaveAs(context.Background(), fileName); err != nil {
		elog.Panic("save excel error", l.E(err))
	}
}

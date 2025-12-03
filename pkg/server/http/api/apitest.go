package api

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/ego/core/elog"
	sdk "github.com/shimo-open/sdk-kit-go"
	"k8s.io/apimachinery/pkg/util/uuid"

	"sdk-demo-go/cmd/sdkctl"
	"sdk-demo-go/pkg/consts"
	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/models/db"
)

var createErr error

func GetAllApiTest(c *gin.Context) {
	testType, ok := c.GetQuery("type")
	if !ok {
		c.JSON(400, gin.H{"message": "type is required"})
		return
	}
	taskId := fmt.Sprintf("%s", uuid.NewUUID())
	resp := struct {
		TaskId string `json:"taskId"`
	}{
		TaskId: taskId,
	}
	c.JSON(200, resp)
	go func() {
		createErr = nil
		res := sdkctl.TestAll(context.Background(), testType)
		// Insert data after the tests finish
		createErr = createTestApi(res, taskId)
		if createErr != nil {
			elog.Error("Create test result", l.S("error: ", createErr.Error()))
			handleDBError(c, createErr)
			return
		}
	}()
}

func CheckETestProgress(c *gin.Context) {
	taskId, ok := c.GetQuery("taskId")
	if !ok {
		c.JSON(400, gin.H{"message": "taskId is undefined"})
		return
	}

	if createErr != nil {
		c.JSON(500, gin.H{
			"status":  consts.TestError,
			"message": createErr,
		})
		return
	}

	testApiList, err := db.GetTestApi(invoker.DB, taskId)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  consts.TestError,
			"message": err,
		})
		return
	}
	if len(testApiList) > 0 {
		resultData := make(map[string][]db.TestApi)
		for _, v := range testApiList {
			resultData[v.TestType] = append(resultData[v.TestType], v)
		}
		c.JSON(200, gin.H{
			"status": consts.TestFinish,
			"result": resultData,
		})
	} else {
		c.JSON(200, gin.H{
			"status":   consts.TestProcessing,
			"progress": sdkctl.TestProgress,
		})
	}
}

// Create test results
func createTestApi(data consts.AllApiTestRes, taskId string) (err error) {
	testTata := handleTestResult(data)
	testApis := make([]db.TestApi, 0)
	for key, value := range testTata {
		for _, v := range value {
			success := 0
			if v.Success {
				success = 1
			} else {
				success = 0
			}
			testApis = append(testApis, db.TestApi{
				TestId:        taskId,
				TestType:      key,
				Success:       success,
				ApiName:       v.ApiName,
				HttpCode:      v.HttpCode,
				HttpResp:      v.HttpResp,
				ErrMsg:        v.ErrMsg,
				PathStr:       v.PathStr,
				BodyReq:       v.BodyReq,
				Query:         v.Query,
				FormData:      v.FormData,
				FileExt:       v.FileExt,
				StartTime:     v.StartTime,
				TimeConsuming: v.TimeConsuming,
			})
		}
	}
	err = db.CreateTestApi(invoker.DB, &testApis)
	if err != nil {
		return err
	}
	return
}

// Retrieve the latest test data
func GetLatestTest(c *gin.Context) {
	result, err := db.GetLatestTest(invoker.DB)
	if err != nil {
		handleDBError(c, err)
		return
	}
	resultData := make(map[string][]db.TestApi)
	for _, v := range result {
		resultData[v.TestType] = append(resultData[v.TestType], v)
	}
	c.JSON(200, gin.H{
		"result": resultData,
	})
}

// Handle the returned test data
func handleTestResult(handleResult consts.AllApiTestRes) map[string][]consts.SingleApiTestRes {
	testResult := map[string][]consts.SingleApiTestRes{}
	for k, _ := range consts.SheetNameMap {
		testResult[k] = make([]consts.SingleApiTestRes, 0)
		switch k {
		case "BaseTestResMap":
			for key, _ := range handleResult.BaseTestResMap {
				if key != sdk.FileTypeTable {
					// Application spreadsheets do not expose these two APIs
					testResult[k] = append(testResult[k], handleResult.BaseTestResMap[key].GetPlainTextRes)
					testResult[k] = append(testResult[k], handleResult.BaseTestResMap[key].GetPlainTextWordCountRes)
				}
				if key != sdk.FileTypeSlide && key != sdk.FileTypeTable {
					// PPT and application spreadsheets do not provide the get-at API
					testResult[k] = append(testResult[k], handleResult.BaseTestResMap[key].GetMentionAtListRes)
				}
				testResult[k] = append(testResult[k], handleResult.BaseTestResMap[key].CreateFileRes)
				testResult[k] = append(testResult[k], handleResult.BaseTestResMap[key].CreateCopyRes)
				testResult[k] = append(testResult[k], handleResult.BaseTestResMap[key].DeleteFileRes)
				testResult[k] = append(testResult[k], handleResult.BaseTestResMap[key].CreatePreviewRes)
				testResult[k] = append(testResult[k], handleResult.BaseTestResMap[key].GetPreviewRes)
				testResult[k] = append(testResult[k], handleResult.BaseTestResMap[key].GetHistoryListRes)
				testResult[k] = append(testResult[k], handleResult.BaseTestResMap[key].GetRevisionListRes)

			}
		case "FileIOResMap":
			for _, statusMap := range handleResult.FileIOResMap {
				for ext, status := range statusMap.ImportFileRes {
					for key := range status {
						testResult[k] = append(testResult[k], consts.SingleApiTestRes{
							ApiName:       "导入文件" + "-" + string(ext),
							FileExt:       key,
							Success:       status[key].Success,
							PathStr:       status[key].PathStr,
							FormData:      status[key].FormData,
							StartTime:     status[key].StartTime,
							TimeConsuming: status[key].TimeConsuming,
							HttpCode:      status[key].HttpCode,
							ErrMsg:        status[key].ErrMsg,
						})
					}
				}
				for extE, statusE := range statusMap.ExportFileRes {
					for keyE := range statusE {
						testResult[k] = append(testResult[k], consts.SingleApiTestRes{
							ApiName:       "导出文件" + "-" + string(extE),
							FileExt:       keyE,
							Success:       statusE[keyE].Success,
							PathStr:       statusE[keyE].PathStr,
							BodyReq:       statusE[keyE].BodyReq,
							StartTime:     statusE[keyE].StartTime,
							TimeConsuming: statusE[keyE].TimeConsuming,
							HttpCode:      statusE[keyE].HttpCode,
							ErrMsg:        statusE[keyE].ErrMsg,
						})
					}
				}
			}
		case "SpreadsheetRes":
			testResult[k] = append(testResult[k], handleResult.SpreadsheetRes.GetTableContentRes)
			testResult[k] = append(testResult[k], handleResult.SpreadsheetRes.UpdateTableContentRes)
			testResult[k] = append(testResult[k], handleResult.SpreadsheetRes.AppendTableContentRes)
			testResult[k] = append(testResult[k], handleResult.SpreadsheetRes.AddTableSheetRes)
			testResult[k] = append(testResult[k], handleResult.SpreadsheetRes.DeleteTableRowRes)
			testResult[k] = append(testResult[k], handleResult.SpreadsheetRes.GetCommentCountRes)

		case "DocumentProRes":
			testResult[k] = append(testResult[k], handleResult.DocumentProRes.ReadBookmarkContentRes)
			testResult[k] = append(testResult[k], handleResult.DocumentProRes.ReplaceBookmarkContentRes)

		case "TableRes":
			testResult[k] = append(testResult[k], handleResult.TableRes.SheetExportToExcelRes)

		case "SystemRes":
			testResult[k] = append(testResult[k], handleResult.SystemRes.GetAppDetailRes)
			// testResult[k] = append(testResult[k], handleResult.SystemRes.UpdateCallbackUrlRes)
			testResult[k] = append(testResult[k], handleResult.SystemRes.ActivateUserSeatRes)
			testResult[k] = append(testResult[k], handleResult.SystemRes.CancelUserSeatRes)
			testResult[k] = append(testResult[k], handleResult.SystemRes.GetUserListAndSeatStatusRes)
			testResult[k] = append(testResult[k], handleResult.SystemRes.BatchSetUserSeatRes)
		}
	}
	return testResult
}

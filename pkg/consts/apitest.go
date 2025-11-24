package consts

// TestStatus represents the status of a test execution
type TestStatus int

const (
	// TestFinish indicates the test has finished
	TestFinish TestStatus = 0
	// TestProcessing indicates the test is in progress
	TestProcessing TestStatus = 1
	// TestError indicates the test has failed
	TestError TestStatus = 2
)

func (t TestStatus) int() int { return int(t) }

// AllApiTestRes holds all API test results organized by functionality
type AllApiTestRes struct {
	// Basic functionality shared by every suite
	BaseTestResMap map[FileType]BaseTestRes `json:"baseTestResMap"`
	// Import/export functionality shared by every suite
	FileIOResMap map[FileType]FileIORes `json:"fileIOResMap"`
	// Spreadsheet-specific functionality
	SpreadsheetRes SpreadsheetRes `json:"spreadsheetRes"`
	// Document Pro-specific functionality
	DocumentProRes DocumentProRes `json:"documentProRes"`
	// Application table-specific functionality
	TableRes TableRes `json:"tableRes"`
	// System functionality
	SystemRes SystemRes `json:"systemRes"`
}

// BaseTestRes holds test results for basic file operations
type BaseTestRes struct {
	// Create document
	CreateFileRes SingleApiTestRes `json:"createFileRes"`
	// Create copy
	CreateCopyRes SingleApiTestRes `json:"createCopyRes"`
	// Delete document
	DeleteFileRes SingleApiTestRes `json:"deleteFileRes"`
	// Create preview
	CreatePreviewRes SingleApiTestRes `json:"createPreviewRes"`
	// Open preview
	GetPreviewRes SingleApiTestRes `json:"getPreviewRes"`
	// Get history list
	GetHistoryListRes SingleApiTestRes `json:"getHistoryListRes"`
	// Get revision list
	GetRevisionListRes SingleApiTestRes `json:"getVersionListRes"`
	// Get plain text
	GetPlainTextRes SingleApiTestRes `json:"getPlainTextRes"`
	// Get plain-text word count
	GetPlainTextWordCountRes SingleApiTestRes `json:"getPlainTextWordCountRes"`
	// Get mention list
	GetMentionAtListRes SingleApiTestRes `json:"getAtListRes"`
}

// FileIORes holds test results for file import and export operations
type FileIORes struct {
	// Import files
	// Import for each supported format
	ImportFileRes map[FileType]map[string]ImportFileRes `json:"importFileRes"`
	// Export files
	// Export for each supported format
	ExportFileRes map[FileType]map[string]ExportFileRes `json:"exportFileRes"`
}

// SpreadsheetRes holds test results for spreadsheet-specific operations
type SpreadsheetRes struct {
	// Get spreadsheet content
	GetTableContentRes SingleApiTestRes `json:"getTableContentRes"`
	// Update spreadsheet content
	UpdateTableContentRes SingleApiTestRes `json:"updateTableContentRes"`
	// Append spreadsheet content
	AppendTableContentRes SingleApiTestRes `json:"appendTableContentRes"`
	// Delete spreadsheet rows
	DeleteTableRowRes SingleApiTestRes `json:"deleteTableRowRes"`
	// Add spreadsheet sheet
	AddTableSheetRes SingleApiTestRes `json:"addTableSheetRes"`
	// Get comment count
	GetCommentCountRes SingleApiTestRes `json:"getCommentCountRes"`
}

// DocumentProRes holds test results for Document Pro-specific operations
type DocumentProRes struct {
	// Read bookmark content
	ReadBookmarkContentRes SingleApiTestRes `json:"readBookmarkContentRes"`
	// Replace bookmark content
	ReplaceBookmarkContentRes SingleApiTestRes `json:"replaceBookmarkContentRes"`
}

// TableRes holds test results for application table-specific operations
type TableRes struct {
	// Export as Excel
	SheetExportToExcelRes SingleApiTestRes `json:"exportToExcelRes"`
}

// SystemRes holds test results for system-level operations
type SystemRes struct {
	// Get app details
	GetAppDetailRes SingleApiTestRes `json:"getAppDetailRes"`
	// // Update callback URL
	// UpdateCallbackUrlRes SingleApiTestRes `json:"updateCallbackUrlRes"`
	// Get users and seat status
	GetUserListAndSeatStatusRes SingleApiTestRes `json:"getUserListAndSeatStatusRes"`
	// Activate user seats
	ActivateUserSeatRes SingleApiTestRes `json:"activateUserSeatRes"`
	// Deactivate user seats
	CancelUserSeatRes SingleApiTestRes `json:"cancelUserSeatRes"`
	// Batch update user seats
	BatchSetUserSeatRes SingleApiTestRes `json:"batchSetUserSeatRes"`
}

// SingleApiTestRes represents the result of a single API test
type SingleApiTestRes struct {
	ApiName       string `json:"apiName" excel:"接口名字" excel_width:"25"`
	Success       bool   `json:"success" excel:"是否成功" excel_width:"15"`
	HttpCode      int    `json:"httpCode" excel:"状态码" excel_width:"10"`
	HttpResp      string `json:"httpResp" excel:"返回结果" excel_width:"25"`
	ErrMsg        string `json:"errMsg" excel:"错误信息" excel_width:"25"`
	PathStr       string `json:"pathStr" excel:"请求地址" excel_width:"25"`
	BodyReq       string `json:"bodyReq" excel:"body传参" excel_width:"25"`
	FormData      string `json:"formData" excel:"formData传参" excel_width:"25"`
	Query         string `json:"query" excel:"地址传参" excel_width:"25"`
	FileExt       string `json:"fileExt" excel:"传入文件后缀/导出文件类型" excel_width:"25"`
	TimeConsuming string `json:"timeConsuming" excel:"耗时" excel_width:"25"`
	StartTime     int64  `json:"startTime" excel:"测试开始时间" excel_width:"25"`
}

// ImportFileRes represents the result of a file import operation
type ImportFileRes struct {
	Success       bool   `json:"success"`
	HttpCode      int    `json:"httpCode"`
	ErrMsg        string `json:"errMsg"`
	PathStr       string `json:"pathStr"`
	FormData      string `json:"formData"`
	TimeConsuming string `json:"timeConsuming"`
	StartTime     int64  `json:"startTime"`
}

// ExportFileRes represents the result of a file export operation
type ExportFileRes struct {
	Success       bool   `json:"success"`
	HttpCode      int    `json:"httpCode"`
	ErrMsg        string `json:"errMsg"`
	PathStr       string `json:"pathStr"`
	BodyReq       string `json:"bodyReq"`
	TimeConsuming string `json:"timeConsuming"`
	StartTime     int64  `json:"startTime"`
}

// ExportBodyReq represents the request body for file export
type ExportBodyReq struct {
	Type FileType `json:"type"`
}

// FileFormData represents multipart form data for file operations
type FileFormData struct {
	FilePath string   `json:"filePath"`
	FileId   string   `json:"fileId"`
	FileType FileType `json:"fileType"`
}

var SheetNameMap = map[string]string{
	"BaseTestResMap": "基础",
	"FileIOResMap":   "导入导出",
	"SpreadsheetRes": "表格",
	"DocumentProRes": "传统文档",
	"TableRes":       "应用表格",
	"SystemRes":      "系统",
}

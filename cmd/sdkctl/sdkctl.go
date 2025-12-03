package sdkctl

import (
	"context"
	"encoding/json"

	"github.com/ego-component/egorm"
	sdk "github.com/shimo-open/sdk-kit-go"
	sdkapi "github.com/shimo-open/sdk-kit-go/api"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	_ "go.uber.org/automaxprocs"

	"sdk-demo-go/cmd"
	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/utils"
)

var (
	fileId       string
	fileName     string
	fileType     string
	shimoType    string
	lang         string
	acLang       string
	filePath     string
	url          string
	taskId       string
	exportType   string
	rg           string
	values2D     []string
	sheetName    string
	index        int
	bookmark     []string
	rbsStr       string
	appId        string
	page         int
	size         int
	userIds      []string
	status       int
	from         string
	to           string
	appIdQuery   string
	userId       int64
	TestProgress float64 // Test progress
)

var SdkCtl = &cobra.Command{
	Use:   "sdk-ctl",
	Short: "sdk demo command-line tool",
	Long:  `sdk demo command-line tool`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// ShimoSDK = shimo.InitTest()
		invoker.DB = egorm.Load("mysql").Build()
		initParams()
		invoker.InitShimo()
	},
}

var ApiTestCtl = &cobra.Command{
	Use:   "api-test",
	Short: "sdk api interface testing",
	Long:  `sdk api interface testing`,
}

var ApiTestCreatePreview = &cobra.Command{
	Use:   "create-preview",
	Short: "Create preview document",
	Long:  `Create preview document, accepts one parameter: the guid of the file to preview`,
	Run: func(cmd *cobra.Command, args []string) {
		TestCreatePreview(context.Background(), fileId)
	},
}

var ApiTestCreateFile = &cobra.Command{
	Use:   "create-file",
	Short: "Create file",
	Long:  `Create file, accepts 4 parameters: file name, file type, Shimo document type, Language and Accept-Language, where file type and Shimo document type are mutually exclusive`,
	Run: func(cmd *cobra.Command, args []string) {
		TestCreateFile(context.Background(), fileName, fileType, shimoType, lang, acLang)
	},
}

var ApiTestCreateFileCopy = &cobra.Command{
	Use:   "create-file-copy",
	Short: "Create file copy",
	Long:  `Create a copy of the file, accepts 1 parameter: source file guid`,
	Run: func(cmd *cobra.Command, args []string) {
		TestCreateFileCopy(context.Background(), fileId)
	},
}

var ApiTestDeleteFile = &cobra.Command{
	Use:   "delete-file",
	Short: "Delete file",
	Long:  `Delete file, accepts 1 parameter: file guid`,
	Run: func(cmd *cobra.Command, args []string) {
		TestDeleteFile(context.Background(), fileId)
	},
}

var ApiTestImportFile = &cobra.Command{
	Use:   "import-file",
	Short: "Upload file",
	Long:  `Upload file, accepts 3 parameters: absolute file path, specified file name, and collaborative document type`,
	Run: func(cmd *cobra.Command, args []string) {
		TestImportFile(context.Background(), filePath, fileName, shimoType)
	},
}

var ApiTestImportFileByUrl = &cobra.Command{
	Use:   "import-file-by-url",
	Short: "Upload file by URL",
	Long:  `Upload file by URL, accepts 4 parameters: file URL, specified file name, file type, and collaborative document type, where file type and collaborative document type are mutually exclusive`,
	Run: func(cmd *cobra.Command, args []string) {
		TestImportFileByUrl(context.Background(), url, fileName, fileType, shimoType)
	},
}

var ApiTestCheckImportFile = &cobra.Command{
	Use:   "check-import-file",
	Short: "Check file upload progress",
	Long:  `Check file upload progress, accepts 1 parameter: taskId`,
	Run: func(cmd *cobra.Command, args []string) {
		TestImportFileProgress(context.Background(), taskId)
	},
}

var ApiTestExportFile = &cobra.Command{
	Use:   "export-file",
	Short: "Export file",
	Long:  `Export file, accepts 2 parameters: file guid and export type`,
	Run: func(cmd *cobra.Command, args []string) {
		TestExportFile(context.Background(), fileId, exportType)
	},
}

var ApiTestCheckExportFile = &cobra.Command{
	Use:   "check-export-file",
	Short: "Check file export progress",
	Long:  `Check file export progress, accepts 1 parameter: taskId`,
	Run: func(cmd *cobra.Command, args []string) {
		TestExportFileProgress(context.Background(), taskId)
	},
}

var ApiTestExportTableSheets = &cobra.Command{
	Use:   "export-table-sheets",
	Short: "Export application table as Excel",
	Long:  `Export application table as .xlsx file and return the download URL.`,
	Run: func(cmd *cobra.Command, args []string) {
		TestExportTableAsSheets(context.Background(), fileId)
	},
}

var ApiTestGetFilePlainText = &cobra.Command{
	Use:   "get-file-plain-text",
	Short: "Get file plain text content",
	Long:  `Get file plain text content, accepts 1 parameter: file guid`,
	Run: func(cmd *cobra.Command, args []string) {
		TestGetFilePlainText(context.Background(), fileId)
	},
}

var ApiTestGetFilePlainTextWC = &cobra.Command{
	Use:   "get-file-plain-text-WC",
	Short: "Get file plain text word count",
	Long:  `Get file plain text word count, accepts 1 parameter: file guid`,
	Run: func(cmd *cobra.Command, args []string) {
		TestGetFilePlainTextWC(context.Background(), fileId)
	},
}

var ApiTestGetDocSideBarInfo = &cobra.Command{
	Use:   "get-doc-sidebar-info",
	Short: "Get file sidebar history list",
	Long:  "Get file sidebar history list, history types include operation history and edit history",
	Run: func(cmd *cobra.Command, args []string) {
		TestGetDocSidebarInfo(context.Background(), fileId, page, size)
	},
}

var ApiTestGetRevision = &cobra.Command{
	Use:   "get-revision",
	Short: "Get revision list",
	Long:  "Get file revision list information",
	Run: func(cmd *cobra.Command, args []string) {
		TestGetRevision(context.Background(), fileId)
	},
}

var ApiTestGetExcelContent = &cobra.Command{
	Use:   "get-excel-content",
	Short: "Get Excel file content",
	Long:  `Get Excel file content, accepts 2 parameters: file ID and query range`,

	Run: func(cmd *cobra.Command, args []string) {
		TestGetExcelContent(context.Background(), fileId, rg, nil)
	},
}

var ApiTestUpdateExcelContent = &cobra.Command{
	Use:   "update-excel-content",
	Short: "Update Excel file content",
	Long:  `Update Excel file content, accepts 3 parameters: file guid, update range, and content to modify`,
	Run: func(cmd *cobra.Command, args []string) {
		TestUpdateExcelContent(context.Background(), fileId, rg, values2D)
	},
}

var ApiTestAppendExcelContent = &cobra.Command{
	Use:   "append-excel-content",
	Short: "Append Excel file content",
	Long:  `Append Excel file content, accepts 3 parameters: file guid, append range, and content to add`,
	Run: func(cmd *cobra.Command, args []string) {
		TestAppendExcelContent(context.Background(), fileId, rg, values2D)
	},
}

var ApiTestDeleteExcelRows = &cobra.Command{
	Use:   "delete-excel-rows",
	Short: "Delete Excel file row",
	Long:  `Delete Excel file row, accepts 3 parameters: file guid, sheet name, and row number`,
	Run: func(cmd *cobra.Command, args []string) {
		TestDeleteExcelRows(context.Background(), fileId, sheetName, index)
	},
}

var ApiTestCreateExcelSheet = &cobra.Command{
	Use:   "create-excel-sheet",
	Short: "Create Excel sheet",
	Long:  `Create Excel sheet, accepts 1 parameter: sheet name`,

	Run: func(cmd *cobra.Command, args []string) {
		TestCreateExcelSheet(context.Background(), fileId, sheetName)
	},
}

var ApiTestGetDocProBookmark = &cobra.Command{
	Use:   "get-doc-pro-bookmark",
	Short: "Get document bookmarks",
	Long:  `Get document bookmarks, accepts 2 parameters: file guid and bookmark list`,
	Run: func(cmd *cobra.Command, args []string) {
		TestGetDocProBookmark(context.Background(), fileId, bookmark)
	},
}

var ApiTestReplaceDocProBookmark = &cobra.Command{
	Use:   "replace-doc-pro-bookmark [fieldId] [replaceBookmarks]",
	Short: "Replace document bookmarks",
	Long:  `Replace document bookmarks, accepts 2 parameters: file guid and bookmark list`,
	Run: func(cmd *cobra.Command, args []string) {
		req := sdkapi.RepBookmarkContentReqBody{}
		err := json.Unmarshal([]byte(rbsStr), &req.Replacements)
		if err != nil {
			errHandler(err)
			return
		}
		TestReplaceDocProBookmark(context.Background(), fileId, req)
	},
}

var ApiTestGetAppDetails = &cobra.Command{
	Use:   "get-app-details",
	Short: "Get APP details",
	Long:  "Get APP details, accepts 1 parameter: appId",
	Run: func(cmd *cobra.Command, args []string) {
		TestGetAppDetails(context.Background(), appId)
	},
}

var ApiTestUpdateAppEndPoint = &cobra.Command{
	Use:   "update-app-endpoint",
	Short: "Update APP callback endpoint",
	Long:  "Update APP callback endpoint, accepts 2 parameters: appId and new callback URL",
	Run: func(cmd *cobra.Command, args []string) {
		TestUpdateAppEndpoint(context.Background(), appId, url)
	},
}

var ApiTestGetUsersWithStatus = &cobra.Command{
	Use:   "get-users-with-status",
	Short: "Get user list and seat status",
	Long:  "Get user list and seat status, accepts 2 parameters: page and size for pagination",
	Run: func(cmd *cobra.Command, args []string) {
		TestGetUsersWithStatus(context.Background(), page, size)
	},
}

var ApiTestActivateUsers = &cobra.Command{
	Use:   "activate-users",
	Short: "Activate user seats",
	Long:  "Activate user seats, accepts 1 parameter: user ID array",
	Run: func(cmd *cobra.Command, args []string) {
		TestActivateUsers(context.Background(), userIds)
	},
}

var ApiTestDeactivateUsers = &cobra.Command{
	Use:   "deactivate-users",
	Short: "Deactivate user seats",
	Long:  "Deactivate user seats, accepts 1 parameter: user ID array",
	Run: func(cmd *cobra.Command, args []string) {
		TestDeactivateUsers(context.Background(), userIds)
	},
}

var ApiTestBatchSetUsersStatus = &cobra.Command{
	Use:   "batch-set-users-status",
	Short: "Batch set user seat status",
	Long:  "Batch set user seat status, accepts 2 parameters: user ID array and seat status",
	Run: func(cmd *cobra.Command, args []string) {
		TestBatchSetUsersStatus(context.Background(), userIds, status)
	},
}

var ApiTestBatch = &cobra.Command{
	Use:   "batch-test",
	Short: "Batch API testing",
	Long:  "Batch API testing",
}

var ApiTestBatch_Base = &cobra.Command{
	Use:   "base",
	Short: "Basic functionality testing",
	Long:  "Basic functionality testing",
	Run: func(cmd *cobra.Command, args []string) {
		TestBase(context.Background(), []sdk.FileType{sdk.GetFileType(getFileTypeStr(args))})
	},
}

var ApiTestBatch_Common = &cobra.Command{
	Use:   "common",
	Short: "Common functionality testing",
	Long:  "Common functionality testing, includes base + import/export",
	Run: func(cmd *cobra.Command, args []string) {
		fileType := "document"
		if len(args) > 0 {
			fileType = cast.ToString(args[0])
		}
		TestCommon(context.Background(), fileType)
	},
}

var ApiTestBatch_System = &cobra.Command{
	Use:   "system",
	Short: "System functionality testing",
	Long:  "System functionality testing, including app info, update callback endpoint, update user status, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		TestSystem(context.Background())
	}}

var ApiTestBatch_All = &cobra.Command{
	Use:   "all",
	Short: "Full functionality testing",
	Long:  "Full functionality testing, includes base + import/export + spreadsheet and document operations",
	Run: func(cmd *cobra.Command, args []string) {
		res := TestAll(context.Background(), getFileTypeStr(args))
		utils.TestSaveExcel(res)
	},
}

func init() {
	SdkCtl.InheritedFlags()
	ApiTestCtl.InheritedFlags()
	cmd.RootCommand.AddCommand(SdkCtl)

	SdkCtl.AddCommand(ApiTestCtl)

	// api-test
	ApiTestCreatePreview.Flags().StringVarP(&fileId, "fileId", "f", "", "File guid (required)")
	ApiTestCreatePreview.MarkFlagRequired("fileId")
	ApiTestCtl.AddCommand(ApiTestCreatePreview)

	ApiTestCreateFile.Flags().StringVar(&fileName, "fileName", "", "File name")
	ApiTestCreateFile.Flags().StringVar(&fileType, "fileType", "", "File type")
	ApiTestCreateFile.Flags().StringVar(&shimoType, "shimoType", "", "Collaborative document type, available values: document, spreadsheet, documentPro, presentation, table")
	ApiTestCreateFile.Flags().StringVar(&lang, "lang", "", "Header Lang field")
	ApiTestCreateFile.Flags().StringVar(&acLang, "acLang", "", "Header Accept-Language field")
	ApiTestCtl.AddCommand(ApiTestCreateFile)

	ApiTestCreateFileCopy.Flags().StringVarP(&fileId, "fileId", "f", "", "File guid (required)")
	ApiTestCreateFileCopy.MarkFlagRequired("fileId")
	ApiTestCtl.AddCommand(ApiTestCreateFileCopy)

	ApiTestDeleteFile.Flags().StringVarP(&fileId, "fileId", "f", "", "File guid (required)")
	ApiTestCtl.AddCommand(ApiTestDeleteFile)

	ApiTestImportFile.Flags().StringVar(&filePath, "filePath", "", "Local file absolute path")
	ApiTestImportFile.Flags().StringVar(&fileName, "fileName", "", "Specified file name, must include correct extension")
	ApiTestImportFile.Flags().StringVar(&shimoType, "shimoType", "", "Collaborative document type, available values: document, spreadsheet, documentPro, presentation, table")
	ApiTestCtl.AddCommand(ApiTestImportFile)

	ApiTestImportFileByUrl.Flags().StringVar(&url, "url", "", "File URL")
	ApiTestImportFileByUrl.Flags().StringVar(&fileName, "fileName", "", "Specified file name")
	ApiTestImportFileByUrl.Flags().StringVar(&fileType, "fileType", "", "File type when not creating collaborative document")
	ApiTestImportFileByUrl.Flags().StringVar(&shimoType, "shimoType", "", "File type when creating collaborative document, available values: document, spreadsheet, documentPro, presentation, table")
	ApiTestCtl.AddCommand(ApiTestImportFileByUrl)

	ApiTestCheckImportFile.Flags().StringVar(&taskId, "taskId", "", "TaskId returned from file upload task")
	ApiTestCtl.AddCommand(ApiTestCheckImportFile)

	ApiTestExportFile.Flags().StringVarP(&fileId, "fileId", "f", "", "File guid")
	ApiTestExportFile.Flags().StringVar(&exportType, "exportType", "", "File export type")
	ApiTestCtl.AddCommand(ApiTestExportFile)

	ApiTestCheckExportFile.Flags().StringVar(&taskId, "taskId", "", "TaskId of file export task")
	ApiTestCtl.AddCommand(ApiTestCheckExportFile)

	ApiTestExportTableSheets.Flags().StringVarP(&fileId, "fileId", "f", "", "File guid")
	ApiTestCtl.AddCommand(ApiTestExportTableSheets)

	ApiTestGetFilePlainText.Flags().StringVarP(&fileId, "fileId", "f", "", "File guid")
	ApiTestCtl.AddCommand(ApiTestGetFilePlainText)

	ApiTestGetFilePlainTextWC.Flags().StringVarP(&fileId, "fileId", "f", "", "File guid")
	ApiTestCtl.AddCommand(ApiTestGetFilePlainTextWC)

	ApiTestGetDocSideBarInfo.Flags().StringVarP(&fileId, "fileId", "f", "", "File guid")
	ApiTestGetDocSideBarInfo.Flags().IntVar(&page, "page", 1, "Page number")
	ApiTestGetDocSideBarInfo.Flags().IntVar(&size, "size", 10, "Page size")
	ApiTestCtl.AddCommand(ApiTestGetDocSideBarInfo)

	ApiTestGetRevision.Flags().StringVarP(&fileId, "fileId", "f", "", "File guid")
	ApiTestCtl.AddCommand(ApiTestGetRevision)

	ApiTestGetExcelContent.Flags().StringVarP(&fileId, "fileId", "f", "", "File guid")
	ApiTestGetExcelContent.Flags().StringVar(&rg, "rg", "", "Query range, format: Sheet!{CellRange}. Format specification:\nCell range format is A1:C10 or A1, starting position cannot exceed maximum rows/columns\nIf sheet name contains special characters !:' they must be wrapped in ''")
	ApiTestCtl.AddCommand(ApiTestGetExcelContent)

	ApiTestUpdateExcelContent.Flags().StringVarP(&fileId, "fileId", "f", "", "File guid")
	ApiTestUpdateExcelContent.Flags().StringVar(&rg, "rg", "", "Query range, format: Sheet!{CellRange}. Format specification:\nCell range format is A1:C10 or A1, starting position cannot exceed maximum rows/columns\nIf sheet name contains special characters !:' they must be wrapped in ''")
	ApiTestUpdateExcelContent.Flags().StringSliceVar(&values2D, "value", []string{}, "Update values, each row in format [a|b|c...]")
	ApiTestCtl.AddCommand(ApiTestUpdateExcelContent)

	ApiTestAppendExcelContent.Flags().StringVarP(&fileId, "fileId", "f", "", "File guid")
	ApiTestAppendExcelContent.Flags().StringVar(&rg, "rg", "", "Query range, format: Sheet!{CellRange}. Format specification:\nCell range format is A1:C10 or A1, starting position cannot exceed maximum rows/columns\nIf sheet name contains special characters !:' they must be wrapped in ''")
	ApiTestAppendExcelContent.Flags().StringSliceVar(&values2D, "value", []string{}, "Update values, each row in format [a|b|c...]")
	ApiTestCtl.AddCommand(ApiTestAppendExcelContent)

	ApiTestDeleteExcelRows.Flags().StringVarP(&fileId, "fileId", "f", "", "File guid")
	ApiTestDeleteExcelRows.Flags().StringVar(&sheetName, "sheetName", "", "Sheet name")
	ApiTestDeleteExcelRows.Flags().IntVar(&index, "index", 10, "Row number")
	ApiTestCtl.AddCommand(ApiTestDeleteExcelRows)

	ApiTestCreateExcelSheet.Flags().StringVarP(&fileId, "fileId", "f", "", "File guid")
	ApiTestCreateExcelSheet.Flags().StringVar(&sheetName, "sheetName", "", "Sheet name")
	ApiTestCtl.AddCommand(ApiTestCreateExcelSheet)

	ApiTestGetDocProBookmark.Flags().StringVarP(&fileId, "fileId", "f", "", "File guid")
	ApiTestGetDocProBookmark.Flags().StringSliceVar(&bookmark, "bookmark", []string{}, "Bookmark list")
	ApiTestCtl.AddCommand(ApiTestGetDocProBookmark)

	ApiTestReplaceDocProBookmark.Flags().StringVarP(&fileId, "fileId", "f", "", "File guid")
	ApiTestReplaceDocProBookmark.Flags().StringVar(&rbsStr, "rbs", "", "Document bookmark list to replace")
	ApiTestCtl.AddCommand(ApiTestReplaceDocProBookmark)

	ApiTestGetAppDetails.Flags().StringVarP(&appId, "appId", "a", "", "appId")
	ApiTestCtl.AddCommand(ApiTestGetAppDetails)

	ApiTestUpdateAppEndPoint.Flags().StringVarP(&appId, "appId", "a", "", "appId")
	ApiTestUpdateAppEndPoint.Flags().StringVar(&url, "url", "", "URL to update")
	ApiTestCtl.AddCommand(ApiTestUpdateAppEndPoint)

	ApiTestGetUsersWithStatus.Flags().IntVar(&page, "page", 1, "Page number, defaults to 1")
	ApiTestGetUsersWithStatus.Flags().IntVar(&size, "size", 10, "Number of items per page, defaults to 10")
	ApiTestCtl.AddCommand(ApiTestGetUsersWithStatus)

	ApiTestActivateUsers.Flags().StringSliceVar(&userIds, "userIds", userIds, "User IDs to activate")
	ApiTestCtl.AddCommand(ApiTestActivateUsers)

	ApiTestDeactivateUsers.Flags().StringSlice("userIds", userIds, "User IDs to deactivate")
	ApiTestCtl.AddCommand(ApiTestDeactivateUsers)

	ApiTestBatchSetUsersStatus.Flags().StringSliceVar(&userIds, "userIds", userIds, "User IDs to modify")
	ApiTestBatchSetUsersStatus.Flags().IntVar(&status, "status", 0, "1=activate, 0=disable, -1=not enabled")
	ApiTestCtl.AddCommand(ApiTestBatchSetUsersStatus)

	ApiTestBatch.Flags().String("fileType", "", "File type")
	ApiTestBatch.AddCommand(ApiTestBatch_Base)
	ApiTestBatch.AddCommand(ApiTestBatch_Common)
	ApiTestBatch.AddCommand(ApiTestBatch_System)
	ApiTestBatch.AddCommand(ApiTestBatch_All)

	ApiTestCtl.AddCommand(ApiTestBatch)
}

func initParams() {
	userId = getUserId()
}

func getFileTypeStr(args []string) (ft string) {
	ft = "document"
	if len(args) > 0 {
		ft = cast.ToString(args[0])
	}
	return
}

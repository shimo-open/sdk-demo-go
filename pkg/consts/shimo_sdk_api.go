package consts

import (
	"net/http"
)

type ShimoSdkApi string

type ShimoSdkPath string

func (s ShimoSdkApi) String() string {
	return string(s)
}

func (s ShimoSdkApi) Name() string {
	return ShimoSdkApiNameMap[s]
}

const (
	// System APIs

	// ShimoSdkApiGetAppDetails fetches app details
	ShimoSdkApiGetAppDetails ShimoSdkApi = "GetAppDetails"
	// ShimoSdkApiUpdateAppEndpoint updates the app callback URL
	ShimoSdkApiUpdateAppEndpoint ShimoSdkApi = "UpdateAppEndpoint"
	// ShimoSdkApiGetUsersWithStatus fetches user seat status
	ShimoSdkApiGetUsersWithStatus ShimoSdkApi = "GetUsersWithStatus"
	// ShimoSdkApiActivateUsers activates user seats
	ShimoSdkApiActivateUsers ShimoSdkApi = "ActivateUsers"
	// ShimoSdkApiDeactivateUsers deactivates user seats
	ShimoSdkApiDeactivateUsers ShimoSdkApi = "DeactivateUsers"
	// ShimoSdkApiBatchSetUsersStatus batch-updates user seats
	ShimoSdkApiBatchSetUsersStatus ShimoSdkApi = "BatchSetUsersStatus"

	// Feature APIs

	// ShimoSdkApiCreatePreview creates a preview
	ShimoSdkApiCreatePreview ShimoSdkApi = "CreatePreview"
	// ShimoSdkApiGetPreview opens a preview
	ShimoSdkApiGetPreview ShimoSdkApi = "GetPreview"
	// ShimoSdkApiGetPreviewDownload downloads the preview file
	ShimoSdkApiGetPreviewDownload ShimoSdkApi = "GetPreviewDownload"
	// ShimoSdkApiCreateFile creates a collaborative document
	ShimoSdkApiCreateFile ShimoSdkApi = "CreateFile"
	// ShimoSdkApiCreateFileCopy duplicates a collaborative document
	ShimoSdkApiCreateFileCopy ShimoSdkApi = "CreateFileCopy"
	// ShimoSdkApiDeleteFile deletes a collaborative document
	ShimoSdkApiDeleteFile ShimoSdkApi = "DeleteFile"
	// ShimoSdkApiImportFile imports a file
	ShimoSdkApiImportFile ShimoSdkApi = "ImportFile"
	// ShimoSdkApiImportFileByUrl imports a file via URL
	ShimoSdkApiImportFileByUrl ShimoSdkApi = "ImportFileByUrl"
	// ShimoSdkApiImportFileProgress checks import progress
	ShimoSdkApiImportFileProgress ShimoSdkApi = "ImportFileProgress"
	// ShimoSdkApiExportFile exports a file
	ShimoSdkApiExportFile ShimoSdkApi = "ExportFile"
	// ShimoSdkApiExportFileProgress checks export progress
	ShimoSdkApiExportFileProgress ShimoSdkApi = "ExportFileProgress"
	// ShimoSdkApiExportTableAsExcel exports a table to Excel
	ShimoSdkApiExportTableAsExcel ShimoSdkApi = "ExportTableAsExcel"
	// ShimoSdkApiGetFilePlainText gets file plain text
	ShimoSdkApiGetFilePlainText ShimoSdkApi = "GetFilePlainText"
	// ShimoSdkApiGetFilePlainTextWordCount gets plain-text word count
	ShimoSdkApiGetFilePlainTextWordCount ShimoSdkApi = "GetFilePlainTextWordCount"
	// ShimoSdkApiGetDocSidebarInfo gets sidebar data
	ShimoSdkApiGetDocSidebarInfo ShimoSdkApi = "GetDocSidebarInfo"
	// ShimoSdkApiGetRevision gets revision history
	ShimoSdkApiGetRevision ShimoSdkApi = "GetRevision"
	// ShimoSdkApiGetExcelContent gets Excel content
	ShimoSdkApiGetExcelContent ShimoSdkApi = "GetExcelContent"
	// ShimoSdkApiUpdateExcelContent updates Excel content
	ShimoSdkApiUpdateExcelContent ShimoSdkApi = "UpdateExcelContent"
	// ShimoSdkApiAppendExcelContent appends Excel content
	ShimoSdkApiAppendExcelContent ShimoSdkApi = "AppendExcelContent"
	// ShimoSdkApiDeleteExcelRows deletes Excel rows
	ShimoSdkApiDeleteExcelRows ShimoSdkApi = "DeleteExcelRows"
	// ShimoSdkApiCreateExcelSheet creates an Excel sheet
	ShimoSdkApiCreateExcelSheet ShimoSdkApi = "CreateExcelSheet"
	// ShimoSdkApiGetDocProBookmark gets a Document Pro bookmark
	ShimoSdkApiGetDocProBookmark ShimoSdkApi = "GetDocProBookmark"
	// ShimoSdkApiReplaceDocProBookmark replaces a Document Pro bookmark
	ShimoSdkApiReplaceDocProBookmark ShimoSdkApi = "ReplaceDocProBookmark"
	// ShimoSdkApiErrorCallback handles error callback
	ShimoSdkApiErrorCallback ShimoSdkApi = "ErrorCallback"
	// ShimoSdkApiGetMentionAt gets mention info
	ShimoSdkApiGetMentionAt ShimoSdkApi = "GetMentionAt"
	// ShimoSdkApiGetCommentCount gets comment count
	ShimoSdkApiGetCommentCount ShimoSdkApi = "GetCommentCount"
)

const (
	ShimoSdkPathImportFile         ShimoSdkPath = "/sdk/v2/api/files/v1/import"
	ShimoSdkPathImportFileProgress ShimoSdkPath = "/sdk/v2/api/files/v1/import/progress"
	ShimoSdkPathExportFile         ShimoSdkPath = "/sdk/v2/api/files/v1/export/"
	ShimoSdkPathExportFileProgress ShimoSdkPath = "/sdk/v2/api/files/v1/export/progress"
)

// ShimoSdkApiExpectCodesMap lists expected HTTP codes
var ShimoSdkApiExpectCodesMap = map[ShimoSdkApi]map[int]struct{}{
	ShimoSdkApiCreatePreview: {
		http.StatusOK: {},
	},
	ShimoSdkApiGetPreview: {
		http.StatusOK: {},
	},
	ShimoSdkApiGetPreviewDownload: {
		http.StatusOK: {},
	},
	ShimoSdkApiCreateFile: {
		http.StatusNoContent: {},
	},
	ShimoSdkApiCreateFileCopy: {
		http.StatusOK:        {},
		http.StatusNoContent: {},
	},
	ShimoSdkApiDeleteFile: {
		http.StatusNoContent: {},
	},
	ShimoSdkApiImportFile: {
		http.StatusOK: {},
	},
	ShimoSdkApiImportFileByUrl: {
		http.StatusOK: {},
	},
	ShimoSdkApiImportFileProgress: {
		http.StatusOK: {},
	},
	ShimoSdkApiExportFile: {
		http.StatusOK: {},
	},
	ShimoSdkApiExportFileProgress: {
		http.StatusOK: {},
	},
	ShimoSdkApiExportTableAsExcel: {
		http.StatusOK: {},
	},
	ShimoSdkApiGetFilePlainText: {
		http.StatusOK: {},
	},
	ShimoSdkApiGetFilePlainTextWordCount: {
		http.StatusOK: {},
	},
	ShimoSdkApiGetDocSidebarInfo: {
		http.StatusOK: {},
	},
	ShimoSdkApiGetRevision: {
		http.StatusOK: {},
	},
	ShimoSdkApiGetExcelContent: {
		http.StatusOK: {},
	},
	ShimoSdkApiUpdateExcelContent: {
		http.StatusNoContent: {},
	},
	ShimoSdkApiAppendExcelContent: {
		http.StatusNoContent: {},
	},
	ShimoSdkApiDeleteExcelRows: {
		http.StatusNoContent: {},
	},
	ShimoSdkApiCreateExcelSheet: {
		http.StatusNoContent: {},
	},
	ShimoSdkApiGetDocProBookmark: {
		http.StatusOK: {},
	},
	ShimoSdkApiReplaceDocProBookmark: {
		http.StatusNoContent: {},
	},
	ShimoSdkApiErrorCallback: {
		http.StatusOK: {},
	},
	ShimoSdkApiGetMentionAt: {
		http.StatusOK: {},
	},
	ShimoSdkApiGetCommentCount: {
		http.StatusOK: {},
	},
	ShimoSdkApiGetAppDetails: {
		http.StatusOK: {},
	},
	ShimoSdkApiUpdateAppEndpoint: {
		http.StatusNoContent: {},
	},
	ShimoSdkApiGetUsersWithStatus: {
		http.StatusOK: {},
	},
	ShimoSdkApiActivateUsers: {
		http.StatusNoContent: {},
	},
	ShimoSdkApiDeactivateUsers: {
		http.StatusNoContent: {},
	},
	ShimoSdkApiBatchSetUsersStatus: {
		http.StatusNoContent: {},
	},
}

var ShimoSdkApiNameMap = map[ShimoSdkApi]string{
	ShimoSdkApiCreatePreview:             "创建预览",
	ShimoSdkApiGetPreview:                "访问预览(后端接口)",
	ShimoSdkApiCreateFile:                "创建协同文档",
	ShimoSdkApiCreateFileCopy:            "创建协同文档副本",
	ShimoSdkApiDeleteFile:                "删除协同文档",
	ShimoSdkApiImportFile:                "导入文件",
	ShimoSdkApiImportFileByUrl:           "通过 url 导入文件",
	ShimoSdkApiImportFileProgress:        "导入文件进度",
	ShimoSdkApiExportFile:                "导出文件",
	ShimoSdkApiExportFileProgress:        "导出文件进度",
	ShimoSdkApiExportTableAsExcel:        "导出表格为 excel",
	ShimoSdkApiGetFilePlainText:          "获取文件纯文本",
	ShimoSdkApiGetFilePlainTextWordCount: "获取文件纯文本字数",
	ShimoSdkApiGetDocSidebarInfo:         "获取文档侧边栏信息",
	ShimoSdkApiGetRevision:               "获取文档历史版本",
	ShimoSdkApiGetExcelContent:           "获取 excel 内容",
	ShimoSdkApiUpdateExcelContent:        "更新 excel 内容",
	ShimoSdkApiAppendExcelContent:        "追加 excel 内容",
	ShimoSdkApiDeleteExcelRows:           "删除 excel 行",
	ShimoSdkApiCreateExcelSheet:          "创建 excel sheet",
	ShimoSdkApiGetDocProBookmark:         "获取传统文档书签",
	ShimoSdkApiReplaceDocProBookmark:     "替换传统文档书签",
	ShimoSdkApiErrorCallback:             "错误回调",
	ShimoSdkApiGetMentionAt:              "获取 @ 信息",
	ShimoSdkApiGetAppDetails:             "获取应用详情",
	ShimoSdkApiUpdateAppEndpoint:         "更新应用回调地址",
	ShimoSdkApiGetUsersWithStatus:        "获取用户状态",
	ShimoSdkApiActivateUsers:             "激活用户席位",
	ShimoSdkApiDeactivateUsers:           "取消用户席位",
	ShimoSdkApiBatchSetUsersStatus:       "批量设置用户席位",
	ShimoSdkApiGetCommentCount:           "获取评论数",
}

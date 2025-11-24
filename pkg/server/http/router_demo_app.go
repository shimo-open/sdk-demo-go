package http

import (
	"github.com/gotomicro/ego/server/egin"

	"sdk-demo-go/pkg/server/http/api"
	"sdk-demo-go/pkg/server/http/middlewares"
)

// registerDemoAppAPIs for the demo project's own APIs
func registerDemoAppAPIs(r *egin.Component) {
	apiGroup := r.Group("/api")
	apiGroup.GET("/sign", api.SignJWT)

	// app api
	apiAppGroup := apiGroup.Group("/apps", middlewares.UserAuthMiddleware)
	apiAppGroup.GET("/detail", api.GetAppDetails)
	apiAppGroup.PUT("/endpoint-url", api.PutEndpointUrl)

	// apiTest api
	appExcelGroup := apiGroup.Group("/apiTest", middlewares.UserAuthMiddleware)
	appExcelGroup.GET("/allApiTest", api.GetAllApiTest)
	appExcelGroup.GET("/testProgress", api.CheckETestProgress)
	appExcelGroup.GET("/getTestApiList", api.GetLatestTest)

	// team api
	apiTeamGroup := apiGroup.Group("/teams", middlewares.UserAuthMiddleware)
	apiTeamGroup.GET("/", api.GetTeams)
	apiTeamGroup.GET("", api.GetTeams)
	apiTeamGroup.GET("/:teamId/members", api.GetTeamMembers)
	apiTeamGroup.POST("/:teamId/members", api.JoinTeam)
	apiTeamGroup.PATCH("/:teamId/role/creator", api.TransferCreator)
	apiTeamGroup.POST("/", api.CreateTeam)
	apiTeamGroup.POST("", api.CreateTeam)
	apiTeamGroup.POST("/:teamId/departments", api.CreateDepartment)
	apiTeamGroup.GET("/:teamId/department-top-tree", api.GetDepartmentTopTree)
	apiTeamGroup.GET("/:teamId/departments/:deptId/members", api.GetDeptMembers)
	apiTeamGroup.GET("/:teamId/departments/:deptId/children", api.GetSubDeptsAndMembers)
	apiTeamGroup.PATCH("/:teamId/departments/:deptId/members", api.PatchMemberToDept)
	apiTeamGroup.POST("/:teamId/departments/:deptId/members", api.AddMemberToDept)
	apiTeamGroup.DELETE("/:teamId/departments/:deptId/members/:userId", api.DeleteMemberFromDept)
	apiTeamGroup.DELETE("/:teamId/departments/:deptId", api.DeleteDept)

	// file api
	apiFileGroup := apiGroup.Group("/files", middlewares.UserAuthMiddleware)
	apiFileGroup.GET("/", api.GetUserFiles)
	apiFileGroup.GET("", api.GetUserFiles)
	apiFileGroup.GET("/:fileGuid/thumbnail", api.GetFileThumbnail)
	apiFileGroup.GET("/:fileGuid/open", api.OpenFile)
	apiFileGroup.GET("/:fileGuid/download-plain-text", api.GetPlainText)
	apiFileGroup.GET("/:fileGuid/revisions", api.GetFileRevision)
	apiFileGroup.GET("/:fileGuid/comment-count", api.CountComments)
	apiFileGroup.GET("/:fileGuid/mention-at-list", api.GetMentionAtList)
	apiFileGroup.POST("/", api.CreateFile)
	apiFileGroup.POST("", api.CreateFile)
	apiFileGroup.POST("/upload", api.UploadFile)
	apiFileGroup.GET("/:fileGuid", api.GetFileInfo)
	apiFileGroup.POST("/import", api.ImportFile)
	apiFileGroup.POST("/import_by_url", api.ImportFileByUrl)
	apiFileGroup.POST("/import_by_url/progress", api.CheckImportUrlProgress)
	apiFileGroup.POST("/import/progress", api.CheckImportProgress)
	apiFileGroup.POST("/:fileGuid/export", api.ExportFile)
	apiFileGroup.POST("/:fileGuid/export/table-sheets", api.ExportTableSheets)
	apiFileGroup.POST("/:fileGuid/duplicate", api.DuplicateFile)
	apiFileGroup.POST("/export/progress", api.CheckExportFileProgress)
	apiFileGroup.DELETE("/:fileGuid", api.DeleteFile)
	apiFileGroup.DELETE("/batch/delete", api.BatchDeleteFile)
	apiFileGroup.PATCH("/:fileGuid", api.RenameFile)
	apiFileGroup.GET("/:fileGuid/collaborators", api.GetCollaborators)
	apiFileGroup.PATCH(":fileGuid/collaborators", api.UpdateCollaborators)
	apiFileGroup.GET("/:fileGuid/doc-sidebar-info", api.GetDocSidebar)
	apiFileGroup.POST("/importUrl", api.GetImportUrl)
	apiFileGroup.GET("/importUrl/redirect", api.GetImportRedirectUrl)

	// user api
	apiUserGroup := apiGroup.Group("/users", middlewares.UserAuthMiddleware)
	apiUserGroup.POST("/auth", api.Auth)
	apiUserGroup.POST("/signin", api.SignIn)
	apiUserGroup.POST("/signup", api.SignUp)
	apiUserGroup.GET("/:userId", api.GetUserById)
	apiUserGroup.GET("/:userId/teams", api.GetTeamsByUserId)
	apiUserGroup.DELETE("/me/teams/:teamId", api.DeleteMeFromTeam)
	apiUserGroup.GET("/", api.GetAllUsers)
	apiUserGroup.GET("", api.GetAllUsers)

	// event api
	apiEventGroup := apiGroup.Group("/events", middlewares.UserAuthMiddleware)
	apiEventGroup.GET("/", api.GetEvents)
	apiEventGroup.GET("", api.GetEvents)
	apiEventGroup.GET("/system-messages", api.GetSystemMessages)
	apiEventGroup.GET("/error_callback", api.ErrorCallback)

	// front inspect api
	apiFrontInspectGroup := apiGroup.Group("/internal", middlewares.FrontInspectAuthMiddleware)
	apiFrontInspectGroup.POST("", api.FrontInspectCreate)

	// AI knowledge base api
	apiKnowledgeBaseGroup := apiGroup.Group("/knowledge", middlewares.UserAuthMiddleware)
	apiKnowledgeBaseGroup.GET("/list", api.ListKnowledgeBases)
	apiKnowledgeBaseGroup.POST("/create", api.CreateKnowledgeBase)
	apiKnowledgeBaseGroup.DELETE("/:knowledgeBaseGuid/delete", api.DeleteKnowledgeBase)
	apiKnowledgeBaseGroup.POST("/import", api.ImportFileToKnowledgeBase)
	apiKnowledgeBaseGroup.POST("/v2/import", api.ImportFileToKnowledgeBaseV2)
	apiKnowledgeBaseGroup.POST("/v2/import/progress", api.ImportFileToKnowledgeBaseV2Progress)
	apiKnowledgeBaseGroup.GET("/:knowledgeBaseGuid/files", api.GetKnowledgeBaseFiles)
	apiKnowledgeBaseGroup.GET("/:knowledgeBaseGuid", api.GetKnowledgeBase)
	apiKnowledgeBaseGroup.DELETE("/:knowledgeBaseGuid/delete/:fileGuid", api.DeleteFileFromKnowledgeBase)
	apiKnowledgeBaseGroup.GET("/ai-assets", api.GetAiAssets)

	healthGroup := r.Group("/health")
	healthGroup.GET("/:path", api.GetWebInspect)
}

package http

import (
	"github.com/gotomicro/ego/server/egin"

	"sdk-demo-go/pkg/server/http/callback"
	"sdk-demo-go/pkg/server/http/middlewares"
)

// registerCallbackAPIs exposes the callback routes called by the Shimo SDK
func registerCallbackAPIs(r *egin.Component) {
	sdkcbGroup := r.Group("/callback", middlewares.CallbackAuthMiddleware)

	// 文件相关接口
	sdkcbGroup.GET("/files/:fileGuid", callback.GetFileInfo)
	sdkcbGroup.GET("/files", callback.GetFilesByUser)
	sdkcbGroup.POST("/files", callback.CreateFiles)
	sdkcbGroup.GET("/files/:fileGuid/collaborators", callback.GetFileCollaborators)
	sdkcbGroup.GET("/admin/files/:fileGuid", callback.AdminGetFileInfo)
	sdkcbGroup.GET("/admin/files/:fileGuid/by-user-id", callback.AdminGetFileInfoByUserId)
	sdkcbGroup.GET("/files/:fileGuid/download", callback.DownloadFile)
	sdkcbGroup.POST("/files/:fileGuid/url", callback.GetFileUrl)
	sdkcbGroup.POST("/files/import", callback.ImportFile)

	// 用户相关接口
	sdkcbGroup.GET("/users/current/info", callback.GetCurrentUser)
	sdkcbGroup.GET("/users/current/team", callback.GetCurrentTeam)
	sdkcbGroup.GET("/users/:userId", callback.GetSpecificUser)
	sdkcbGroup.GET("/users/:userId/watermark", callback.GetUserWatermark)
	sdkcbGroup.GET("/users/:userId/department-paths", callback.GetDepartmentPath)
	sdkcbGroup.POST("/users/batch/get", callback.BatchGetUserInfo)
	sdkcbGroup.POST("/admin/users/batch/get", callback.BatchGetUserInfo)

	// 团队及团队下部门相关接口
	sdkcbGroup.GET("/teams/:teamGuid/members", callback.GetTeamMembers)
	sdkcbGroup.GET("/departments/:departmentId", callback.GetDepartmentInfo)
	sdkcbGroup.GET("/departments/:departmentId/children", callback.GetSubDepts)
	sdkcbGroup.GET("/departments/:departmentId/members", callback.GetDeptMembers)

	// 搜索相关接口
	sdkcbGroup.GET("/search/users/recent", callback.SearchRelatedUsers)
	sdkcbGroup.GET("/search/files/recent", callback.SearchRelatedFiles)
	sdkcbGroup.POST("/search", callback.SearchByKeyword)

	// 其他杂项接口
	sdkcbGroup.GET("/link/queryUrl")
	sdkcbGroup.POST("/events", callback.PushEvent)
	sdkcbGroup.GET("/knowledgeBases", callback.GetKnowledgeBases)
}

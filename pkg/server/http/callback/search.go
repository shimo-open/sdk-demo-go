package callback

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/models/db"
)

func SearchRelatedUsers(c *gin.Context) {
	users, err := db.FindAllUsers(invoker.DB)
	if err != nil {
		handleDBError(c, err)
		return
	}

	userId := getUserIdFromToken(c)
	// This fetches all users and filters out the current user; redesign the logic for production use
	res := make([]UserInfo, 0)
	for i := range users {
		if users[i].ID == userId {
			continue
		}
		res = append(res, UserInfo{
			Id:        strconv.FormatInt(users[i].ID, 10),
			Name:      users[i].Name,
			Avatar:    users[i].Avatar,
			Email:     users[i].Email,
			CanBother: users[i].CanBother,
		})
	}

	c.JSON(200, res)
}

func SearchRelatedFiles(c *gin.Context) {
	userId := getUserIdFromToken(c)
	files, err := db.FindFileByUserId(invoker.DB, userId, 0, "")
	if err != nil {
		handleDBError(c, err)
		return
	}

	fileInfos := make([]FileInfo, 0)
	for i := range files {
		fileInfos = append(fileInfos, *loadFileInfo(&files[i]))
	}

	c.JSON(200, fileInfos)
}

type SearchByKeywordReq struct {
	FileId   string `json:"fileId"`
	Keyword  string `json:"keyword"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Type     string `json:"type"`
}

func SearchByKeyword(c *gin.Context) {
	userId := getUserIdFromToken(c)
	var req SearchByKeywordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "invalid request"})
		return
	}
	if req.PageSize < 1 {
		req.PageSize = 6
	}

	file, err := db.FindFileByGuid(invoker.DB, req.FileId)
	if err != nil {
		handleDBError(c, err)
		return
	}

	// Search collaborators
	var resCollaborators []UserInfo
	if strings.Contains(req.Type, "collaborator") {
		fps, err := db.FindFilePermissionsByFileId(invoker.DB, file.ID)
		if err != nil {
			handleDBError(c, err)
			return
		}
		userIds := make([]int64, 0)
		// Find users with read permission
		for i := range fps {
			if fps[i].Permissions["readable"] {
				userIds = append(userIds, fps[i].UserId)
			}
		}
		users, err := db.FindUsersByIds(invoker.DB, userIds)
		if err != nil {
			handleDBError(c, err)
			return
		}
		// Filter by keyword
		for _, user := range users {
			if strings.Contains(user.Name, req.Keyword) {
				resCollaborators = append(resCollaborators, UserInfo{
					Id:        strconv.FormatInt(user.ID, 10),
					Name:      user.Name,
					Avatar:    user.Avatar,
					Email:     user.Email,
					CanBother: user.CanBother,
				})
			} else if strings.Contains(req.Keyword, user.Email) {
				resCollaborators = append(resCollaborators, UserInfo{
					Id:        strconv.FormatInt(user.ID, 10),
					Name:      user.Name,
					Avatar:    user.Avatar,
					Email:     user.Email,
					CanBother: user.CanBother,
				})
			}
		}
	}

	// Search files
	var resFiles []FileInfo
	if strings.Contains(req.Type, "file") {
		files, err := db.FindFileByUserId(invoker.DB, userId, 100, "")
		if err != nil {
			handleDBError(c, err)
			return
		}
		for _, file := range files {
			if strings.Contains(file.Name, req.Keyword) {
				resFiles = append(resFiles, *loadFileInfo(&file))
			}
		}
	}

	// Search departments
	var resDepts []DeptInfo
	if strings.Contains(req.Type, "department") {
		teams, err := db.FindTeamsByUserId(invoker.DB, userId)
		if err != nil {
			handleDBError(c, err)
			return
		}
		if len(teams) > 0 {
			teamId := teams[0].ID
			depts, err := db.FindAllDepartmentsByTeamID(invoker.DB, teamId)
			if err != nil {
				handleDBError(c, err)
				return
			}

			var hitDeptIds []int64
			for _, dept := range depts {
				if strings.Contains(dept.Name, req.Keyword) {
					hitDeptIds = append(hitDeptIds, dept.ID)
				}
			}

			counts, err := db.CountDeptMembersByIds(invoker.DB, hitDeptIds)
			if err != nil {
				handleDBError(c, err)
				return
			}

			resDepts = make([]DeptInfo, len(depts))
			for i, dept := range depts {
				ancestors, err := db.FindDepartmentAllAncestorsById(invoker.DB, dept.ID, false)
				if err != nil {
					handleDBError(c, err)
					return
				}
				ancestorsInfo := make([]AncestorInfo, len(ancestors))
				for i := range ancestors {
					ancestorsInfo[i] = AncestorInfo{
						Id:   strconv.FormatInt(ancestors[i].ID, 10),
						Name: ancestors[i].Name,
					}
				}

				resDepts[i] = DeptInfo{
					Id:                strconv.FormatInt(dept.ID, 10),
					Name:              dept.Name,
					AllMemberCount:    counts[dept.ID],
					ParentDepartments: ancestorsInfo,
					CanBother:         dept.CanBother,
				}
			}
		}
	}

	// Search team members
	var resTeamMembers []UserInfo
	if strings.Contains(req.Type, "team_member") {
		teams, err := db.FindTeamsByUserId(invoker.DB, userId)
		if err != nil {
			handleDBError(c, err)
			return
		}
		if len(teams) > 0 {
			teamId := teams[0].ID
			userIds, err := db.FindTeamAllMembersByTeamId(invoker.DB, teamId)
			if err != nil {
				handleDBError(c, err)
				return
			}

			// Remove the current user
			flag := false
			for i := range userIds {
				if userIds[i] == userId {
					flag = true
					userIds[i] = userIds[0]
				}
			}
			if flag {
				userIds = userIds[1:]
			}

			members, err := db.FindUsersByIds(invoker.DB, userIds)
			if err != nil {
				handleDBError(c, err)
				return
			}

			for _, member := range members {
				if strings.Contains(member.Name, req.Keyword) {
					resTeamMembers = append(resTeamMembers, UserInfo{
						Id:        strconv.FormatInt(member.ID, 10),
						Name:      member.Name,
						Avatar:    member.Avatar,
						Email:     member.Email,
						CanBother: member.CanBother,
					})
				} else if strings.Contains(req.Keyword, member.Email) {
					resTeamMembers = append(resTeamMembers, UserInfo{
						Id:        strconv.FormatInt(member.ID, 10),
						Name:      member.Name,
						Avatar:    member.Avatar,
						Email:     member.Email,
						CanBother: member.CanBother,
					})
				}
			}
		}
	}

	// Search recent contacts
	var resRecentUsers []UserInfo
	if strings.Contains(req.Type, "recent") {
		users, err := db.FindAllUsers(invoker.DB)
		if err != nil {
			handleDBError(c, err)
			return
		}

		for _, user := range users {
			if strings.Contains(user.Name, req.Keyword) {
				resRecentUsers = append(resRecentUsers, UserInfo{
					Id:        strconv.FormatInt(user.ID, 10),
					Name:      user.Name,
					Avatar:    user.Avatar,
					Email:     user.Email,
					CanBother: user.CanBother,
				})
			} else if strings.Contains(req.Keyword, user.Email) {
				resRecentUsers = append(resRecentUsers, UserInfo{
					Id:        strconv.FormatInt(user.ID, 10),
					Name:      user.Name,
					Avatar:    user.Avatar,
					Email:     user.Email,
					CanBother: user.CanBother,
				})
			}
		}
	}

	if len(resTeamMembers) > req.PageSize {
		resTeamMembers = resTeamMembers[:req.PageSize]
	}
	if len(resCollaborators) > req.PageSize {
		resCollaborators = resCollaborators[:req.PageSize]
	}
	if len(resDepts) > req.PageSize {
		resDepts = resDepts[:req.PageSize]
	}
	if len(resFiles) > req.PageSize {
		resFiles = resFiles[:req.PageSize]
	}
	if len(resRecentUsers) > req.PageSize {
		resRecentUsers = resRecentUsers[:req.PageSize]
	}

	c.JSON(200, gin.H{
		"files": gin.H{
			"count":     len(resFiles),
			"page":      0,
			"pageSize":  req.PageSize,
			"pageCount": 1,
			"results":   resFiles,
		},
		"recentUsers": gin.H{
			"count":     len(resRecentUsers),
			"page":      0,
			"pageSize":  req.PageSize,
			"pageCount": 1,
			"results":   resRecentUsers,
		},
		"collaborators": gin.H{
			"count":     len(resCollaborators),
			"page":      0,
			"pageSize":  req.PageSize,
			"pageCount": 1,
			"results":   resCollaborators,
		},
		"department": gin.H{
			"count":     len(resDepts),
			"page":      0,
			"pageSize":  req.PageSize,
			"pageCount": 1,
			"results":   resDepts,
		},
		"teamMembers": gin.H{
			"count":     len(resTeamMembers),
			"page":      0,
			"pageSize":  req.PageSize,
			"pageCount": 1,
			"results":   resTeamMembers,
		},
	})
}

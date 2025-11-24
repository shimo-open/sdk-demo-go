package api

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/models/db"
)

// GetTeams retrieves all teams in the system
func GetTeams(c *gin.Context) {
	teams, err := db.FindAllTeams(invoker.DB)
	if err != nil {
		handleDBError(c, err)
		return
	}
	c.JSON(200, teams)
}

// GetTeamMembers retrieves all members of a specific team
func GetTeamMembers(c *gin.Context) {
	teamId := getInt64FromParam(c, "teamId")

	userIds, err := db.FindTeamAllMembersByTeamId(invoker.DB, teamId)
	if err != nil {
		handleDBError(c, err)
		return
	}

	users, err := db.FindUsersByIds(invoker.DB, userIds)
	if err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(200, users)
}

// JoinTeam adds a user to a team
func JoinTeam(c *gin.Context) {
	teamId := getInt64FromParam(c, "teamId")
	body := struct {
		UserId int64 `json:"userId"`
	}{}
	err := c.Bind(&body)
	if err != nil {
		return
	}
	userId := body.UserId
	if userId == 0 {
		userId = getUserIdFromToken(c)
	}
	err = db.JoinTeam(invoker.DB, teamId, userId)

	if err != nil {
		handleDBError(c, err)
		return
	}
	c.JSON(204, nil)
}

// TransferCreator transfers team creator role to another user
func TransferCreator(c *gin.Context) {
	teamId := getInt64FromParam(c, "teamId")

	body := struct {
		NewCreatorId int64 `json:"newCreatorId"`
	}{}
	err := c.BindJSON(&body)
	if err != nil {
		return
	}

	userId := getUserIdFromToken(c)
	oldCreatorId, err := db.FindTeamCreator(invoker.DB, teamId)
	if err != nil {
		handleDBError(c, err)
		return
	}
	if oldCreatorId != userId {
		c.JSON(403, gin.H{"message": "only creator can transfer creator"})
		return
	}

	err = db.TransferTeam(invoker.DB, teamId, body.NewCreatorId)
	if err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(204, nil)
}

// CreateTeam creates a new team with the current user as creator
func CreateTeam(c *gin.Context) {
	body := struct {
		Name string `json:"name"`
	}{}
	_ = c.BindJSON(&body)

	team := db.Team{
		Name: body.Name,
	}
	userId := getUserIdFromToken(c)
	err := db.CreateTeam(invoker.DB, &team, userId)
	if err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(204, nil)
}

// CreateDepartment creates a new department within a team
func CreateDepartment(c *gin.Context) {
	teamId := getInt64FromParam(c, "teamId")

	body := struct {
		Name     string `json:"name"`
		ParentId string `json:"parentId"`
	}{}
	err := c.BindJSON(&body)
	if err != nil {
		return
	}

	var parentId int64
	if body.ParentId != "root" {
		parentId, _ = strconv.ParseInt(body.ParentId, 10, 64)
	}

	err = db.CreateDepartment(invoker.DB, body.Name, parentId, teamId)
	if err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(204, nil)
}

// GetDepartmentTopTree retrieves the department tree structure for a team
func GetDepartmentTopTree(c *gin.Context) {
	teamId := getInt64FromParam(c, "teamId")

	root, err := db.FindDeptTree(invoker.DB, teamId)
	if err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(200, []db.DeptTreeNode{*root})
}

func GetDeptMembers(c *gin.Context) {
	dept := c.Param("deptId")
	if dept == "root" {
		c.JSON(200, []db.User{})
		return
	}

	deptId := getInt64FromParam(c, "deptId")

	members, err := db.FindDepartmentMembers(invoker.DB, deptId)
	if err != nil {
		handleDBError(c, err)
		return
	}

	userIds := make([]int64, len(members))
	for i := range members {
		userIds[i] = members[i].UserID
	}
	users, err := db.FindUsersByIds(invoker.DB, userIds)
	c.JSON(200, users)
}

func GetSubDeptsAndMembers(c *gin.Context) {
	teamId := getInt64FromParam(c, "teamId")

	dept := c.Param("deptId")
	var subDepts []db.Department
	var users []db.User
	var err error
	if dept == "root" {
		subDepts, err = db.FindRootDepartment(invoker.DB, teamId)
	} else {
		deptId, _ := strconv.ParseInt(dept, 10, 64)
		subDepts, err = db.FindSubDepartmentsByParentId(invoker.DB, deptId)
		if err != nil {
			handleDBError(c, err)
			return
		}

		members, err := db.FindDepartmentMembers(invoker.DB, deptId)
		if err != nil {
			handleDBError(c, err)
			return
		}

		userIds := make([]int64, len(members))
		for i := range members {
			userIds[i] = members[i].UserID
		}

		users, err = db.FindUsersByIds(invoker.DB, userIds)
	}

	resNodes := make([]db.DeptTreeNode, len(subDepts)+len(users))
	for i, dept := range subDepts {
		resNodes[i] = db.DeptTreeNode{
			Node:     dept,
			Type:     "department",
			Children: nil,
		}
	}
	for i, user := range users {
		resNodes[i+len(subDepts)] = db.DeptTreeNode{
			Node:     user,
			Type:     "user",
			Children: nil,
		}
	}

	c.JSON(200, resNodes)
}

func PatchMemberToDept(c *gin.Context) {
	deptId := getInt64FromParam(c, "deptId")
	teamId := getInt64FromParam(c, "teamId")

	body := struct {
		UserId int64 `json:"userId"`
	}{}
	_ = c.BindJSON(&body)

	_, err := db.FindTeamById(invoker.DB, teamId)
	if err != nil {
		handleDBError(c, err)
		return
	}

	err = db.UpdateUserDepartment(invoker.DB, deptId, body.UserId)
	if err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(204, nil)
}

func AddMemberToDept(c *gin.Context) {
	deptId := getInt64FromParam(c, "deptId")

	body := struct {
		UserId string `json:"userId"`
	}{}
	_ = c.BindJSON(&body)
	userId, _ := strconv.ParseInt(body.UserId, 10, 64)

	err := db.JoinDepartment(invoker.DB, deptId, userId)
	if err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(204, nil)
}

func DeleteMemberFromDept(c *gin.Context) {
	userId := getInt64FromParam(c, "userId")

	err := db.LeaveDepartment(invoker.DB, userId)
	if err != nil {
		handleDBError(c, err)
		return
	}
	c.JSON(204, nil)
}

func DeleteDept(c *gin.Context) {
	deptId := getInt64FromParam(c, "deptId")

	// TODO this ideally should be wrapped in a transaction
	err := db.RemoveDepartmentWithMembers(invoker.DB, deptId)
	if err != nil {
		handleDBError(c, err)
		return
	}

	err = db.RemoveSubDepartmentsWithMembers(invoker.DB, deptId)
	if err != nil {
		handleDBError(c, err)
		return
	}
	c.JSON(204, nil)
}

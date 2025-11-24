package callback

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/models/db"
)

// DeptInfo represents department information returned to Shimo
type DeptInfo struct {
	// Id is the department ID
	Id string `json:"id"`
	// Name is the department name
	Name string `json:"name"`
	// AllMemberCount is the total number of members in the department
	AllMemberCount int `json:"allMemberCount"`
	// ParentDepartments is the list of parent departments in the hierarchy
	ParentDepartments []AncestorInfo `json:"parentDepartments"`
	// CanBother indicates whether members can be disturbed
	CanBother bool `json:"canBother"`
}

// AncestorInfo represents parent department information in the hierarchy
type AncestorInfo struct {
	// Id is the parent department ID
	Id string `json:"id"`
	// Name is the parent department name
	Name string `json:"name"`
}

// TODO currently counts only the specified department and ignores sub-departments
func GetDepartmentInfo(c *gin.Context) {
	deptId := c.Param("departmentId")
	// If the ID starts with TEAM_, fetch every member under that team
	if strings.HasPrefix(deptId, "TEAM_") {
		tid, err := strconv.ParseInt(strings.TrimPrefix(deptId, "TEAM_"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "deptId format error",
			})
			return
		}
		team, err := db.FindTeamById(invoker.DB, tid)
		if err != nil {
			handleDBError(c, err)
			return
		}
		members, err := db.FindTeamAllMembersByTeamId(invoker.DB, tid)
		if err != nil {
			handleDBError(c, err)
			return
		}
		c.JSON(200, DeptInfo{
			Id:             strconv.FormatInt(tid, 10),
			Name:           team.Name,
			AllMemberCount: len(members),
			CanBother:      true,
		})
	} else {
		did, err := strconv.ParseInt(deptId, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "deptId format error",
			})
			return
		}
		dept, err := db.FindDepartmentById(invoker.DB, did)
		if err != nil {
			handleDBError(c, err)
			return
		}

		count, err := db.CountDeptMembersById(invoker.DB, did)
		if err != nil {
			handleDBError(c, err)
			return
		}

		c.JSON(200, DeptInfo{
			Id:             strconv.FormatInt(dept.ID, 10),
			Name:           dept.Name,
			AllMemberCount: int(count),
			CanBother:      dept.CanBother,
		})
	}

}

func GetSubDepts(c *gin.Context) {
	deptId := c.Param("departmentId")
	var subDepts []db.Department
	// If the ID starts with TEAM_, fetch every department under that team
	if strings.HasPrefix(deptId, "TEAM_") {
		tid, err := strconv.ParseInt(strings.TrimPrefix(deptId, "TEAM_"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "deptId format error",
			})
			return
		}
		subDepts, err = db.FindAllDepartmentsByTeamID(invoker.DB, tid)
		if err != nil {
			handleDBError(c, err)
			return
		}
	} else {
		did, err := strconv.ParseInt(deptId, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "deptId format error",
			})
			return
		}
		subDepts, err = db.FindSubDepartmentsByParentId(invoker.DB, did)
		if err != nil {
			handleDBError(c, err)
			return
		}
	}

	res := make([]DeptInfo, len(subDepts))
	for i := range subDepts {
		count, err := db.CountDeptMembersById(invoker.DB, subDepts[i].ID)
		if err != nil {
			handleDBError(c, err)
			return
		}
		res[i] = DeptInfo{
			Id:             strconv.FormatInt(subDepts[i].ID, 10),
			Name:           subDepts[i].Name,
			AllMemberCount: int(count),
			CanBother:      subDepts[i].CanBother,
		}
	}

	c.JSON(200, res)

}

func GetDeptMembers(c *gin.Context) {
	deptId := c.Param("departmentId")
	query := PaginationQuery{}
	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.JSON(400, gin.H{"message": "query params error"})
		return
	}
	type Users struct {
		UserID int64
	}
	var userIds []int64
	var tid, did int64
	// If the ID starts with TEAM_, fetch every member under that team
	if strings.HasPrefix(deptId, "TEAM_") {
		tid, err = strconv.ParseInt(strings.TrimPrefix(deptId, "TEAM_"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "deptId format error",
			})
		}
		members, err := db.FindTeamWithPagination(invoker.DB, tid, query.Page, query.PageSize)
		if err != nil {
			handleDBError(c, err)
			return
		}
		userIds = make([]int64, len(members))
		for i := range members {
			userIds[i] = members[i].UserID
		}
	} else {
		did, err = strconv.ParseInt(deptId, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "deptId format error",
			})
			return
		}
		members, err := db.FindDepartmentMembersWithPagination(invoker.DB, did, query.Page, query.PageSize)
		if err != nil {
			handleDBError(c, err)
			return
		}
		userIds = make([]int64, len(members))
		for i := range members {
			userIds[i] = members[i].UserID
		}
	}

	users, err := db.FindUsersByIds(invoker.DB, userIds)
	if err != nil {
		handleDBError(c, err)
		return
	}

	res := make([]UserInfo, len(users))
	for i := range users {
		res[i] = UserInfo{
			Id:        strconv.FormatInt(users[i].ID, 10),
			Name:      users[i].Name,
			Avatar:    users[i].Avatar,
			Email:     users[i].Email,
			CanBother: users[i].CanBother,
		}
	}

	c.JSON(200, gin.H{
		"total":   len(res),
		"members": res,
	})
}

func getDeptIdFromParams(c *gin.Context) int64 {
	_deptId := c.Param("departmentId")
	deptId, err := strconv.ParseInt(_deptId, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "deptId format error",
		})
		c.Abort()
	}
	return deptId
}

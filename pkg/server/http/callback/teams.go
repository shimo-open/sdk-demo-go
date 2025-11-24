package callback

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/models/db"
)

type PaginationQuery struct {
	Pagination bool `form:"pagination"`
	Page       int  `form:"page" binding:"required"`
	PageSize   int  `form:"pageSize" binding:"required"`
}

func GetTeamMembers(c *gin.Context) {
	teamID := getInt64FromParam(c, "teamGuid")
	query := PaginationQuery{}
	err := c.BindQuery(&query)
	if err != nil {
		return
	}

	var members []db.User
	if query.Pagination {
		trs, err := db.FindTeamWithPagination(invoker.DB, teamID, query.Page, query.PageSize)
		if err != nil {
			handleDBError(c, err)
			return
		}
		userIds := make([]int64, len(trs))
		for _, tr := range trs {
			userIds = append(userIds, tr.UserID)
		}

		members, err = db.FindUsersByIds(invoker.DB, userIds)
		if err != nil {
			handleDBError(c, err)
			return
		}
	} else {
		userIds, err := db.FindTeamAllMembersByTeamId(invoker.DB, teamID)
		if err != nil {
			handleDBError(c, err)
			return
		}
		members, err = db.FindUsersByIds(invoker.DB, userIds)
	}

	res := make([]UserInfo, len(members))
	for i := range members {
		res[i] = UserInfo{
			Id:        strconv.FormatInt(members[i].ID, 10),
			Name:      members[i].Name,
			Avatar:    members[i].Avatar,
			Email:     members[i].Email,
			CanBother: members[i].CanBother,
		}
	}
	c.JSON(http.StatusOK, res)
}

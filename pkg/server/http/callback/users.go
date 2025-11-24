package callback

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/econf"

	"sdk-demo-go/pkg/consts"
	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/models/db"
	"sdk-demo-go/pkg/server/http/api"
)

type UserInfo struct {
	Id              string            `json:"id"`
	Name            string            `json:"name"`
	Avatar          string            `json:"avatar"`
	Email           string            `json:"email"`
	ExtraAttributes map[string]string `json:"extraAttributes"`
	CanBother       bool              `json:"canBother"`
}

func GetCurrentUser(c *gin.Context) {
	userId := getUserIdFromToken(c)

	if userId < 0 {
		_anonUser := api.LoadAnonymousUser(consts.ANONYMOUS)
		anonUser := db.User{
			Name:   _anonUser.User.Name,
			Avatar: _anonUser.User.Avatar,
			Email:  _anonUser.User.Email,
		}
		anonUser.ID = _anonUser.User.ID
		sendUserInfo(c, &anonUser)
	} else {
		user, err := db.FindUserById(invoker.DB, userId)
		if err != nil {
			handleDBError(c, err)
			return
		}
		sendUserInfo(c, user)
	}
}

func GetCurrentTeam(c *gin.Context) {
	userId := getUserIdFromToken(c)
	teams, err := db.FindTeamsByUserId(invoker.DB, userId)
	if err != nil {
		handleDBError(c, err)
		return
	}

	if len(teams) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "team not found",
		})
		return
	}

	count, err := db.CountTeamMembers(invoker.DB, teams[0].ID)
	if err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"id":          strconv.FormatInt(teams[0].ID, 10),
		"name":        teams[0].Name,
		"memberCount": count,
	})
}

func GetSpecificUser(c *gin.Context) {
	userId := getInt64FromParam(c, "userId")

	user, err := db.FindUserById(invoker.DB, userId)
	if err != nil {
		handleDBError(c, err)
		return
	}

	sendUserInfo(c, user)
}

func GetUserWatermark(c *gin.Context) {
	userId := getInt64FromParam(c, "userId")
	// enableDownloadWatermark := c.Query("enableDownloadWatermark")
	exportType := c.Query("exportType")

	user, err := db.FindUserById(invoker.DB, userId)
	if err != nil {
		handleDBError(c, err)
		return
	}

	var watermarks []string

	if econf.GetBool("watermark.enable") {
		fields := econf.GetStringSlice("watermark.fields")
		watermarks = make([]string, len(fields))
		for i, field := range fields {
			watermarks[i] = user.GetField(field)
		}
	}

	if exportType != "" {
		dWatermarks := append(watermarks, "download")
		c.JSON(200, map[string]interface{}{
			"watermarks":         watermarks,
			"downloadWatermarks": dWatermarks,
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"watermarks": watermarks,
	})
}

func GetDepartmentPath(c *gin.Context) {
	userId := getInt64FromParam(c, "userId")

	depts, err := db.FindDeptsByUserId(invoker.DB, userId)
	if err != nil {
		handleDBError(c, err)
		return
	}

	var res [][]map[string]string

	for _, dept := range depts {
		path, err := db.FindDepartmentAllAncestorsById(invoker.DB, dept.ID, false)
		if err != nil {
			handleDBError(c, err)
			return
		}

		var temp []map[string]string
		for i := range path {
			temp = append(temp, map[string]string{
				"id":   strconv.FormatInt(path[len(path)-i-1].ID, 10),
				"name": path[len(path)-i-1].Name,
			})
		}
		res = append(res, temp)
	}
	c.JSON(200, res)
}

func BatchGetUserInfo(c *gin.Context) {
	body := struct {
		IDs []string `json:"ids"`
	}{}
	err := c.BindJSON(&body)
	if err != nil {
		return
	}

	ids := make([]int64, len(body.IDs))
	for i, id := range body.IDs {
		ids[i], err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "invalid body params -> ids",
			})
			return
		}
	}

	users, err := db.FindUsersByIds(invoker.DB, ids)
	if err != nil {
		handleDBError(c, err)
		return
	}

	var res []UserInfo
	for _, user := range users {
		res = append(res, UserInfo{
			Id:        strconv.FormatInt(user.ID, 10),
			Name:      user.Name,
			Avatar:    user.Avatar,
			Email:     user.Email,
			CanBother: user.CanBother,
		})

	}
	c.JSON(200, res)
}

func sendUserInfo(c *gin.Context, user *db.User) {
	c.JSON(200, UserInfo{
		Id:     strconv.FormatInt(user.ID, 10),
		Name:   user.Name,
		Avatar: user.Avatar,
		Email:  user.Email,
		ExtraAttributes: map[string]string{
			"附加值": "value1",
		},
	})
}

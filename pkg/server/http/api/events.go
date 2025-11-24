package api

import (
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/elog"

	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/models/db"
)

type EventParam struct {
	FileId string `json:"fileId"`
	UserId string `json:"userId"`
}

func GetEvents(c *gin.Context) {
	_page, _ := c.GetQuery("page")
	_size, _ := c.GetQuery("size")
	fileId, _ := c.GetQuery("fileId")
	userId, _ := c.GetQuery("userId")

	var page int
	if _page == "" || _page == "0" {
		page = 1
	} else {
		page, _ = strconv.Atoi(_page)
	}

	var size int
	if _size == "" || _size == "0" {
		size = 50
	} else {
		size, _ = strconv.Atoi(_size)
	}

	event := &db.Event{
		FileId: fileId,
		UserId: userId,
	}

	var wg sync.WaitGroup
	var events []db.Event
	var count int64

	wg.Add(2)
	go func() {
		events, _ = db.FindAllEvents(invoker.DB, event, page, size, "")
		wg.Done()
	}()
	go func() {
		count, _ = db.CountEvents(invoker.DB, event)
		wg.Done()
	}()
	wg.Wait()

	es, fileIds, userIds, err := db.GetEventFileAndUserIds(events)
	if err != nil {
		handleDBError(c, err)
		return
	}

	var files []db.File
	var users []db.User
	wg.Add(2)
	go func() {
		files, _ = db.FindFilesByGuids(invoker.DB, fileIds)
		wg.Done()
	}()
	go func() {
		uids := make([]int64, len(userIds))
		for i := range userIds {
			uids[i], _ = strconv.ParseInt(userIds[i], 10, 64)
		}
		users, _ = db.FindUsersByIds(invoker.DB, uids)
		wg.Done()
	}()

	wg.Wait()

	userMap := map[string]db.User{}
	for _, u := range users {
		uid := strconv.FormatInt(u.ID, 10)
		userMap[uid] = u
	}
	fileMap := map[string]db.File{}
	for _, f := range files {
		fileMap[f.Guid] = f
	}

	eventList := db.BindUserAndFilesInfo(es, userMap, fileMap)

	c.JSON(200, gin.H{
		"list":  eventList,
		"count": count,
		"page":  page,
		"size":  size,
	})
}

func GetSystemMessages(c *gin.Context) {
	from, _ := c.GetQuery("from")
	to, _ := c.GetQuery("to")
	appId, _ := c.GetQuery("AppId")
	elog.Info("from: " + from + "to: " + to + "appId: " + appId)
	c.JSON(200, nil)
}

func ErrorCallback(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "success",
	})
}

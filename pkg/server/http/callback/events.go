package callback

import (
	"encoding/json"
	"io"
	"strconv"

	sdkapi "github.com/shimo-open/sdk-kit-go/api"

	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/models/db"

	"github.com/gin-gonic/gin"
)

type EventReq interface {
	GetFileId() string
	GetUserId() string
}

type Comment struct {
	Guid             string   `json:"guid"`
	Content          string   `json:"content"`
	UserIds          []string `json:"userIds"`
	SelectionGuid    string   `json:"selectionGuid"`
	SelectionContent string   `json:"selectionContent"`
}

type Discussion struct {
	Id       string `json:"id"`
	Unixus   int64  `json:"unixus"`
	Content  string `json:"content"`
	Position string `json:"position"`
}

type Mention struct {
	Guid             string   `json:"guid"`
	Content          string   `json:"content"`
	UserIds          []string `json:"userIds"`
	SelectionGuid    string   `json:"selectionGuid"`
	SelectionContent string   `json:"selectionContent"`
}

type EventBody struct {
	Kind   string `json:"kind"`
	Type   string `json:"type"`
	Action string `json:"action"`
	FileId string `json:"fileId"`
	UserId string `json:"userId"`
}

func (e EventBody) GetFileId() string {
	return e.FileId
}

func (e EventBody) GetUserId() string {
	return e.UserId
}

func PushEvent(c *gin.Context) {
	eventType := c.GetHeader(sdkapi.HeaderShimoSdkEvent)
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	body := &EventBody{}
	err = json.Unmarshal(rawBody, body)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	rawHeader := c.Request.Header
	_rawHeader := make(map[string]string)
	for k, v := range rawHeader {
		_rawHeader[k] = v[0]
	}
	header, _ := json.Marshal(&_rawHeader)

	// Handle auto-mention events separately since they do not include userId
	uid := body.UserId
	if body.Type == "auto_mention" {
		var raw map[string]any
		if err := json.Unmarshal(rawBody, &raw); err == nil {
			if m, ok := raw["autoMention"].(map[string]any); ok {
				switch v := m["editorProviderId"].(type) {
				case string:
					uid = v
				case float64:
					uid = strconv.FormatInt(int64(v), 10)
				}
			}
		}
	}

	err = db.SaveEvent(invoker.DB, &db.Event{
		Type:    eventType,
		FileId:  body.FileId,
		UserId:  uid,
		RawData: string(rawBody),
		Headers: string(header),
	})
	if err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(204, nil)
}

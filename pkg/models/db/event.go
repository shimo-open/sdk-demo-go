package db

import (
	"encoding/json"
	"strconv"

	"github.com/gotomicro/cetus/l"
	"gorm.io/gorm"

	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
)

// Event represents a system event
type Event struct {
	BaseModel
	// Type is the event type (e.g., Comment, Discussion, MentionAt, etc.)
	Type string `gorm:"comment:'Event type'" json:"type"`
	// FileId is the file ID associated with this event
	FileId string `gorm:"comment:'File ID'" json:"fileId"`
	// UserId is the user ID who triggered this event
	UserId string `gorm:"index:idx_event_user_id;comment:'Related user ID'" json:"userId"`
	// RawData is the raw event data in JSON format
	RawData string `gorm:"comment:'Message content'" json:"rawData"`
	// Headers is the event headers in JSON format
	Headers string `gorm:"comment:'Event headers'" json:"headers"`
}

// SaveEvent stores an event
func SaveEvent(db *gorm.DB, e *Event) error {
	return db.Create(&e).Error
}

// FindAllEvents queries events (defaults: page=1, limit=10, orderBy=created_at)
func FindAllEvents(db *gorm.DB, e *Event, page int, limit int, orderBy string) (events []Event, err error) {
	if page <= 0 {
		page = 1
	}
	if orderBy == "" {
		orderBy = "created_at desc"
	}
	query := db.Table("events as e").Select("e.*")
	if e.FileId != "" {
		query = query.Where("e.file_id = ?", e.FileId)
	}

	if e.UserId != "" {
		query = query.Where("e.user_id = ?", e.UserId)
	}
	err = query.Joins("left join users as u on u.id = e.user_id").Joins("left join files as f on f.id = e.file_id").Where("u.app_id = ?", econf.GetString("shimoSDK.appId")).Order(orderBy).Offset((page - 1) * limit).Limit(limit).Find(&events).Error
	return
}

// CountEvents returns the event count
// Filters by FileId or UserId depending on which field is provided
func CountEvents(db *gorm.DB, e *Event) (count int64, err error) {
	query := db.Table("events as e")
	if e.FileId != "" {
		query = query.Where("e.file_id = ?", e.FileId)
	}

	if e.UserId != "" {
		query = query.Where("e.user_id = ?", e.UserId)
	}
	err = query.Joins("left join users as u on u.id = e.user_id").Joins("left join files as f on f.id = e.file_id").Where("u.app_id = ?", econf.GetString("shimoSDK.appId")).Count(&count).Error
	return
}

// EventWithDetails extends Event with related file and user information
type EventWithDetails struct {
	Event
	// FileIds is the list of file IDs related to this event
	FileIds []string `json:"fileIds"`
	// UserIds is the list of user IDs related to this event
	UserIds []string `json:"userIds"`
	// Files is a map of file ID to File for related files
	Files map[string]File `json:"files"`
	// Users is a map of user ID to User for related users
	Users map[string]User `json:"users"`
}

// GetEventFileAndUserIds augments each event with its related file and user IDs
func GetEventFileAndUserIds(events []Event) (es []*EventWithDetails, fileIds []string, userIds []string, err error) {
	for _, e := range events {
		ew, fIds, uIds, err := GetFileAndUserIdsFromEvent(&e)
		if err != nil {
			return nil, nil, nil, err
		}
		fileIds = append(fileIds, fIds...)
		userIds = append(userIds, uIds...)
		es = append(es, ew)
	}
	return
}

// GetFileAndUserIdsFromEvent extracts the file and user IDs for a single event
func GetFileAndUserIdsFromEvent(e *Event) (event *EventWithDetails, fileIds []string, userIds []string, err error) {
	var data map[string]interface{}
	err = json.Unmarshal([]byte(e.RawData), &data)
	if err != nil {
		return
	}

	switch e.Type {
	case "Comment":
		fileIds, userIds = idsInCommit(data)
	case "Discussion":
		fileIds, userIds = idsInDiscussion(data)
	case "MentionAt":
		fileIds, userIds = idsInMentionAt(data)
	case "DateMention":
		fileIds, userIds = idsInDateMention(data)
	case "FileContent":
		fileIds, userIds = idsInFileContent(data)
	case "Collaborator":
		fileIds, userIds = idsInCollaborator(data)
	default:
		if userId, ok := data["userId"].(string); ok {
			userIds = append(userIds, userId)
		}
		if fileId, ok := data["fileId"].(string); ok {
			fileIds = append(fileIds, fileId)
		}
	}

	event = &EventWithDetails{
		Event:   *e,
		FileIds: fileIds,
		UserIds: userIds,
	}
	return
}

// BindUserAndFilesInfo binds user and file information to events
func BindUserAndFilesInfo(events []*EventWithDetails, userMap map[string]User, fileMap map[string]File) (resEvents []*EventWithDetails) {
	resEvents = make([]*EventWithDetails, 0)
	for _, e := range events {
		valid := true
		if len(e.FileIds) > 0 {
			e.Files = make(map[string]File)
			for _, fileId := range e.FileIds {
				if f, ok := fileMap[fileId]; ok && f.ID > 0 {
					e.Files[fileId] = f
				} else {
					elog.Warn("BindUserAndFilesInfo file not found", elog.FieldComponent("db"), l.S("fileId", fileId))
					valid = false
				}
			}
		}

		if len(e.UserIds) > 0 {
			e.Users = make(map[string]User)
			for _, userId := range e.UserIds {
				if u, ok := userMap[userId]; ok && u.ID > 0 {
					e.Users[userId] = u
				} else {
					elog.Warn("BindUserAndFilesInfo user not found", elog.FieldComponent("db"), l.S("userId", userId))
					valid = false
				}
			}
		}

		if valid {
			resEvents = append(resEvents, e)
		}
	}
	return
}

func idsInCommit(data map[string]interface{}) (fileIds []string, userIds []string) {
	if fileId, ok := data["fileId"].(string); ok {
		fileIds = append(fileIds, fileId)
	}
	if userId, ok := data["userId"].(string); ok {
		userIds = append(userIds, userId)
	}

	action, ok := data["action"].(string)
	if !ok {
		return
	}

	switch action {
	case "update":
	case "create":
		uids := getStringSlice(data, "comment", "userIds")
		userIds = append(userIds, uids...)
	}
	return
}

func idsInDiscussion(data map[string]interface{}) (fileIds []string, userIds []string) {
	if fileId, ok := data["fileId"].(string); ok {
		fileIds = append(fileIds, fileId)
	}
	if userId, ok := data["userId"].(string); ok {
		userIds = append(userIds, userId)
	}
	return
}

func idsInMentionAt(data map[string]interface{}) (fileIds []string, userIds []string) {
	if fileId, ok := data["fileId"].(string); ok {
		fileIds = append(fileIds, fileId)
	}
	if userId, ok := data["userId"].(string); ok {
		userIds = append(userIds, userId)
	}

	typ, ok := data["type"].(string)
	if !ok {
		return
	}

	switch typ {
	case "comment":
		uids := getStringSlice(data, "comment", "userIds")
		userIds = append(userIds, uids...)
	case "discussion":
		uids := getStringSlice(data, "discussion", "userIds")
		userIds = append(userIds, uids...)
	case "mention_at":
		uids := getStringSlice(data, "mentionAt", "userId")
		userIds = append(userIds, uids...)
	}

	return
}

func idsInDateMention(data map[string]interface{}) (fileIds []string, userIds []string) {
	if fileId, ok := data["fileId"].(string); ok {
		fileIds = append(fileIds, fileId)
	}
	if userId, ok := data["userId"].(string); ok {
		userIds = append(userIds, userId)
	}

	action, ok := data["action"].(string)
	if !ok {
		return
	}

	switch action {
	case "create":
		fids := getStringSlice(data, "createData", "fileId")
		fileIds = append(fileIds, fids...)

		uids := getStringSlice(data, "createData", "authorId")
		userIds = append(userIds, uids...)

		uids = getStringSlice(data, "createData", "remindUserIds")
		userIds = append(userIds, uids...)

	case "update":
		uids := getStringSlice(data, "updateData", "remindUserIds")
		userIds = append(userIds, uids...)
	}
	return
}

func idsInFileContent(data map[string]interface{}) (fileIds []string, userIds []string) {
	if fileId, ok := data["fileId"].(string); ok {
		fileIds = append(fileIds, fileId)
	}
	if userId, ok := data["userId"].(string); ok {
		userIds = append(userIds, userId)
	}
	if m, ok := data["autoMention"].(map[string]any); ok {
		switch v := m["editorProviderId"].(type) {
		case string:
			userIds = append(userIds, v)
		case float64:
			userIds = append(userIds, strconv.FormatInt(int64(v), 10))
		}
	}
	return
}

func idsInCollaborator(data map[string]interface{}) (fileIds []string, userIds []string) {
	if fileId, ok := data["fileId"].(string); ok {
		fileIds = append(fileIds, fileId)
	}
	if userId, ok := data["userId"].(string); ok {
		userIds = append(userIds, userId)
	}
	return
}

func getStringSlice(data map[string]interface{}, t1, t2 string) []string {
	var res []string
	if comment, ok := data[t1].(map[string]interface{}); ok {
		if t1 == "mentionAt" {
			if str, ok := comment[t2].(string); ok {
				res = append(res, str)
			}
		} else {
			if str, ok := comment[t2].([]interface{}); ok {
				strSlice := make([]string, len(str))
				for i, v := range str {
					if str, ok := v.(string); ok {
						strSlice[i] = str
					}
				}
				res = append(res, strSlice...)
			}
		}
	}
	return res
}

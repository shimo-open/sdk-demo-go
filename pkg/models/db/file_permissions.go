package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// FilePermissions represents the permissions a user has for a file
type FilePermissions struct {
	BaseModel
	// FileId is the file ID
	FileId int64 `gorm:"uniqueIndex:uniq_file_id_user_id;comment:'File ID'" json:"fileId"`
	// UserId is the user ID
	UserId int64 `gorm:"uniqueIndex:uniq_file_id_user_id;comment:'User ID'" json:"userId"`
	// Role is the user's role (owner/collaborator)
	Role string `gorm:"default:'collaborator';comment:'Role (owner/collaborator)'" json:"role"`
	// Permissions is a JSON object of permission flags (stored as TEXT in DB)
	Permissions Permissions `gorm:"comment:'Permissions'" json:"permissions"`
}

// Permissions is a map of permission names to boolean values
type Permissions map[string]bool

func (fp *FilePermissions) TableName() string {
	return "file_permissions"
}

func (p Permissions) Value() (driver.Value, error) {
	v, err := json.Marshal(p)
	return string(v), err
}

func (p *Permissions) Scan(value interface{}) error {
	if value == nil {
		// If the value is NULL, return an empty Permissions struct
		*p = Permissions{}
		return nil
	}

	// Check whether the underlying type is []byte
	switch v := value.(type) {
	case []byte:
		// Decode JSON when the value is []byte
		return json.Unmarshal(v, p)
	case string:
		// Convert strings to []byte before decoding
		return json.Unmarshal([]byte(v), p)
	default:
		// Return an error for unsupported types
		return fmt.Errorf("failed to scan Permissions, unexpected type: %T", v)
	}
}

// FindFilePermissionsByFileId fetches every permission row for a file
func FindFilePermissionsByFileId(db *gorm.DB, fileId int64) (fps []FilePermissions, err error) {
	err = db.Where("file_id = ?", fileId).Find(&fps).Error
	return
}

// BatchSavePermissions upserts file permissions in bulk
// If a [userId, fileId] conflict occurs, overwrite the existing record instead of adding a new one
func BatchSavePermissions(db *gorm.DB, fps []FilePermissions) (err error) {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "file_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"role", "permissions"}),
	}).Create(&fps).Error
}

func GetUsersPermissions(db *gorm.DB, fileId int64, userId int64) (fpMap map[int64]FilePermissions, err error) {
	fpMap = make(map[int64]FilePermissions)
	fps := make([]FilePermissions, 0)
	err = db.Where("file_id = ? and user_id = ?", fileId, userId).Find(&fps).Error
	if err != nil {
		return
	}

	for _, fp := range fps {
		fpMap[fp.UserId] = fp
	}

	return
}

func RemovePermissionsByFileId(db *gorm.DB, fileId int64, userId int64) (err error) {
	return db.Unscoped().
		Where("file_id = ? AND user_id = ?", fileId, userId).
		Delete(&FilePermissions{}).Error
}

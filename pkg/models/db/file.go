package db

import (
	"errors"

	"github.com/gotomicro/cetus/l"
	"github.com/shimo-open/sdk-kit-go"

	"sdk-demo-go/pkg/utils"

	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"gorm.io/gorm"
)

// File represents a file in the system
type File struct {
	BaseModel
	// Guid is the unique file identifier
	Guid string `gorm:"uniqueIndex:uniq_guid;comment:'File GUID (unique identifier)'" json:"id"`
	// Name is the file name
	Name string `gorm:"comment:'File name'" json:"name"`
	// Type is the file type (document, spreadsheet, etc.)
	Type string `gorm:"comment:'File type'" json:"type"`
	// FilePath is the storage path of the file
	FilePath string `gorm:"comment:'File path'" json:"filePath"`
	// CreatorId is the user ID who created this file
	CreatorId int64 `gorm:"index:files_id_creator_id_index;comment:'Creator ID'" json:"creatorId"`
	// IsShimoFile indicates whether this is a Shimo collaborative file (1) or uploaded file (0)
	IsShimoFile int `gorm:"comment:'Is Shimo file'" json:"isShimoFile"`
	// ShimoType is the Shimo file type (document, spreadsheet, presentation, etc.)
	ShimoType string `gorm:"comment:'Shimo file type'" json:"shimoType"`
	// Permissions contains the file permissions (populated on demand, not stored in DB)
	Permissions Permissions `gorm:"-" json:"permissions"`
	// Role is the user's role for this file (populated on demand, not stored in DB)
	Role string `gorm:"-" json:"role"`
}

func (f *File) TableName() string {
	return "files"
}

// FindFileByUserId fetches up to 100 files for a user, ordered by creation time (desc)
// Because a user ID is provided, permissions are loaded by default
func FindFileByUserId(db *gorm.DB, userId int64, limit int, orderBy string) (files []File, err error) {
	if orderBy == "" {
		orderBy = "created_at"
	}
	if limit == 0 {
		limit = 100
	}

	var fps []FilePermissions
	err = db.Model(&FilePermissions{}).Where("user_id = ?", userId).Order("id desc").Limit(limit).Find(&fps).Error
	if err != nil {
		return
	}

	fileIds := make([]int64, len(fps))
	permissions := map[int64]map[string]bool{} // Holds permissions for each file
	role := map[int64]string{}
	for i := range fps {
		fileIds[i] = fps[i].FileId
		permissions[fps[i].FileId] = fps[i].Permissions
		role[fps[i].FileId] = fps[i].Role
	}

	err = db.Where("id IN ?", fileIds).Order(orderBy + " DESC").Find(&files).Error
	if err != nil {
		return
	}

	for i := range files {
		files[i].Permissions = permissions[files[i].ID]
		files[i].Role = role[files[i].ID]
	}
	return
}

// FindAllFiles fetches all files
func FindAllFiles(db *gorm.DB) (files []File, err error) {
	err = db.Find(&files).Error
	return
}

// FindFileById retrieves a file by ID
func FindFileById(db *gorm.DB, fileId int64) (file *File, err error) {
	err = db.Where("id = ?", fileId).First(&file).Error
	return
}

// FindFilesByIds fetches files by a list of IDs
func FindFilesByIds(db *gorm.DB, fileIds []int64) (files []File, err error) {
	err = db.Where("id IN ?", fileIds).Find(&files).Error
	return
}

// FindFileByGuid retrieves a file by GUID
func FindFileByGuid(db *gorm.DB, fileGuid string) (file *File, err error) {
	err = db.Where("guid = ?", fileGuid).First(&file).Error
	return
}

// FindFilesByGuids fetches files by a list of GUIDs
func FindFilesByGuids(db *gorm.DB, fileGuids []string) (files []File, err error) {
	err = db.Where("guid IN ?", fileGuids).Find(&files).Error
	return
}

// RenameFile updates a file name
func RenameFile(db *gorm.DB, guid string, newName string) (err error) {
	res := db.Model(&File{}).Where("guid = ?", guid).Updates(map[string]interface{}{
		"name": newName,
	})
	if res.RowsAffected == 0 {
		err = errors.New("no such file, check file id")
	}
	err = res.Error
	return
}

// CreateFile inserts a file
func CreateFile(db *gorm.DB, file *File, userId int64, permissions ...map[string]bool) (err error, fileId int64) {
	err = db.Transaction(func(tx *gorm.DB) error {
		if file.Guid == "" {
			file.Guid = utils.GenFileGuid()
		}

		if file.IsShimoFile == 0 {
			file.FilePath = file.Guid
		}

		err = tx.Create(&file).Error
		if err != nil {
			return err
		}

		fileId = file.ID

		filePermissionsList := make(map[string]bool)

		if len(permissions) > 0 {
			filePermissionsList = permissions[0]
		} else {
			filePermissionsList = sdk.HandleBasicFilePermission(true)

			if econf.GetBool("permissions.setNewFilePermission") {
				filePermissionsList = sdk.HandleFilePermission(true)
			}
		}

		fp := FilePermissions{
			FileId:      file.ID,
			UserId:      userId,
			Role:        "owner",
			Permissions: filePermissionsList,
		}

		file.Permissions = fp.Permissions

		err = tx.Create(&fp).Error
		return err
	})
	return
}

// RemoveFileByGuid deletes a file (and its permissions) by GUID
func RemoveFileByGuid(db *gorm.DB, fileGuid string) (err error) {
	f := &File{}
	err = db.Where("guid = ?", fileGuid).First(f).Error
	if err != nil || f.ID == 0 {
		elog.Warn("file guid is invalid: ", l.E(err))
		return nil
	}
	err = db.Where("guid = ?", fileGuid).Delete(&File{}).Error
	if err != nil {
		return err
	}
	return db.Where("file_id = ?", f.ID).Delete(&FilePermissions{}).Error
}

// RemoveFileByGuids deletes multiple files by GUID
func RemoveFileByGuids(db *gorm.DB, fileGuids []string) (err error) {
	err = db.Transaction(func(tx *gorm.DB) error {
		files := make([]File, 0)
		err = db.Where("guid IN ?", fileGuids).Find(&files).Error
		ids := make([]int64, len(files))
		for _, v := range files {
			ids = append(ids, v.ID)
		}
		err = db.Where("guid IN ?", fileGuids).Delete(&File{}).Error
		if err != nil {
			return err
		}
		err = db.Where("file_id IN ?", ids).Delete(&FilePermissions{}).Error
		return err
	})
	return
}

// RemoveFileById deletes a file (and its permissions) by ID
func RemoveFileById(db *gorm.DB, fileId int64) (err error) {
	err = db.Where("id = ?", fileId).Delete(&File{}).Error
	if err != nil {
		return err
	}
	return db.Where("file_id = ?", fileId).Delete(&FilePermissions{}).Error
}

// FindFileByGuidAndUserId fetches a file scoped to a user
// Loads permissions because the user ID is known
func FindFileByGuidAndUserId(db *gorm.DB, userId int64, guid string) (file *File, err error) {
	file, err = FindFileByGuid(db, guid)
	if err != nil {
		return
	}
	// Special handling for anonymous form access
	if userId < 0 && file.IsShimoFile == 1 && file.ShimoType == "form" {
		elog.Debug("file is form, anonymous user", elog.FieldComponent("db"))
		return
	}

	fp := FilePermissions{}
	err = db.Where("user_id = ? AND file_id = ?", userId, file.ID).First(&fp).Error
	if err != nil {
		// Ignore missing permission records for forms
		if errors.Is(err, gorm.ErrRecordNotFound) {
			file.Permissions = Permissions{}
			return file, nil
		}
		return
	}

	file.Permissions = fp.Permissions
	return
}

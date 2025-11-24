package db

import (
	"gorm.io/plugin/soft_delete"
)

// BaseModel is the base model for all database tables
type BaseModel struct {
	// ID is the primary key with auto increment
	ID int64 `gorm:"primaryKey; auto_increment" json:"id"`
	// CreatedAt is the Unix timestamp when the record was created
	CreatedAt int64 `gorm:"comment:'Created timestamp';autoCreateTime" json:"createdAt"`
	// UpdatedAt is the Unix timestamp when the record was last updated
	UpdatedAt int64 `gorm:"comment:'Updated timestamp';autoUpdateTime'" json:"updatedAt"`
	// DeletedAt is the Unix timestamp when the record was soft deleted (0 means not deleted)
	DeletedAt soft_delete.DeletedAt `gorm:"default:0;index;comment:'Deleted timestamp'" json:"-"`
}

package db

import (
	"gorm.io/gorm"
)

// KnowledgeBase represents a knowledge base entry
type KnowledgeBase struct {
	BaseModel
	// Guid is the unique knowledge base identifier
	Guid string `json:"guid"`
	// FileGuid is the file GUID associated with this knowledge base entry
	FileGuid string `json:"fileGuid"`
	// CreateBy is the user ID who created this entry
	CreateBy int64 `json:"createBy"`
	// DeleteAt is the Unix timestamp when this entry was deleted (0 means not deleted)
	DeleteAt int64 `json:"deleteAt"`
	// Name is the knowledge base name
	Name string `json:"name"`
}

// TableName returns the database table name for KnowledgeBase
func (kb *KnowledgeBase) TableName() string {
	return "knowledge_bases"
}

// CreateKnowledgeBase inserts a knowledge base record
func CreateKnowledgeBase(db *gorm.DB, kb *KnowledgeBase) error {
	return db.Create(kb).Error
}

// FindKnowledgeBasesByGuid fetches all records with the given knowledge base GUID
func FindKnowledgeBasesByGuid(db *gorm.DB, guid string) (kbs []KnowledgeBase, err error) {
	err = db.Where("guid = ? AND delete_at = 0", guid).Find(&kbs).Error
	return
}

// FindAllKnowledgeBaseGuids returns every unique knowledge base GUID
func FindAllKnowledgeBaseGuids(db *gorm.DB) (guids []string, err error) {
	err = db.Model(&KnowledgeBase{}).
		Where("delete_at = 0").
		Distinct().
		Pluck("guid", &guids).Error
	return
}

// FindAllKnowledgeBases returns every knowledge base grouped by GUID
func FindAllKnowledgeBases(db *gorm.DB) (kbs []KnowledgeBase, err error) {
	err = db.Model(&KnowledgeBase{}).
		Where("delete_at = 0").
		Select("guid, ANY_VALUE(name) as name").
		Group("guid").
		Find(&kbs).Error
	return
}

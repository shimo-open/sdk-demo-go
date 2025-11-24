package db

import (
	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/ego/core/elog"
	"gorm.io/gorm"
)

// AppClient represents an application client with credentials
type AppClient struct {
	BaseModel
	// AppID is the application identifier
	AppID string `json:"appId"`
	// AppSecret is the application secret key for authentication
	AppSecret string `json:"appSecret"`
}

// TableName returns the database table name for AppClient
func (t *AppClient) TableName() string {
	return "app_clients"
}

// AppClientFindById finds an AppClient by appId
func AppClientFindById(db *gorm.DB, appId string) (ac *AppClient, err error) {
	if err = db.Where("app_id = ?", appId).First(&ac).Error; err != nil {
		elog.Error("AppClient FindById db error", l.E(err))
		return nil, err
	}
	return
}

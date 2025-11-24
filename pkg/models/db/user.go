package db

import (
	"strconv"

	"github.com/gotomicro/ego/core/econf"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	BaseModel
	// Name is the user's display name
	Name string `gorm:"comment:'User name'" json:"name"`
	// Email is the user's email address (unique)
	Email string `gorm:"index:idx_email;unique;comment:'Email address'" json:"email"`
	// Avatar is the URL to the user's avatar image
	Avatar string `gorm:"comment:'Avatar URL'" json:"avatar"`
	// Password is the hashed password
	Password string `gorm:"comment:'Password'" json:"password"`
	// AppID is the application ID this user belongs to
	AppID string `gorm:"comment:'appId'" json:"appId"`
	// CanBother indicates whether the user can be disturbed
	CanBother bool `gorm:"comment:'Can bother'" json:"canBother" default:"true"`
}

// AllUser extends User with team information
type AllUser struct {
	User
	// TeamId is the team ID the user belongs to
	TeamId string `gorm:"comment:'Team ID'" json:"teamId"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) GetField(field string) string {
	switch field {
	case "id":
		return strconv.FormatInt(u.ID, 10)
	case "name":
		return u.Name
	case "avatar":
		return u.Avatar
	case "email":
		return u.Email
	default:
		return ""
	}
}

// FindAllUsers retrieves every user
func FindAllUsers(db *gorm.DB) (users []AllUser, err error) {
	// err = db.Find(&users).Error
	err = db.Table("users u").
		Select("u.*, t.team_id").
		Where("u.app_id = ?", econf.GetString("shimoSDK.appId")).
		Joins("LEFT JOIN team_role t ON t.user_id = u.id").
		Scan(&users).Error
	return
}

// FindUsersByAppId fetches users by app ID
func FindUsersByAppId(db *gorm.DB, appId string) (users []User, err error) {
	err = db.Where("app_id = ?", appId).Find(&users).Error
	return
}

// FindUserById fetches a user by ID
func FindUserById(db *gorm.DB, id int64) (user *User, err error) {
	err = db.Where("app_id = ? and id = ?", econf.GetString("shimoSDK.appId"), id).First(&user).Error
	return
}

// FindUsersByIds fetches users in bulk by their IDs
func FindUsersByIds(db *gorm.DB, ids []int64) (users []User, err error) {
	err = db.Where("app_id = ? and id IN ?", econf.GetString("shimoSDK.appId"), ids).Find(&users).Error
	return
}

// FindUserByInstance loads the full user info that matches the given instance; returns the first result on conflicts
func FindUserByInstance(db *gorm.DB, _user *User) (user *User, err error) {
	err = db.Where(_user).First(&user).Error
	return
}

// CheckUserExist determines whether a user exists based on the given instance
func CheckUserExist(db *gorm.DB, user *User) (exist bool, err error) {
	var cnt int64
	err = db.Model(&User{}).Where(&user).Count(&cnt).Error
	exist = cnt != 0
	return
}

// CreateUser inserts a new user
func CreateUser(db *gorm.DB, user *User) (err error) {
	res := db.Create(&user)
	err = res.Error
	return
}

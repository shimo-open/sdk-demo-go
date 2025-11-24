package api

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/econf"
	sdkapi "github.com/shimo-open/sdk-kit-go/model/api"
	"golang.org/x/crypto/bcrypt"

	"sdk-demo-go/pkg/consts"
	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/models/db"
	"sdk-demo-go/pkg/server/http/middlewares"
	"sdk-demo-go/pkg/utils"
)

// Auth authenticates a user and returns their information with app details
func Auth(c *gin.Context) {
	anonUser := LoadAnonymousUser(consts.ANONYMOUS)
	token := middlewares.FindAccessToken(c)
	if token == "" {
		c.JSON(http.StatusOK, anonUser)
		return
	}

	err := middlewares.ValidateUserToken(c, token)
	if err != nil {
		c.JSON(http.StatusOK, anonUser)
		return
	}

	userId := getUserIdFromToken(c)
	user, err := db.FindUserById(invoker.DB, userId)
	if err != nil {
		c.JSON(http.StatusOK, anonUser)
		return
	}

	auth := utils.GetAuth(user.ID, true)
	params := sdkapi.GetAppDetailParams{
		Auth:  auth,
		AppId: econf.GetString("shimoSDK.appId"),
	}
	detailsRes, err := invoker.SdkMgr.GetAppDetail(params)
	if err != nil {
		c.JSON(detailsRes.Resp.StatusCode(), gin.H{"message": "failed to get app details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":                 user,
		"token":                token,
		"availableFileTypes":   detailsRes.AvailableFileTypes,
		"premiumFileTypes":     detailsRes.PremiumFileTypes,
		"hostPath":             econf.GetString("host.addr"),
		"allowBoardAndMindmap": econf.GetBool("shimoSDK.allowBoardAndMindmap"),
	})
}

// GetUserById retrieves a user by their ID or "me" for the current user
func GetUserById(c *gin.Context) {
	id := c.Param("userId")

	var userId int64
	if id == consts.ME {
		userId = getUserIdFromToken(c)
	} else {
		userId = getInt64FromParam(c, "userId")
	}

	user, err := db.FindUserById(invoker.DB, userId)
	if err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetTeamsByUserId retrieves all teams associated with a user
func GetTeamsByUserId(c *gin.Context) {
	userId := getInt64FromParam(c, "userId")

	teams, err := db.FindTeamsByUserId(invoker.DB, userId)
	if err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(http.StatusOK, teams)
}

// DeleteMeFromTeam removes the current user from a team
func DeleteMeFromTeam(c *gin.Context) {
	teamId := getInt64FromParam(c, "teamId")

	creatorId, err := db.FindTeamCreator(invoker.DB, teamId)
	if err != nil {
		handleDBError(c, err)
		return
	}

	if creatorId == getUserIdFromToken(c) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "creator should transfer team before leaving",
		})
		return
	}

	err = db.LeaveTeam(invoker.DB, teamId, getUserIdFromToken(c))
	if err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// SignIn authenticates a user with email and password
func SignIn(c *gin.Context) {
	requestBody := struct {
		Email    string
		Password string
		AppId    string `json:"app_id"`
	}{}
	err := c.BindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	appId := requestBody.AppId
	if !c.GetBool("multipleClientMode") {
		client, _ := c.Get("appClient")
		appId = client.(db.AppClient).AppID
	}

	if appId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "missing app id"})
		return
	}

	if requestBody.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "missing email"})
		return
	}

	if requestBody.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "missing password"})
		return
	}

	user, err := db.FindUserByInstance(invoker.DB, &db.User{Email: requestBody.Email, AppID: appId})
	if err != nil {
		handleDBError(c, err)
		return
	}

	if !checkPassword(user.Password, requestBody.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid password"})
		return
	}

	tokenStr := utils.SignUserJWT(user.ID)

	auth := utils.GetAuth(user.ID, true)
	params := sdkapi.GetAppDetailParams{
		Auth:  auth,
		AppId: econf.GetString("shimoSDK.appId"),
	}
	appDetails, _ := invoker.SdkMgr.GetAppDetail(params)

	c.JSON(http.StatusOK, gin.H{
		"user":                 user,
		"token":                tokenStr,
		"availableFileTypes":   appDetails.AvailableFileTypes,
		"premiumFileTypes":     appDetails.PremiumFileTypes,
		"hostPath":             econf.GetString("host.addr"),
		"allowBoardAndMindmap": econf.GetBool("shimoSDK.allowBoardAndMindmap"),
	})
}

// SignUp creates a new user account
func SignUp(c *gin.Context) {
	if econf.GetBool("registrationDisabled") {
		c.JSON(http.StatusBadRequest, gin.H{"message": "registration is disabled"})
		return
	}

	requestBody := struct {
		Email    string
		Password string
		AppId    string `json:"app_id"`
	}{}
	err := c.BindJSON(&requestBody)
	if err != nil {
		return
	}

	appId := requestBody.AppId
	if !c.GetBool("multipleClientMode") {
		client, _ := c.Get("appClient")
		appId = client.(db.AppClient).AppID
	}

	if appId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "missing app id"})
		return
	}

	if requestBody.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "missing email"})
		return
	}

	if requestBody.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "missing password"})
		return
	}

	exist, err := db.CheckUserExist(invoker.DB, &db.User{Email: requestBody.Email, AppID: appId})
	if err != nil {
		handleDBError(c, err)
		return
	}

	if exist {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("%s already taken", requestBody.Email),
		})
		return
	}

	user := &db.User{
		Email:    requestBody.Email,
		Password: utils.HashPassword(requestBody.Password),
		Avatar:   getDefaultAvatar(),
		AppID:    appId,
		Name:     strings.Split(requestBody.Email, "@")[0],
	}
	err = db.CreateUser(invoker.DB, user)
	if err != nil {
		handleDBError(c, err)
		return
	}

	tokenStr := utils.SignUserJWT(user.ID)
	auth := utils.GetAuth(user.ID, true)
	params := sdkapi.GetAppDetailParams{
		Auth:  auth,
		AppId: econf.GetString("shimoSDK.appId"),
	}
	appDetails, _ := invoker.SdkMgr.GetAppDetail(params)

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":        user.ID,
			"name":      user.Name,
			"email":     user.Email,
			"appId":     user.AppID,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
			"avatar":    user.Avatar,
		},
		"token":                tokenStr,
		"availableFileTypes":   appDetails.AvailableFileTypes,
		"premiumFileTypes":     appDetails.PremiumFileTypes,
		"hostPath":             econf.GetString("host.addr"),
		"allowBoardAndMindmap": econf.GetBool("shimoSDK.allowBoardAndMindmap"),
	})
}

// GetAllUsers retrieves all users in the system
func GetAllUsers(c *gin.Context) {
	users, err := db.FindAllUsers(invoker.DB)
	if err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(http.StatusOK, users)
}

func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func getDefaultAvatar() string {
	if rand.Intn(10) >= 5 {
		return fmt.Sprintf("%sstatic/img/default-avatar-moke.png", econf.GetString("publicPath.publicPath"))
	} else {
		return fmt.Sprintf("%sstatic/img/default-avatar-moke-2.png", econf.GetString("publicPath.publicPath"))
	}
}

type AnonymousUserInfo struct {
	User     AnonymousUser `json:"user"`
	Token    string        `json:"token"`
	HostPath string        `json:"hostPath"`
}

type AnonymousUser struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Avatar          string `json:"avatar"`
	IsAnonymousUser bool   `json:"isAnonymousUser"`
}

// LoadAnonymousUser creates a random anonymous user
// LoadAnonymousUser creates an anonymous user info structure
func LoadAnonymousUser(userId int64) AnonymousUserInfo {
	return AnonymousUserInfo{
		User: AnonymousUser{
			ID:              userId,
			Name:            "匿名用户",
			Avatar:          getDefaultAvatar(),
			Email:           "anonymous@shimo.im",
			IsAnonymousUser: true,
		},
		Token:    utils.SignUserJWT(userId),
		HostPath: econf.GetString("host.addr"),
	}
}

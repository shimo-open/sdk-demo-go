package db

import (
	"gorm.io/gorm"
)

// TestApi represents an API test record
type TestApi struct {
	BaseModel
	// TestId is the unique test UUID
	TestId string `gorm:"comment:'Test UUID'" json:"testId"`
	// TestType is the type of test being run
	TestType string `gorm:"comment:'Test type'" json:"testType"`
	// ApiName is the name of the API being tested
	ApiName string `gorm:"comment:'API name'" json:"apiName"`
	// Success indicates whether the test passed (0=false, 1=true)
	Success int `gorm:"comment:'Success (0-false;1-true)'" json:"success"`
	// HttpCode is the HTTP status code returned
	HttpCode int `gorm:"comment:'Status code'" json:"httpCode"`
	// HttpResp is the HTTP response body
	HttpResp string `gorm:"comment:'Response result'" json:"httpResp"`
	// ErrMsg is the error message if the test failed
	ErrMsg string `gorm:"comment:'Error message'" json:"errMsg"`
	// PathStr is the API request path
	PathStr string `gorm:"comment:'API request path'" json:"pathStr"`
	// BodyReq is the request body parameters
	BodyReq string `gorm:"comment:'Body parameters'" json:"bodyReq"`
	// Query is the query string parameters
	Query string `gorm:"comment:'Query parameters'" json:"query"`
	// FormData is the form data parameters
	FormData string `gorm:"comment:'Form data parameters'" json:"formData"`
	// FileExt is the file extension or export file type
	FileExt string `gorm:"comment:'File extension/Export file type'" json:"fileExt"`
	// TimeConsuming is the time taken to complete the test
	TimeConsuming string `gorm:"comment:'Time consuming'" json:"timeConsuming"`
	// StartTime is the Unix timestamp when the test started
	StartTime int64 `gorm:"comment:'Test start time'" json:"startTime"`
}

func (t *TestApi) TableName() string {
	return "test_api"
}

func CreateTestApi(db *gorm.DB, testApi *[]TestApi) (err error) {
	err = db.Create(testApi).Error
	if err != nil {
		return err
	}
	return
}

func GetTestApi(db *gorm.DB, testId string) (testApis []TestApi, err error) {
	err = db.Where("test_id = ?", testId).Find(&testApis).Error
	return
}

func GetLatestTest(db *gorm.DB) (testApis []TestApi, err error) {
	tr := TestApi{}
	err1 := db.Order("created_at desc").First(&tr).Error
	if err1 != nil {
		return
	} else {
		err = db.Where("test_id = ?", tr.TestId).Find(&testApis).Error
		return
	}
}

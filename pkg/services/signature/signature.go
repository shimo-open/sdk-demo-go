package signature

import (
	"time"

	"sdk-demo-go/pkg/utils"
)

// SignatureService provides signature generation functionality
type SignatureService struct{}

// Init initializes a new SignatureService instance
func Init() *SignatureService {
	s := &SignatureService{}
	s.registerClients()
	return s
}

// Sign generates a JWT signature with the given credentials and expiration policy
func (s *SignatureService) Sign(appID, appSecret string, strict bool) string {
	nowTime := time.Now()
	var addTime time.Duration
	if strict {
		addTime = time.Minute * 4
	} else {
		addTime = time.Hour * 24 * 365
	}
	exp := nowTime.Add(addTime).Unix()

	if strict {
		return utils.SignJWT(appID, appSecret, exp, true)
	} else {
		return utils.SignJWT(appID, appSecret, exp, false)
	}

}

// ClientInfo holds application client credentials
type ClientInfo struct {
	AppID     string `json:"appId"`
	AppSecret string `json:"appSecret"`
}

func (s *SignatureService) registerClients() {
}

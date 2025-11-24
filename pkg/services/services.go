package services

import (
	"sdk-demo-go/pkg/services/awos"
	"sdk-demo-go/pkg/services/signature"

	"github.com/gotomicro/ego/client/ehttp"
)

// Services holds all service instances
type Services struct {
	// SignatureService handles signature generation
	SignatureService *signature.SignatureService
	// AwosService handles object storage operations
	AwosService *awos.AwosService
	// InspectHttp is the HTTP client for inspection service
	InspectHttp *ehttp.Component
}

// NewServices creates and initializes a new Services instance
func NewServices() *Services {
	return &Services{
		SignatureService: signature.Init(),
		AwosService:      awos.Init(),
		InspectHttp:      ehttp.Load("frontInspect.http").Build(),
	}
}

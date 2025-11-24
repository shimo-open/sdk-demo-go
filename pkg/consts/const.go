package consts

const (
	// ME represents the current user in API requests
	ME = "me"
	// TOKEN is the header name for authentication token
	TOKEN = "X-Shimo-Token"
	// SIGNATURE is the header name for request signature
	SIGNATURE = "X-Shimo-Signature"
	// CREDENTIAL is the header name for credential type
	CREDENTIAL = "X-Shimo-Credential-Type"
	// SDKEVENT is the header name for SDK events
	SDKEVENT = "X-Shimo-Sdk-Event"
	// ANONYMOUS is the user ID for anonymous users
	ANONYMOUS = -1
	// ANONYMOUSTOKEN is the token used for anonymous users
	ANONYMOUSTOKEN = "pseudonymoustoken"
)

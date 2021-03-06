package utility

// APIResponse is the type to return to the callers
type APIResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Result     interface{} `json:"result"`
	Code       string      `json:"code"`
}

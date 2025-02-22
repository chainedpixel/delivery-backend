package responser

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   *APIError   `json:"error,omitempty"`
}

type APIErrorResponse struct {
	Error *APIError `json:"error"`
}

type APIError struct {
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
	Code    string   `json:"code"`
}

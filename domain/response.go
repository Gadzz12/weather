package domain

// WebResponse represents a standard unified JSON response
type WebResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// ErrorResponse represents an error payload
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

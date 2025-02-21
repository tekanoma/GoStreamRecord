package web_response

// Response is a generic response structure for our API endpoints.
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

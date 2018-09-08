package transport

// APIResponse represents a generic API response
type APIResponse struct {
	Data interface{}            `json:"data"`
	Meta map[string]interface{} `json:"meta"`
}

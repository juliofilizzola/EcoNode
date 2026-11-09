package dto

type ErrorResponse struct {
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

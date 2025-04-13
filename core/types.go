package core

import "fmt"

type (
	// ErrorResponse represents the structure of the error response
	ErrorResponse struct {
		Error ErrorDetails `json:"error"` // Embed the ErrorDetails struct
	}

	// ErrorDetails holds the error code and message
	ErrorDetails struct {
		Code    int    `json:"code"`    // Error code (e.g., 404)
		Message string `json:"message"` // Error message (e.g., "Resource not found")
	}

	NotFoundError struct {
		ID       string
	}
)

// Error implements the error interface for NotFoundError
func (e NotFoundError) Error() string {
	return fmt.Sprintf("VM with ID %s not found", e.ID)
}

// NewNotFoundError creates a new NotFoundError
func NewNotFoundError(id string) *NotFoundError {
	return &NotFoundError{
		ID:       id,
	}
}
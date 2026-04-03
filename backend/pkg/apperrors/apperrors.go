package apperrors

import "net/http"

// AppError represents a structured API error standardising the JSON output
type AppError struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Error implements the standard Go error interface
func (e *AppError) Error() string {
	return e.Message
}

func NewValidationError(details interface{}) *AppError {
	return &AppError{
		Status:  http.StatusBadRequest,
		Message: "Validation Failed",
		Details: details,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Status:  http.StatusNotFound,
		Message: message,
	}
}

func NewInternalError(details interface{}) *AppError {
	return &AppError{
		Status:  http.StatusInternalServerError,
		Message: "Internal Server Error",
		Details: details,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Status:  http.StatusUnauthorized,
		Message: message,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Status:  http.StatusForbidden,
		Message: message,
	}
}

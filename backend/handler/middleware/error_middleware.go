package middleware

import (
	"errors"
	"net/http"

	"backend/pkg/apperrors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrorHandler is a global middleware catching deferred panics and intercepting c.Errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Proceed through handlers
		c.Next()

		// If there's an error detected in the queue, process it cleanly
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// 1. Is it a custom application error?
			var appErr *apperrors.AppError
			if errors.As(err, &appErr) {
				c.JSON(appErr.Status, appErr)
				return
			}

			// 2. Is it a request body JSON validation error?
			var valErrs validator.ValidationErrors
			if errors.As(err, &valErrs) {
				formattedErrors := make(map[string]string)
				for _, e := range valErrs {
					// Creates nice field-by-field breakdowns
					formattedErrors[e.Field()] = "Failed validation: requirement '" + e.ActualTag() + "' not met."
				}

				c.JSON(http.StatusBadRequest, apperrors.AppError{
					Status:  http.StatusBadRequest,
					Message: "Input Validation Failed",
					Details: formattedErrors,
				})
				return
			}

			// 3. Fallback standard error
			c.JSON(http.StatusInternalServerError, apperrors.AppError{
				Status:  http.StatusInternalServerError,
				Message: "An unexpected internal server error occurred",
				Details: err.Error(),
			})
		}
	}
}

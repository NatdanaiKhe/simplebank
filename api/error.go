package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/NatdanaiKhe/simplebank/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// Error codes returned to API clients. Constants prevent typos and
// let client code switch on the code field.
const (
	ErrCodeValidation = "VALIDATION_ERROR"
	ErrCodeNotFound   = "NOT_FOUND"
	ErrCodeConflict   = "CONFLICT"
	ErrCodeBadRequest = "BAD_REQUEST"
	ErrCodeInternal   = "INTERNAL_ERROR"
)

// FieldError describes a single validation failure on a specific field.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorResponse is the standardized JSON error body for every endpoint.
type ErrorResponse struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Fields  []FieldError `json:"fields,omitempty"`
}

// errorResponse inspects the error, maps it to an appropriate HTTP status
// and user-safe message, and writes the JSON response to the client.
func errorResponse(c *gin.Context, err error) {
	status, resp := mapError(c, err)
	c.JSON(status, resp)
}

// mapError classifies the input error and returns the correct HTTP status
// and client-safe ErrorResponse. Unknown errors are logged internally but
// never leaked to the client.
func mapError(c *gin.Context, err error) (int, ErrorResponse) {
	// 1. Keep the Gin validation logic here (it's an API-level concern)
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		return http.StatusBadRequest, ErrorResponse{
			Code:    ErrCodeValidation,
			Message: "Validation failed",
			Fields:  formatValidationErrors(ve),
		}
	}

	// 2. Use a map or a simple switch for Domain Errors
	// This is the ONLY place the API looks at service errors
	switch {
	case errors.Is(err, service.ErrAccountNotFound):
		return http.StatusNotFound, ErrorResponse{Code: ErrCodeNotFound, Message: err.Error()}
	case errors.Is(err, service.ErrInsufficientFunds):
		return http.StatusBadRequest, ErrorResponse{Code: ErrCodeBadRequest, Message: err.Error()}
	case errors.Is(err, service.ErrAccountAlreadyExists):
		return http.StatusConflict, ErrorResponse{Code: ErrCodeConflict, Message: err.Error()}
	case errors.Is(err, service.ErrUnsupportedCurrency):
		return http.StatusUnprocessableEntity, ErrorResponse{Code: ErrCodeBadRequest, Message: err.Error()}
	case errors.Is(err, service.ErrForeignKeyViolation):
		return http.StatusBadRequest, ErrorResponse{Code: ErrCodeBadRequest, Message: err.Error()}
	case errors.Is(err, service.ErrDuplicateUsername):
		return http.StatusConflict, ErrorResponse{Code: ErrCodeConflict, Message: err.Error()}
	}

	// Transfer specific errors

	switch {
	case errors.Is(err, service.ErrTransferCurrencyMismatch):
		return http.StatusBadRequest, ErrorResponse{Code: ErrCodeBadRequest, Message: err.Error()}
	}

	// 3. Fallback for anything that wasn't translated by the service
	GetLogger(c).Error("internal error",
		zap.String("error_message", err.Error()),
	)
	return http.StatusInternalServerError, ErrorResponse{
		Code:    ErrCodeInternal,
		Message: "An unexpected error occurred",
	}
}

func formatValidationErrors(fe validator.ValidationErrors) []FieldError {
	fields := make([]FieldError, 0, len(fe))
	for _, e := range fe {
		fields = append(fields, FieldError{
			Field:   e.Field(),
			Message: formatValidationError(e),
		})
	}
	return fields
}

// formatValidationError translates a validator tag into a human-readable message.
// Extend this as you add more binding rules.
func formatValidationError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return fmt.Sprintf("Must be at least %s", fe.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", fe.Param())
	default:
		return fmt.Sprintf("Failed on '%s' validation", fe.Tag())
	}
}

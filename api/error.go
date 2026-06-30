package api

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
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
	status, resp := mapError(err)
	c.JSON(status, resp)
}

// mapError classifies the input error and returns the correct HTTP status
// and client-safe ErrorResponse. Unknown errors are logged internally but
// never leaked to the client.
func mapError(err error) (int, ErrorResponse) {
	// 1. Validation errors from Gin binding (ShouldBindJSON / ShouldBindUri / ShouldBindQuery)
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		fields := make([]FieldError, 0, len(ve))
		for _, fe := range ve {
			fields = append(fields, FieldError{
				Field:   fe.Field(),
				Message: formatValidationError(fe),
			})
		}
		return http.StatusBadRequest, ErrorResponse{
			Code:    ErrCodeValidation,
			Message: "Validation failed",
			Fields:  fields,
		}
	}

	// 2. Database "not found" — sql.ErrNoRows
	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound, ErrorResponse{
			Code:    ErrCodeNotFound,
			Message: "Resource not found",
		}
	}

	// 3. PostgreSQL driver errors — we map pg error codes to
	//    meaningful HTTP statuses.
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505": // unique_violation
			return http.StatusConflict, ErrorResponse{
				Code:    ErrCodeConflict,
				Message: "Resource already exists",
			}
		case "23503": // foreign_key_violation
			return http.StatusBadRequest, ErrorResponse{
				Code:    ErrCodeBadRequest,
				Message: "Referenced resource does not exist",
			}
		case "23502": // not_null_violation
			return http.StatusBadRequest, ErrorResponse{
				Code:    ErrCodeBadRequest,
				Message: "A required field is missing",
			}
		}
	}

	// 4. Everything else — log the real error on the server,
	//    return a generic message to the client.
	log.Printf("internal error: %v", err)
	return http.StatusInternalServerError, ErrorResponse{
		Code:    ErrCodeInternal,
		Message: "An unexpected error occurred",
	}
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

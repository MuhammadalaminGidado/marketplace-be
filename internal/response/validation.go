package response

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidationError(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		messages := make([]string, 0, len(ve))
		for _, fe := range ve {
			messages = append(messages, fieldErrorMessage(fe))
		}
		Error(c, http.StatusBadRequest, strings.Join(messages, "; "))
		return
	}

	// json.UnmarshalTypeError, json.SyntaxError, io.EOF (empty body), etc.
	// all collapse to one generic message — never expose Go type/field internals.
	Error(c, http.StatusBadRequest, "invalid request body")
}

func fieldErrorMessage(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())
	switch fe.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email"
	case "min":
		return field + " must be at least " + fe.Param() + " characters"
	case "max":
		return field + " must be at most " + fe.Param() + " characters"
	case "eqfield":
		return field + " must match " + strings.ToLower(fe.Param())
	default:
		return field + " is invalid"
	}
}

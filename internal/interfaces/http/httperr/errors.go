package httperr

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"reader-club/internal/application/apperr"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Respond(c *gin.Context, status int, code, message string) {
	c.JSON(status, APIError{Code: code, Message: message})
}

func Abort(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, APIError{Code: code, Message: message})
}

func RespondInternalError(c *gin.Context) {
	Respond(c, http.StatusInternalServerError, "internal_error", "an unexpected error occurred")
}

func RespondValidation(c *gin.Context, err error) {
	Respond(c, http.StatusBadRequest, "validation_error", err.Error())
}

// RespondError maps well-known application errors to HTTP status + code.
// Callers can bypass this and call Respond directly when a specific message is required.
func RespondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperr.ErrNotFound):
		Respond(c, http.StatusNotFound, "not_found", "resource not found")
	case errors.Is(err, apperr.ErrAlreadyExists):
		Respond(c, http.StatusConflict, "conflict", "resource already exists")
	case errors.Is(err, apperr.ErrForbidden):
		Respond(c, http.StatusForbidden, "forbidden", "forbidden")
	case errors.Is(err, apperr.ErrInvalidCredentials):
		Respond(c, http.StatusUnauthorized, "unauthorized", "invalid credentials")
	case errors.Is(err, apperr.ErrThemeAlreadyExists):
		Respond(c, http.StatusConflict, "conflict", "a theme already exists for this month")
	case errors.Is(err, apperr.ErrNoEligibleBooks):
		Respond(c, http.StatusUnprocessableEntity, "unprocessable_entity", "no eligible book suggestions to draw from")
	case errors.Is(err, apperr.ErrMemberLimitReached):
		Respond(c, http.StatusUnprocessableEntity, "unprocessable_entity", "club member limit has been reached")
	default:
		RespondInternalError(c)
	}
}

package utils

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/refine-software/afrad-api/internal/database"
)

type APIError struct {
	Message string `json:"msg"`
	Code    int    `json:"code"`
}

func (e *APIError) Error() string {
	return e.Message
}

func NewAPIError(code int, message string) *APIError {
	return &APIError{Code: code, Message: message}
}

var (
	ErrNotFound     = NewAPIError(http.StatusNotFound, "Requested resource not found.")
	ErrUnauthorized = NewAPIError(http.StatusUnauthorized, "Invalid credentials.")
	ErrBadRequest   = NewAPIError(http.StatusBadRequest, "Invalid request data.")
	ErrInternal     = NewAPIError(
		http.StatusInternalServerError,
		"Something went wrong, please try again later.",
	)
	ErrStolenToken = NewAPIError(http.StatusUnauthorized, "sus behavior")
	ErrForbidden   = NewAPIError(
		http.StatusForbidden,
		"you're not allowed to access this resource",
	)
	ErrHeaderMissing = func(headerName string) *APIError {
		return NewAPIError(http.StatusBadRequest, fmt.Sprintf("%s header is missing", headerName))
	}
	ErrInvalidCredentials = NewAPIError(http.StatusUnauthorized, "invalid credentials")

	// DB Errors
	ErrForeignKeyViolation = func(columnName string) *APIError {
		return NewAPIError(http.StatusBadRequest, fmt.Sprintf("No such %s", columnName))
	}

	ErrUniqueViolation = func(columnName string) *APIError {
		return NewAPIError(http.StatusConflict, fmt.Sprintf("%s already exists", columnName))
	}

	ErrStringTooLong = func(columnName string) *APIError {
		return NewAPIError(http.StatusBadRequest, fmt.Sprintf("%s is too long", columnName))
	}

	ErrInvalidFormat = func(columnName string) *APIError {
		return NewAPIError(http.StatusBadRequest, fmt.Sprintf("Invalid format for %s", columnName))
	}

	ErrNullViolation = func(columnName string) *APIError {
		return NewAPIError(http.StatusBadRequest, fmt.Sprintf("%s is required", columnName))
	}

	ErrCheckViolation = func(columnName string) *APIError {
		return NewAPIError(http.StatusBadRequest, fmt.Sprintf("Invalid value for %s", columnName))
	}

	// Token Errors
	ErrTokenExpired       = NewAPIError(http.StatusUnauthorized, jwt.ErrTokenExpired.Error())
	ErrTokenInvalidClaims = NewAPIError(http.StatusUnauthorized, jwt.ErrTokenInvalidClaims.Error())
	ErrParsingToken       = NewAPIError(http.StatusUnauthorized, "unable to parse token")
	ErrInvalidToken       = NewAPIError(http.StatusUnauthorized, "invalid token")

	// Role Errors
	ErrRoleNotAllowed = NewAPIError(http.StatusForbidden, "role not allowed")
)

func MapDBErrorToAPIError(err error, columnName string) *APIError {
	switch err {
	case database.ErrDuplicate:
		return ErrUniqueViolation(columnName)
	case database.ErrNotFound:
		return ErrNotFound
	case database.ErrInvalidInput:
		return ErrBadRequest
	case database.ErrForeignKey:
		return ErrForeignKeyViolation(columnName)
	case database.ErrStringTooLong:
		return ErrStringTooLong(columnName)
	case database.ErrInvalidFormat:
		return ErrInvalidFormat(columnName)
	case database.ErrNullViolation:
		return ErrNullViolation(columnName)
	case database.ErrCheckViolation:
		return ErrCheckViolation(columnName)

	case database.ErrUnknown:
		return ErrInternal
	}

	return ErrInternal
}

package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func Respond(
	c *gin.Context,
	status int,
	success bool,
	message string,
	data any,
	err any,
) {
	if err != nil {
		log.Printf(
			"Internal Error: %v | Path: %s | Method: %s\n",
			err, c.Request.URL.Path, c.Request.Method,
		)
	}
	c.JSON(status, APIResponse{
		Success: success,
		Message: message,
		Data:    data,
		Error:   err,
	})
}

func RespondAndAbort(
	c *gin.Context,
	status int,
	success bool,
	message string,
	data any,
	err error,
) {
	if err != nil {
		log.Printf(
			"Internal Error: %v | Path: %s | Method: %s\n",
			err, c.Request.URL.Path, c.Request.Method,
		)
	}
	c.AbortWithStatusJSON(status, APIResponse{
		Success: success,
		Message: message,
		Data:    data,
		Error:   err,
	})
}

func Success(c *gin.Context, message string, data any) {
	Respond(c, http.StatusOK, true, message, data, nil)
}

func Created(c *gin.Context, message string, data any) {
	Respond(c, http.StatusCreated, true, message, data, nil)
}

func Fail(c *gin.Context, apiErr error, loggedError error) {
	switch e := apiErr.(type) {
	case *APIError:
		Respond(c, e.Code, false, e.Message, nil, loggedError)
	case error:
		Respond(c, http.StatusInternalServerError, false, e.Error(), nil, loggedError)
	default:
		Respond(c, http.StatusInternalServerError, false, "Unknown error", nil, loggedError)
	}
}

func FailAndAbort(c *gin.Context, apiErr error, loggedError error) {
	switch e := apiErr.(type) {
	case *APIError:
		RespondAndAbort(c, e.Code, false, e.Message, nil, loggedError)
	case error:
		RespondAndAbort(c, http.StatusInternalServerError, false, e.Error(), nil, loggedError)
	default:
		RespondAndAbort(c, http.StatusInternalServerError, false, "Unknown error", nil, loggedError)
	}
}

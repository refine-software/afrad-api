package utils

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/database"
)

func Success(c *gin.Context, data any) {
	if data == nil {
		c.Status(http.StatusOK)
	}
	c.JSON(http.StatusOK, data)
}

func Created(c *gin.Context, data any) {
	if data == nil {
		c.Status(http.StatusCreated)
	}
	c.JSON(http.StatusCreated, data)
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Fail(c *gin.Context, apiErr *APIError, loggedError error) {
	var dbErr database.DBError
	if errors.As(loggedError, &dbErr) {
		log.Printf(
			"Database Error: %s | Repository: %s | Method: %s | Code: %s | Error: %s",
			dbErr.Message,
			dbErr.Repo,
			dbErr.Method,
			dbErr.Code,
			dbErr.Err.Error(),
		)
	} else {
		log.Printf(
			"Internal Error: %v | Path: %s | Method: %s\n",
			loggedError,
			c.Request.URL.Path,
			c.Request.Method,
		)
	}

	c.JSON(apiErr.Code, apiErr)
}

func FailAndAbort(c *gin.Context, apiErr *APIError, loggedError error) {
	log.Printf(
		"Internal Error: %v | Path: %s | Method: %s\n",
		loggedError,
		c.Request.URL.Path,
		c.Request.Method,
	)

	c.AbortWithStatusJSON(apiErr.Code, apiErr)
}

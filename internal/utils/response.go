package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Fail(c *gin.Context, apiErr *APIError, loggedError error) {
	log.Printf(
		"Internal Error: %v | Path: %s | Method: %s\n",
		loggedError,
		c.Request.URL.Path,
		c.Request.Method,
	)

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

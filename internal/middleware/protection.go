package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/auth"
	"github.com/refine-software/afrad-api/internal/utils"
)

func AdminOnly(accessTokenSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := auth.GetAccessClaimsFromAuthHeader(c, accessTokenSecret)

		if claims == nil {
			return
		}

		if claims.Role != "admin" {
			utils.FailAndAbort(
				c,
				utils.ErrForbidden,
				nil,
			)
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func AuthRequired(accessTokenSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := auth.GetAccessClaimsFromAuthHeader(c, accessTokenSecret)

		if claims == nil {
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

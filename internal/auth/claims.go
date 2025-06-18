package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/refine-software/afrad-api/internal/utils"
)

type AccessClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GetAccessClaims(c *gin.Context) *AccessClaims {
	claimsInterface, exists := c.Get("claims")
	if !exists {
		utils.Fail(c, utils.ErrUnauthorized, errors.New("no claims in token"))
		return nil
	}

	claims, ok := claimsInterface.(*AccessClaims)
	if !ok {
		utils.Fail(c, utils.ErrUnauthorized, errors.New("bad claims type"))
		return nil
	}
	return claims
}

func GetAccessClaimsFromAuthHeader(c *gin.Context, accessTokenSecret string) *AccessClaims {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.FailAndAbort(
			c,
			utils.NewAPIError(http.StatusUnauthorized, "Authorization header missing"),
			nil,
		)
		return nil
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := ParseAccessToken(tokenString, accessTokenSecret)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			utils.FailAndAbort(c, utils.ErrTokenExpired, nil)
		case errors.Is(err, jwt.ErrTokenInvalidClaims):
			utils.FailAndAbort(c, utils.ErrTokenInvalidClaims, nil)
		case errors.Is(err, utils.ErrParsingToken):
			utils.FailAndAbort(c, utils.ErrParsingToken, err)
		case errors.Is(err, utils.ErrInvalidToken):
			utils.FailAndAbort(c, utils.ErrInvalidToken, err)
		default:
			utils.FailAndAbort(c, utils.ErrUnauthorized, err)
		}
		return nil
	}

	return claims
}

func GetClaimsFromAuthHeader(c *gin.Context, refreshTokenSecret string) *jwt.RegisteredClaims {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.FailAndAbort(
			c,
			utils.NewAPIError(http.StatusUnauthorized, "Authorization header missing"),
			nil,
		)
		return nil
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := ParseRefreshToken(tokenString, refreshTokenSecret)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			utils.FailAndAbort(c, utils.ErrTokenExpired, nil)
		case errors.Is(err, jwt.ErrTokenInvalidClaims):
			utils.FailAndAbort(c, utils.ErrTokenInvalidClaims, nil)
		case errors.Is(err, utils.ErrParsingToken):
			utils.FailAndAbort(c, utils.ErrParsingToken, err)
		case errors.Is(err, utils.ErrInvalidToken):
			utils.FailAndAbort(c, utils.ErrInvalidToken, err)
		default:
			utils.FailAndAbort(c, utils.ErrUnauthorized, err)
		}
		return nil
	}

	return claims
}

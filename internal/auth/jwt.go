package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/refine-software/afrad-api/internal/utils"
)

func GenerateToken(tokenSecret string, expInMin int) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(
			time.Now().Add(time.Minute * time.Duration(expInMin)),
		),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenSecret))
}

func GenerateAccessToken(userID, userRole, accessTokenSecret string, expInMin int) (string, error) {
	expirationTime := time.Now().Add((time.Minute * time.Duration(expInMin)))

	claims := &AccessClaims{
		Role: userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(accessTokenSecret))
}

func GenerateRefreshToken(userID, refreshTokenSecret string, expInDays int) (string, error) {
	expirationTime := time.Now().Add((time.Hour * 24) * time.Duration(expInDays))

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Subject:   userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(refreshTokenSecret))
}

// this func will parse the jwt token and return the claims,
// in case of an error we'll return one of three errors
// ErrParsingToken, ErrInvalidToken or ErrInvalidClaims
// Very self explanatory.
func ParseAccessToken(token, accessTokenSecret string) (claims *AccessClaims, err error) {
	parsedToken, err := jwt.ParseWithClaims(
		token,
		&AccessClaims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(accessTokenSecret), nil
		},
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, jwt.ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenInvalidClaims) {
			return nil, jwt.ErrTokenInvalidClaims
		}
		return nil, utils.ErrParsingToken
	}

	if !parsedToken.Valid {
		return nil, utils.ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(*AccessClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

func ParseRefreshToken(token, refreshTokenSecret string) (claims *jwt.RegisteredClaims, err error) {
	parsedToken, err := jwt.ParseWithClaims(
		token,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(refreshTokenSecret), nil
		},
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, jwt.ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenInvalidClaims) {
			return nil, jwt.ErrTokenInvalidClaims
		}
		return nil, utils.ErrParsingToken
	}

	if !parsedToken.Valid {
		return nil, utils.ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

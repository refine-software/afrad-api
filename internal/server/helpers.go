package server

import (
	"errors"
	"mime/multipart"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/auth"
	"github.com/refine-software/afrad-api/internal/models"
	"github.com/refine-software/afrad-api/internal/utils"
)

func getUserRole(email string) models.Role {
	admins := []string{
		"ali93456@gmail.com",
		"bruhgg596@gmail.com",
	}
	if slices.Contains(admins, email) {
		return models.RoleAdmin
	}
	return models.RoleUser
}

func getNameFallback(firstName, name string) string {
	if firstName == "" && name != "" {
		return name
	}
	return firstName
}

func (s *Server) generateTokens(userID, role string) (access, refresh string, err error) {
	access, err = auth.GenerateAccessToken(
		userID,
		role,
		s.env.AccessTokenSecret,
		s.env.AccessTokenExpInMin,
	)
	if err != nil {
		return
	}
	refresh, err = auth.GenerateRefreshToken(
		userID,
		s.env.RefreshTokenSecret,
		s.env.RefreshTokenExpInDays,
	)
	return
}

func getHeader(c *gin.Context, key string) string {
	header := strings.TrimSpace(c.GetHeader(key))
	if header == "" {
		utils.Fail(
			c,
			utils.ErrHeaderMissing(key),
			nil,
		)
		return ""
	}
	return header
}

func (s *Server) setCookie(c *gin.Context, cookieName, cookieVal string) {
	var secure bool
	if s.env.Environment == "prod" {
		secure = true
	}

	expTimeInSec := int((time.Hour * 24 * time.Duration(s.env.RefreshTokenExpInDays)).Seconds())

	c.SetCookie(
		cookieName,
		cookieVal,
		expTimeInSec,
		"/",
		"",
		secure,
		true,
	)
}

func (s *Server) setRefreshCookie(c *gin.Context, refreshToken string) {
	s.setCookie(c, "refresh_token", refreshToken)
}

func setEmptyCookie(c *gin.Context) {
	c.SetCookie(
		"",
		"",
		0,
		"/",
		"",
		false,
		true,
	)
}

type ImageUpload struct {
	File   multipart.File
	Header *multipart.FileHeader
}

func getImageFile(ctx *gin.Context) (*ImageUpload, error) {
	// set the form to take 5MB or less files
	const maxUpload = 5 << 20 // 5 MiB
	_ = ctx.Request.ParseMultipartForm(maxUpload)

	// fetch image from request
	file, header, err := ctx.Request.FormFile("image")
	if errors.Is(err, http.ErrMissingFile) {
		return nil, nil
	}

	if err != nil {
		return nil, &utils.APIError{
			Code:    http.StatusInternalServerError,
			Message: "unable to parse form file",
		}
	}

	if header.Size > maxUpload {
		return nil, &utils.APIError{
			Code:    utils.ErrBadRequest.Code,
			Message: "image size is bigger than allowed",
		}
	}

	contentType := header.Header.Get("Content-Type")
	if !slices.Contains([]string{"image/png", "image/jpeg", "image/webp"}, contentType) {
		return nil, &utils.APIError{
			Code:    utils.ErrBadRequest.Code,
			Message: "this type of file is not allowed",
		}
	}

	return &ImageUpload{
		File:   file,
		Header: header,
	}, nil
}

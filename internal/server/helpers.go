package server

import (
	"errors"
	"mime/multipart"
	"net/http"
	"slices"
	"strconv"
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
		s.Env.AccessTokenSecret,
		s.Env.AccessTokenExpInMin,
	)
	if err != nil {
		return
	}
	refresh, err = auth.GenerateRefreshToken(
		userID,
		s.Env.RefreshTokenSecret,
		s.Env.RefreshTokenExpInDays,
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
	if s.Env.Environment == "prod" {
		secure = true
	}

	expTimeInSec := int((time.Hour * 24 * time.Duration(s.Env.RefreshTokenExpInDays)).Seconds())

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

type ImageUpload struct {
	File   multipart.File
	Header *multipart.FileHeader
}

// this function will fetch the image from the request
// and return the file and header.
// it will return an api error,
// but be careful this function might not return an error nither a file.
// it treats the file as an optional form feild
func getImageFile(ctx *gin.Context) (*ImageUpload, *utils.APIError) {
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
			Code:    http.StatusBadRequest,
			Message: "image size is bigger than allowed",
		}
	}

	contentType := header.Header.Get("Content-Type")
	if !slices.Contains([]string{"image/png", "image/jpeg", "image/webp"}, contentType) {
		return nil, &utils.APIError{
			Code:    http.StatusBadRequest,
			Message: "this type of file is not allowed",
		}
	}

	return &ImageUpload{
		File:   file,
		Header: header,
	}, nil
}

func convStrToInt(c *gin.Context, numAsStr string, fieldName string) int {
	val, err := strconv.Atoi(numAsStr)
	if err != nil {
		utils.Fail(
			c,
			&utils.APIError{Code: http.StatusBadRequest, Message: "Invalid " + fieldName},
			err,
		)
		return 0
	}

	return val
}

func getRequiredQuery(c *gin.Context, queryName string) string {
	query := c.Query(queryName)
	if query == "" {
		utils.Fail(
			c,
			&utils.APIError{Code: http.StatusBadRequest, Message: queryName + " query required"},
			nil,
		)
		return ""
	}

	return query
}

func getRequiredQueryInt(c *gin.Context, queryName string) int {
	queryStr := getRequiredQuery(c, queryName)
	if queryStr == "" {
		return 0
	}

	query := convStrToInt(c, queryStr, queryName)
	if query == 0 {
		return 0
	}

	return query
}

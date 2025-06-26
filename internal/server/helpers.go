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

// getImageFile extracts a single image file from a multipart form request.
//
// Parameters:
//   - ctx:         The Gin context containing the HTTP request.
//   - formImgName: The name of the form field containing the image.
//   - maxUploadSize: The maximum allowed size for the uploaded image (in bytes).
//
// Behavior:
//   - Returns nil, nil if the file is not present in the form (i.e., optional field).
//   - Returns an APIError if the file exceeds the size limit, has an unsupported MIME type,
//     or if any other error occurs while parsing the form or reading the file.
//   - If successful, returns an ImageUpload struct containing the opened file and its header.
//
// Notes:
//   - This function does not close the returned file. The caller is responsible for closing it.
//   - Only "image/png", "image/jpeg", and "image/webp" content types are accepted.
func getImageFile(
	ctx *gin.Context,
	formImgName string,
	maxUploadSize int64,
) (*ImageUpload, *utils.APIError) {
	_ = ctx.Request.ParseMultipartForm(maxUploadSize)

	// fetch image from request
	file, header, err := ctx.Request.FormFile(formImgName)
	if errors.Is(err, http.ErrMissingFile) {
		return nil, nil
	}

	if err != nil {
		return nil, &utils.APIError{
			Code:    http.StatusInternalServerError,
			Message: "unable to parse form file",
		}
	}

	if header.Size > maxUploadSize {
		return nil, &utils.APIError{
			Code:    http.StatusBadRequest,
			Message: "image size is bigger than allowed",
		}
	}

	allowedTypes := map[string]bool{
		"image/png":  true,
		"image/jpeg": true,
		"image/webp": true,
	}

	contentType := header.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
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

// getImageFiles retrieves multiple image files from a multipart form request.
//
// Parameters:
//   - ctx:            The Gin context containing the HTTP request.
//   - formImgName:    The name of the form field containing the image files.
//   - maxUploadSize:  The maximum allowed size for each uploaded image (in bytes).
//
// Behavior:
//   - Parses the multipart form data from the request.
//   - Returns nil, nil if no files are provided under the specified form field (i.e., optional).
//   - Validates each file for size and MIME type ("image/png", "image/jpeg", "image/webp").
//   - Opens each valid image file and returns a slice of ImageUpload structs.
//   - Returns an APIError if any file fails validation or opening.
//
// Notes:
//   - The caller is responsible for closing all returned files to prevent memory/resource leaks.
//   - This function stops and returns an error as soon as one invalid or unreadable file is encountered.
//   - The total memory used for parsing is limited to maxUploadSize.
func getImageFiles(
	ctx *gin.Context,
	formImgName string,
	maxUploadSize int64,
) ([]ImageUpload, *utils.APIError) {
	err := ctx.Request.ParseMultipartForm(maxUploadSize)
	if err != nil {
		return nil, &utils.APIError{
			Code:    http.StatusBadRequest,
			Message: "unable to parse multipart form",
		}
	}

	form := ctx.Request.MultipartForm
	if form == nil || form.File == nil {
		return nil, nil
	}

	files := form.File[formImgName]
	if len(files) == 0 {
		return nil, nil
	}

	allowedTypes := map[string]bool{
		"image/png":  true,
		"image/jpeg": true,
		"image/webp": true,
	}

	var uploads []ImageUpload
	for _, header := range files {
		if header.Size > maxUploadSize {
			return nil, &utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "one of the images exceeds the max size of 5MB",
			}
		}

		contentType := header.Header.Get("Content-Type")
		if !allowedTypes[contentType] {
			return nil, &utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "one of the files is not an allowed image type",
			}
		}

		file, err := header.Open()
		if err != nil {
			return nil, &utils.APIError{
				Code:    http.StatusInternalServerError,
				Message: "unable to open uploaded image file",
			}
		}

		uploads = append(uploads, ImageUpload{
			File:   file,
			Header: header,
		})
	}

	return uploads, nil
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

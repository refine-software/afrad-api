package test

import (
	"mime/multipart"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/config"
	"github.com/refine-software/afrad-api/internal/database"
	"github.com/refine-software/afrad-api/internal/server"
)

type MockS3 struct{}

func (m *MockS3) UploadImage(
	ctx *gin.Context,
	file multipart.File,
	header *multipart.FileHeader,
) (string, error) {
	return "https://mock-bucket/image.jpg", nil
}

func (m *MockS3) DeleteImageByURL(ctx *gin.Context, url string) error {
	return nil
}

type MockEmail struct{}

func (m *MockEmail) SendOtpEmail(userEmail, otp string) error {
	return nil
}

func setupTestServer(t *testing.T) *gin.Engine {
	t.Helper()

	env := config.NewTestEnv()
	db := database.New(env)

	s := &server.Server{
		Env:   env,
		DB:    db,
		S3:    &MockS3{},
		Email: &MockEmail{},
	}
	gin.SetMode(gin.TestMode)
	return s.RegisterRoutes().(*gin.Engine)
}

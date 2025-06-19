package test

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	tests := []struct {
		name       string
		fields     map[string]string
		wantStatus int
		wantBody   string
	}{
		// ✅ Success Cases
		{
			name: "Valid user without phone",
			fields: map[string]string{
				"firstName": "Ali",
				"lastName":  "Test",
				"email":     "ali1@example.com",
				"password":  "supersecure123",
			},
			wantStatus: http.StatusCreated,
			wantBody:   "user created",
		},
		{
			name: "Valid user with phone",
			fields: map[string]string{
				"firstName":   "Ali",
				"lastName":    "Test",
				"email":       "ali2@example.com",
				"password":    "supersecure123",
				"phoneNumber": "1234567890",
			},
			wantStatus: http.StatusCreated,
			wantBody:   "user created",
		},

		// ❌ Missing Required Fields
		{
			name: "Missing firstName",
			fields: map[string]string{
				"lastName": "Test",
				"email":    "missing1@example.com",
				"password": "supersecure123",
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "Invalid",
		},
		{
			name: "Missing lastName",
			fields: map[string]string{
				"firstName": "Ali",
				"email":     "missing2@example.com",
				"password":  "supersecure123",
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "Invalid",
		},
		{
			name: "Missing email",
			fields: map[string]string{
				"firstName": "Ali",
				"lastName":  "Test",
				"password":  "supersecure123",
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "Invalid",
		},
		{
			name: "Missing password",
			fields: map[string]string{
				"firstName": "Ali",
				"lastName":  "Test",
				"email":     "missing4@example.com",
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "Invalid",
		},

		// ❌ Invalid Formats
		{
			name: "Invalid email format",
			fields: map[string]string{
				"firstName": "Ali",
				"lastName":  "Test",
				"email":     "not-an-email",
				"password":  "supersecure123",
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "email",
		},
		{
			name: "Empty password",
			fields: map[string]string{
				"firstName": "Ali",
				"lastName":  "Test",
				"email":     "empty@example.com",
				"password":  "",
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "Invalid",
		},

		// ❌ Duplicate User
		{
			name: "Duplicate email registration",
			fields: map[string]string{
				"firstName": "Ali",
				"lastName":  "Test",
				"email":     "ali2@example.com",
				"password":  "supersecure123",
			},
			wantStatus: http.StatusConflict,
			wantBody:   "user already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestServer(t)

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			for k, v := range tt.fields {
				_ = writer.WriteField(k, v)
			}

			writer.Close()

			req, _ := http.NewRequest("POST", "/auth/register", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.wantStatus, resp.Code, "status code mismatch")
			assert.Contains(t, resp.Body.String(), tt.wantBody, "body mismatch")
		})
	}
}

func TestRegisterUserWithImage(t *testing.T) {
	router := setupTestServer(t)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("firstName", "Ali")
	_ = writer.WriteField("lastName", "Test")
	_ = writer.WriteField("email", "aliimage@example.com")
	_ = writer.WriteField("password", "secure123")

	// Create a real in-memory JPEG image
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255})

	var imgBuf bytes.Buffer
	err := jpeg.Encode(&imgBuf, img, nil)
	assert.NoError(t, err)

	// Manually set Content-Type for image
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="image"; filename="test.jpg"`)
	h.Set("Content-Type", "image/jpeg")

	part, err := writer.CreatePart(h)
	assert.NoError(t, err)
	_, err = part.Write(imgBuf.Bytes())
	assert.NoError(t, err)

	writer.Close()

	req, _ := http.NewRequest("POST", "/auth/register", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code, resp.Body.String())
	assert.Contains(t, resp.Body.String(), "user created")
}

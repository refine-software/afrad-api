// Package server handles starting the HTTP server and routing.
package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/refine-software/afrad-api/config"
	"github.com/refine-software/afrad-api/internal/auth"
	"github.com/refine-software/afrad-api/internal/database"
	"github.com/refine-software/afrad-api/internal/s3"
)

type Server struct {
	port int

	DB    database.Service
	Env   *config.Env
	S3    s3.S3
	Email auth.EmailSender
}

func NewServer() *http.Server {
	env := config.NewEnv()

	if env.Environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	s3Storage, err := s3.NewS3Storage(env)
	if err != nil {
		log.Println("couldn't create s3 storage")
		log.Fatalln(err)
	}

	auth.InitOauth(env)

	NewServer := &Server{
		port: env.Port,

		DB:    database.New(env),
		Env:   env,
		S3:    s3Storage,
		Email: auth.NewEmailService(env.Email, env.Password),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/refine-software/afrad-api/config"
	"github.com/refine-software/afrad-api/internal/auth"
	"github.com/refine-software/afrad-api/internal/database"
	"github.com/refine-software/afrad-api/internal/s3"
)

type Server struct {
	port int

	db     database.Service
	env    *config.Env
	s3     *s3.S3Storage
	twilio auth.Twilio
}

func NewServer() *http.Server {
	env := config.NewEnv()

	s3Storage, err := s3.NewS3Storage(env)
	if err != nil {
		log.Println("couldn't create s3 storage")
		log.Fatalln(err)
	}

	auth.InitOauth(env)

	NewServer := &Server{
		port: env.Port,

		db:     database.New(env),
		env:    env,
		s3:     s3Storage,
		twilio: auth.NewTwilio(env),
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

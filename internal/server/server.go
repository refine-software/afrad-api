package server

import (
	"fmt"
	"net/http"
	"time"

	"afrad-api/config"
	"afrad-api/internal/database"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int

	db  database.Service
	env *config.Env
}

func NewServer() *http.Server {
	env := config.NewEnv()

	NewServer := &Server{
		port: env.Port,

		db:  database.New(env),
		env: env,
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

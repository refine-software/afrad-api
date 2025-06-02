package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/middleware"

	"github.com/coder/websocket"
)

func (s *Server) RegisterRoutes() http.Handler {
	engine := gin.Default()
	engine.Use(middleware.SetupCors())

	engine.GET("/websocket", s.websocketHandler)

	s.registerPublicRoutes(engine)
	s.registerUserRoutes(engine)
	s.registerAdminRoutes(engine)

	return engine
}

func (s *Server) registerPublicRoutes(e *gin.Engine) {
	auth := e.Group("/auth")
}

func (s *Server) registerUserRoutes(e *gin.Engine) {
}

func (s *Server) registerAdminRoutes(e *gin.Engine) {
}

func (s *Server) websocketHandler(c *gin.Context) {
	w := c.Writer
	r := c.Request
	socket, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Printf("could not open websocket: %v", err)
		_, _ = w.Write([]byte("could not open websocket"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer socket.Close(websocket.StatusGoingAway, "server closing websocket")

	ctx := r.Context()
	socketCtx := socket.CloseRead(ctx)

	for {
		payload := fmt.Sprintf("server timestamp: %d", time.Now().UnixNano())
		err := socket.Write(socketCtx, websocket.MessageText, []byte(payload))
		if err != nil {
			break
		}
		time.Sleep(time.Second * 2)
	}
}

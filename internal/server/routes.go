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
	oauth := e.Group("/oauth")
	{
		oauth.GET("/google/login", s.loginWithGoogle)
		oauth.GET("/google/callback", s.googleCallback)
		oauth.POST("/refresh", s.refreshTokenOauth)
		oauth.POST("/logout")
		oauth.GET("/me")
	}

	auth := e.Group("/auth")
	{
		auth.POST("/register")
		auth.POST("/verify-phone-number")
		auth.POST("/resend-verification-otp")
		auth.POST("/login")
		auth.POST("/reset-password")
		auth.POST("/reset-password/confirm")
		auth.POST("/refresh-tokens")
	}

	products := e.Group("/products")
	{
		products.GET("")
		products.GET("/:id")
		products.GET("/categories")
	}
}

func (s *Server) registerUserRoutes(e *gin.Engine) {
	user := e.Group("/user")
	{
		user.GET("")
		user.PUT("")
		user.DELETE("")
		user.POST("/review")
		user.PUT("/review")
		user.GET("/review")
		user.PATCH("/user/notificatoin-preferences")
		user.POST("/logout")
		user.POST("/logout-all")
	}

	cart := e.Group("/cart")
	{
		cart.GET("")
		cart.POST("/item")
		cart.PATCH("/:id")
		cart.DELETE("/:id")
	}

	wishlist := e.Group("/wishlist")
	{
		wishlist.GET("")
		wishlist.POST("/:id")
		wishlist.DELETE("/:id")
	}

	orders := e.Group("/orders")
	{
		orders.GET("")
		orders.POST("")
		orders.GET("/:id")
		orders.PATCH("/:id/cancel")
	}
}

func (s *Server) registerAdminRoutes(e *gin.Engine) {
	admin := e.Group("/admin")

	product := admin.Group("/product")
	{
		product.PUT("/:id")
		product.DELETE("/:id")
		product.POST("")
	}

	category := admin.Group("/category")
	{
		category.POST("")
		category.PATCH("/:id")
		category.DELETE("/:id")
	}

	discount := admin.Group("/discounts")
	{
		discount.POST("/product")
		discount.POST("/variant")
		discount.GET("") // ??
		discount.PUT("/:id")
		discount.DELETE("/:id")
	}

	orders := admin.Group("/orders")
	{
		orders.GET("")
	}
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

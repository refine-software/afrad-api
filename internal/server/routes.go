package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/middleware"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/coder/websocket"
	swaggerFiles "github.com/swaggo/files"
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
	}

	auth := e.Group("/auth")
	{
		auth.POST("/register", s.register)
		auth.POST("/verify-account", s.verifyAccount)
		auth.POST("/resend-verification", s.resendVerification)
		auth.POST("/login", s.login)
		auth.POST("/reset-password", s.passwordReset)
		auth.POST("/reset-password/confirm", s.resetPasswordConfirm)
		auth.POST("/refresh", s.refreshTokens)
	}

	products := e.Group("/products")
	{
		products.GET("", s.getAllProducts)
		products.GET("/:id", s.getProduct)
		products.GET("/categories", s.getCategories)
	}
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (s *Server) registerUserRoutes(e *gin.Engine) {
	protected := e.Group("")
	protected.Use(middleware.AuthRequired(s.Env.AccessTokenSecret))

	user := protected.Group("/user")
	{
		user.GET("", s.getUser)
		user.PUT("", s.updateUser)
		user.DELETE("", s.deleteUser)
		user.POST("/review", s.postReview)
		user.PUT("/review")
		user.GET("/review")
		user.PATCH("/user/notificatoin-preferences")
		user.POST("/logout", s.logout)
		user.POST("/logout/all", s.logoutFromAllSessions)
	}

	cart := protected.Group("/cart")
	{
		cart.GET("")
		cart.POST("/item")
		cart.PATCH("/:id")
		cart.DELETE("/:id")
	}

	wishlist := protected.Group("/wishlist")
	{
		wishlist.GET("")
		wishlist.POST("/:id")
		wishlist.DELETE("/:id")
	}

	orders := protected.Group("/orders")
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
		category.POST("", s.createCategory)
		category.PATCH("/:id", s.updateCategory)
		category.DELETE("/:id", s.deleteCategory)
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

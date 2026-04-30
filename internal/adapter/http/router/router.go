package router

import (
	"jsunnykhan/go-clean-template/internal/core/config"
	"log/slog"

	"jsunnykhan/go-clean-template/internal/adapter/http/handler"

	"jsunnykhan/go-clean-template/internal/adapter/http/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	userH *handler.UserHandler,
	jwt config.JWTConfig,
	log *slog.Logger,
) *gin.Engine {

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	// ------health check------
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// ------Group routes------
	api := r.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", userH.Register)
		auth.POST("/login", userH.Login)
		auth.POST("/refresh", userH.Refresh)
	}

	// ------Protected routes------
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(jwt.SecretKey))
	{
		protected.GET("/auth/me", func(ctx *gin.Context) {})
	}

	return r

}

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/auth"
	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/bootstrap"
	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/config"
	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/handler"
	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/k8s"
)

func main() {
	cfg := config.Load()

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	clientset, err := k8s.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}
	log.Println("Connected to Kubernetes")

	jwtSecret, err := k8s.GetJWTKey(context.Background(), clientset, cfg.Namespace)
	if err != nil {
		log.Fatalf("Failed to get JWT secret: %v", err)
	}
	log.Println("JWT secret initialized")

	ctx := context.Background()
	if err := bootstrap.InitializeAdmin(ctx, clientset, cfg.Namespace, cfg.InitAdminUsername, cfg.InitAdminPassword); err != nil {
		log.Fatalf("Failed to initialize admin: %v", err)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	authHandler := handler.NewAuthHandler(clientset, cfg.Namespace, jwtSecret, cfg.JWTExpiry)
	userHandler := handler.NewUserHandler(clientset, cfg.Namespace)

	public := router.Group("/api")
	{
		public.POST("/auth/login", auth.RateLimitMiddleware(1, 5), authHandler.Login)
	}

	protected := router.Group("/api")
	protected.Use(auth.AuthMiddleware(jwtSecret))
	{
		protected.GET("/users", userHandler.ListUsers)
		protected.POST("/users", userHandler.CreateUser)
		protected.DELETE("/users/:username", userHandler.DeleteUser)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	addr := ":" + cfg.Port
	log.Printf("Server starting on %s", addr)

	go func() {
		if err := router.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}

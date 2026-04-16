package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/auth"
	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/k8s"
	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/model"
	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/validator"
	"k8s.io/client-go/kubernetes"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	clientset *kubernetes.Clientset
	namespace string
	jwtSecret string
	jwtExpiry time.Duration
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(clientset *kubernetes.Clientset, namespace, jwtSecret string, jwtExpiry time.Duration) *AuthHandler {
	return &AuthHandler{
		clientset: clientset,
		namespace: namespace,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validator.ValidateUsername(req.Username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userSecret, err := k8s.GetUserSecret(c.Request.Context(), h.clientset, h.namespace, req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !auth.CheckPassword(req.Password, userSecret.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := auth.GenerateToken(req.Username, h.jwtSecret, h.jwtExpiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, model.LoginResponse{
		Token: token,
		User: model.User{
			Username:  req.Username,
			CreatedAt: userSecret.CreatedAt,
		},
	})
}

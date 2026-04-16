package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/auth"
	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/k8s"
	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/model"
	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/validator"
	"k8s.io/client-go/kubernetes"
)

// UserHandler handles user management endpoints.
type UserHandler struct {
	clientset *kubernetes.Clientset
	namespace string
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(clientset *kubernetes.Clientset, namespace string) *UserHandler {
	return &UserHandler{
		clientset: clientset,
		namespace: namespace,
	}
}

// ListUsers handles GET /api/users
func (h *UserHandler) ListUsers(c *gin.Context) {
	usernames, err := k8s.ListUserSecrets(c.Request.Context(), h.clientset, h.namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}

	users := make([]model.User, 0, len(usernames))
	for _, username := range usernames {
		userSecret, err := k8s.GetUserSecret(c.Request.Context(), h.clientset, h.namespace, username)
		if err == nil {
			users = append(users, model.User{
				Username:  username,
				CreatedAt: userSecret.CreatedAt,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": users})
}

// CreateUser handles POST /api/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validator.ValidateUsername(req.Username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 8 characters"})
		return
	}

	if len(req.Password) > 128 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at most 128 characters"})
		return
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	if err := k8s.CreateUserSecret(c.Request.Context(), h.clientset, h.namespace, req.Username, passwordHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	userSecret, err := k8s.GetUserSecret(c.Request.Context(), h.clientset, h.namespace, req.Username)
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"data": model.User{
				Username: req.Username,
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": model.User{
			Username:  req.Username,
			CreatedAt: userSecret.CreatedAt,
		},
	})
}

// DeleteUser handles DELETE /api/users/:username
func (h *UserHandler) DeleteUser(c *gin.Context) {
	targetUsername := c.Param("username")

	currentUsername, exists := auth.GetUsername(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if targetUsername == currentUsername {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete your own account"})
		return
	}

	usernames, err := k8s.ListUserSecrets(c.Request.Context(), h.clientset, h.namespace)
	if err == nil && len(usernames) <= 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete the last user"})
		return
	}

	if err := k8s.DeleteUserSecret(c.Request.Context(), h.clientset, h.namespace, targetUsername); err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "user deleted"})
}

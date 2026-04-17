package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/k8s"
	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/model"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// HtpasswdHandler handles htpasswd secret management endpoints.
type HtpasswdHandler struct {
	clientset *kubernetes.Clientset
	dynClient dynamic.Interface
}

// NewHtpasswdHandler creates a new HtpasswdHandler.
func NewHtpasswdHandler(clientset *kubernetes.Clientset, dynClient dynamic.Interface) *HtpasswdHandler {
	return &HtpasswdHandler{
		clientset: clientset,
		dynClient: dynClient,
	}
}

// ListHtpasswdSecrets handles GET /api/namespaces/:namespace/htpasswd
func (h *HtpasswdHandler) ListHtpasswdSecrets(c *gin.Context) {
	namespace := c.Param("namespace")
	ctx := c.Request.Context()

	secrets, err := k8s.ListHtpasswdSecrets(ctx, h.clientset, namespace)
	if err != nil {
		log.Printf("Failed to list htpasswd secrets: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list htpasswd secrets"})
		return
	}

	summaries := make([]model.HtpasswdSecretSummary, 0, len(secrets))
	for _, s := range secrets {
		summaries = append(summaries, model.HtpasswdSecretSummary{
			Name:      s.Name,
			Namespace: s.Namespace,
			UserCount: len(s.Users),
			CreatedAt: s.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": summaries})
}

// GetHtpasswdSecret handles GET /api/namespaces/:namespace/htpasswd/:name
func (h *HtpasswdHandler) GetHtpasswdSecret(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	ctx := c.Request.Context()

	secret, err := k8s.GetHtpasswdSecret(ctx, h.clientset, namespace, name)
	if err != nil {
		if errors.Is(err, k8s.ErrHtpasswdNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "htpasswd secret not found"})
			return
		}
		log.Printf("Failed to get htpasswd secret: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get htpasswd secret"})
		return
	}

	// Find linked SecurityPolicies
	linkedPolicies := h.findLinkedPolicies(c.Request.Context(), namespace, name)

	users := make([]model.HtpasswdUserEntry, 0, len(secret.Users))
	for _, u := range secret.Users {
		users = append(users, model.HtpasswdUserEntry{Username: u})
	}

	detail := model.HtpasswdSecretDetail{
		Name:                   secret.Name,
		Namespace:              secret.Namespace,
		Users:                  users,
		UserCount:              len(users),
		CreatedAt:              secret.CreatedAt,
		LinkedSecurityPolicies: linkedPolicies,
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": detail})
}

// CreateHtpasswdSecret handles POST /api/namespaces/:namespace/htpasswd
func (h *HtpasswdHandler) CreateHtpasswdSecret(c *gin.Context) {
	namespace := c.Param("namespace")
	ctx := c.Request.Context()

	var req model.CreateHtpasswdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Check for duplicate usernames in request
	seen := make(map[string]bool, len(req.Users))
	for _, u := range req.Users {
		if seen[u.Username] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "duplicate username: " + u.Username})
			return
		}
		seen[u.Username] = true
	}

	users := make(map[string]string, len(req.Users))
	for _, u := range req.Users {
		users[u.Username] = u.Password
	}

	if err := k8s.CreateHtpasswdSecret(ctx, h.clientset, namespace, req.Name, users); err != nil {
		log.Printf("Failed to create htpasswd secret: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create htpasswd secret"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": gin.H{
			"name":      req.Name,
			"namespace": namespace,
			"userCount": len(users),
		},
	})
}

// AddUser handles POST /api/namespaces/:namespace/htpasswd/:name/users
func (h *HtpasswdHandler) AddUser(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	ctx := c.Request.Context()

	var req model.AddHtpasswdUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := k8s.AddUserToHtpasswd(ctx, h.clientset, namespace, name, req.Username, req.Password); err != nil {
		if errors.Is(err, k8s.ErrHtpasswdNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "htpasswd secret not found"})
			return
		}
		log.Printf("Failed to add user to htpasswd: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "user added"})
}

// RemoveUser handles DELETE /api/namespaces/:namespace/htpasswd/:name/users/:username
func (h *HtpasswdHandler) RemoveUser(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	username := c.Param("username")
	ctx := c.Request.Context()

	if err := k8s.RemoveUserFromHtpasswd(ctx, h.clientset, namespace, name, username); err != nil {
		if errors.Is(err, k8s.ErrHtpasswdNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "htpasswd secret not found"})
			return
		}
		if errors.Is(err, k8s.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		log.Printf("Failed to remove user from htpasswd: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "user removed"})
}

// DeleteHtpasswdSecret handles DELETE /api/namespaces/:namespace/htpasswd/:name
func (h *HtpasswdHandler) DeleteHtpasswdSecret(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	ctx := c.Request.Context()

	if err := k8s.DeleteHtpasswdSecret(ctx, h.clientset, namespace, name); err != nil {
		if errors.Is(err, k8s.ErrHtpasswdNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "htpasswd secret not found"})
			return
		}
		log.Printf("Failed to delete htpasswd secret: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete htpasswd secret"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "htpasswd secret deleted"})
}

// findLinkedPolicies queries SecurityPolicies in the namespace that reference the given htpasswd secret.
func (h *HtpasswdHandler) findLinkedPolicies(ctx context.Context, namespace, secretName string) []model.LinkedSecurityPolicy {
	policies, err := k8s.ListSecurityPolicies(ctx, h.dynClient, namespace)
	if err != nil {
		log.Printf("Failed to list security policies for association: %v", err)
		return nil
	}

	var linked []model.LinkedSecurityPolicy
	for _, p := range policies {
		if p.BasicAuthSecretName == secretName {
			linked = append(linked, model.LinkedSecurityPolicy{
				Name:      p.Name,
				Namespace: p.Namespace,
				TargetRef: model.PolicyTargetRef{
					Name: p.TargetRef.Name,
					Kind: p.TargetRef.Kind,
				},
			})
		}
	}
	return linked
}

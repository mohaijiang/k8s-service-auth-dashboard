package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/k8s"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// ServiceHandler handles service overview endpoints.
type ServiceHandler struct {
	clientset *kubernetes.Clientset
	dynClient dynamic.Interface
}

// NewServiceHandler creates a new ServiceHandler.
func NewServiceHandler(clientset *kubernetes.Clientset, dynClient dynamic.Interface) *ServiceHandler {
	return &ServiceHandler{
		clientset: clientset,
		dynClient: dynClient,
	}
}

// ListServices handles GET /api/services
func (h *ServiceHandler) ListServices(c *gin.Context) {
	namespace := c.Query("namespace")
	ctx := c.Request.Context()

	services, err := k8s.ListServices(ctx, h.clientset, namespace)
	if err != nil {
		log.Printf("Failed to list services: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list services"})
		return
	}

	httpRoutes, err := k8s.ListHTTPRoutes(ctx, h.dynClient, namespace)
	if err != nil {
		log.Printf("Failed to list httproutes (CRD may not be installed): %v", err)
		httpRoutes = []k8s.HTTPRouteData{}
	}

	securityPolicies, err := k8s.ListSecurityPolicies(ctx, h.dynClient, namespace)
	if err != nil {
		log.Printf("Failed to list securitypolicies (CRD may not be installed): %v", err)
		securityPolicies = []k8s.SecurityPolicyData{}
	}

	overviews := BuildServiceOverviews(services, httpRoutes, securityPolicies)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    overviews,
	})
}

// ListNamespaces handles GET /api/namespaces
func (h *ServiceHandler) ListNamespaces(c *gin.Context) {
	namespaces, err := k8s.ListNamespaces(c.Request.Context(), h.clientset)
	if err != nil {
		log.Printf("Failed to list namespaces: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list namespaces"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    namespaces,
	})
}

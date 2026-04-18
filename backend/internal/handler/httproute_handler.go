package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/k8s"
	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/model"
	"k8s.io/client-go/dynamic"
)

// HTTPRouteHandler handles HTTPRoute-related endpoints.
type HTTPRouteHandler struct {
	dynClient dynamic.Interface
}

// NewHTTPRouteHandler creates a new HTTPRouteHandler.
func NewHTTPRouteHandler(dynClient dynamic.Interface) *HTTPRouteHandler {
	return &HTTPRouteHandler{
		dynClient: dynClient,
	}
}

// ListByService handles GET /api/services/:namespace/:service/httproutes
func (h *HTTPRouteHandler) ListByService(c *gin.Context) {
	namespace := c.Param("namespace")
	serviceName := c.Param("service")
	ctx := c.Request.Context()

	routes, err := k8s.GetHTTPRoutesByService(ctx, h.dynClient, namespace, serviceName)
	if err != nil {
		log.Printf("Failed to list httproutes for service %s/%s: %v", namespace, serviceName, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list httproutes"})
		return
	}

	// Convert to detail format
	details := make([]model.HTTPRouteDetail, 0, len(routes))
	for _, route := range routes {
		parentRefs := make([]model.ParentRef, 0, len(route.ParentRefs))
		for _, ref := range route.ParentRefs {
			parentRefs = append(parentRefs, model.ParentRef{
				Name:      ref.Name,
				Namespace: ref.Namespace,
			})
		}

		details = append(details, model.HTTPRouteDetail{
			Name:       route.Name,
			Namespace:  route.Namespace,
			Hostnames:  route.Hostnames,
			ParentRefs: parentRefs,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    details,
	})
}

// Create handles POST /api/httproutes
func (h *HTTPRouteHandler) Create(c *gin.Context) {
	var req model.HTTPRouteCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// Build parent refs from request
	parentRefs := make([]k8s.ParentRefData, 0, len(req.ParentRefs))
	for _, ref := range req.ParentRefs {
		parentRefs = append(parentRefs, k8s.ParentRefData{
			Name:      ref.Name,
			Namespace: ref.Namespace,
		})
	}

	spec := k8s.HTTPRouteSpec{
		Hostnames:   req.Hostnames,
		ServiceName: req.ServiceName,
		ServicePort: req.ServicePort,
		ParentRefs:  parentRefs,
	}

	err := k8s.CreateHTTPRoute(ctx, h.dynClient, req.Namespace, req.Name, spec)
	if err != nil {
		log.Printf("Failed to create httproute %s/%s: %v", req.Namespace, req.Name, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create httproute"})
		return
	}

	// Create SecurityPolicy if requested
	if req.SecurityPolicy != nil {
		securityPolicySpec := k8s.SecurityPolicySpec{
			TargetHTTPRouteName: req.Name,
			BasicAuthSecretName:   req.SecurityPolicy.BasicAuthSecretName,
		}

		err = k8s.CreateSecurityPolicy(ctx, h.dynClient, req.Namespace, req.Name, securityPolicySpec)
		if err != nil {
			log.Printf("Failed to create securitypolicy %s/%s: %v", req.Namespace, req.Name, err)
			// Rollback: delete the created HTTPRoute
			_ = k8s.DeleteHTTPRoute(ctx, h.dynClient, req.Namespace, req.Name)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create securitypolicy"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": model.HTTPRouteDetail{
			Name:        req.Name,
			Namespace:   req.Namespace,
			Hostnames:   req.Hostnames,
			ParentRefs:  req.ParentRefs,
			ServiceName: req.ServiceName,
		},
	})
}

// Delete handles DELETE /api/httproutes/:namespace/:name
func (h *HTTPRouteHandler) Delete(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	ctx := c.Request.Context()

	// Delete the associated SecurityPolicy if it exists (ignore errors if not found)
	_ = k8s.DeleteSecurityPolicy(ctx, h.dynClient, namespace, name)

	err := k8s.DeleteHTTPRoute(ctx, h.dynClient, namespace, name)
	if err != nil {
		log.Printf("Failed to delete httproute %s/%s: %v", namespace, name, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete httproute"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "HTTPRoute deleted successfully",
	})
}

package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/k8s"
	"k8s.io/client-go/dynamic"
)

// GatewayHandler handles gateway-related endpoints.
type GatewayHandler struct {
	dynClient dynamic.Interface
}

// NewGatewayHandler creates a new GatewayHandler.
func NewGatewayHandler(dynClient dynamic.Interface) *GatewayHandler {
	return &GatewayHandler{
		dynClient: dynClient,
	}
}

// ListGateways handles GET /api/gateways
func (h *GatewayHandler) ListGateways(c *gin.Context) {
	namespace := c.Query("namespace")
	ctx := c.Request.Context()

	gateways, err := k8s.ListGateways(ctx, h.dynClient, namespace)
	if err != nil {
		log.Printf("Failed to list gateways (CRD may not be installed): %v", err)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    []k8s.GatewayData{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    gateways,
	})
}

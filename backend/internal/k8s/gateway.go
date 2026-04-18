package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// GatewayData holds parsed Gateway information.
type GatewayData struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

var gatewayGVR = schema.GroupVersionResource{
	Group:    "gateway.networking.k8s.io",
	Version:  "v1",
	Resource: "gateways",
}

// ListGateways lists Gateway resources, optionally filtered by namespace.
func ListGateways(ctx context.Context, dynClient dynamic.Interface, namespace string) ([]GatewayData, error) {
	list, err := dynClient.Resource(gatewayGVR).Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list gateways: %w", err)
	}

	result := make([]GatewayData, 0, len(list.Items))
	for i := range list.Items {
		item := &list.Items[i]
		result = append(result, GatewayData{
			Name:      item.GetName(),
			Namespace: item.GetNamespace(),
		})
	}
	return result, nil
}

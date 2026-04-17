package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ServiceInfo holds the K8s Service data needed for the overview.
type ServiceInfo struct {
	Name      string
	Namespace string
	ClusterIP string
	Ports     []corev1.ServicePort
	Selector  map[string]string
}

// ListServices lists Services, optionally filtered by namespace.
// If namespace is empty, lists across all namespaces.
func ListServices(ctx context.Context, clientset *kubernetes.Clientset, namespace string) ([]ServiceInfo, error) {
	services, err := clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	result := make([]ServiceInfo, 0, len(services.Items))
	for _, svc := range services.Items {
		result = append(result, ServiceInfo{
			Name:      svc.Name,
			Namespace: svc.Namespace,
			ClusterIP: svc.Spec.ClusterIP,
			Ports:     svc.Spec.Ports,
			Selector:  svc.Spec.Selector,
		})
	}
	return result, nil
}

// ListNamespaces returns all namespace names.
func ListNamespaces(ctx context.Context, clientset *kubernetes.Clientset) ([]string, error) {
	namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	result := make([]string, 0, len(namespaces.Items))
	for _, ns := range namespaces.Items {
		result = append(result, ns.Name)
	}
	return result, nil
}

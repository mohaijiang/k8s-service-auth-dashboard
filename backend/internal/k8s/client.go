package k8s

import (
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewConfig creates a Kubernetes rest.Config with dual-mode credential loading.
func NewConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		log.Printf("Using local kubeconfig: %s", kubeconfig)
	} else {
		log.Println("Using in-cluster configuration")
	}
	return config, nil
}

// NewClient creates a Kubernetes typed client.
func NewClient() (*kubernetes.Clientset, error) {
	config, err := NewConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

// NewDynamicClient creates a Kubernetes dynamic client for CRD access.
func NewDynamicClient() (dynamic.Interface, error) {
	config, err := NewConfig()
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(config)
}

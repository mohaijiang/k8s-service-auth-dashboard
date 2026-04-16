package k8s

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	JWTSecretName = "dashboard-jwt-secret"
	JWTSecretKey  = "jwt-key"
)

// GetJWTKey retrieves or creates the JWT signing key from K8s Secret.
func GetJWTKey(ctx context.Context, clientset *kubernetes.Clientset, namespace string) (string, error) {
	secret, err := clientset.CoreV1().Secrets(namespace).Get(ctx, JWTSecretName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return createJWTSecret(ctx, clientset, namespace)
		}
		return "", fmt.Errorf("failed to get JWT secret: %w", err)
	}

	keyBytes, ok := secret.Data[JWTSecretKey]
	if !ok || len(keyBytes) == 0 {
		return "", fmt.Errorf("JWT secret exists but key is empty")
	}

	return string(keyBytes), nil
}

func createJWTSecret(ctx context.Context, clientset *kubernetes.Clientset, namespace string) (string, error) {
	keyBytes := make([]byte, 32)
	if _, err := rand.Read(keyBytes); err != nil {
		return "", fmt.Errorf("failed to generate random key: %w", err)
	}
	keyHex := hex.EncodeToString(keyBytes)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      JWTSecretName,
			Namespace: namespace,
			Labels: map[string]string{
				"app.kubernetes.io/part-of": "k8s-service-auth-dashboard",
			},
		},
		Data: map[string][]byte{
			JWTSecretKey: []byte(keyHex),
		},
	}

	created, err := clientset.CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create JWT secret: %w", err)
	}

	log.Printf("Created JWT secret: %s/%s", namespace, created.Name)
	return keyHex, nil
}

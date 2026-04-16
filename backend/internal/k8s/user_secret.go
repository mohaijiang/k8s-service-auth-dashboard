package k8s

import (
	"context"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	UserSecretLabelKey   = "app.kubernetes.io/part-of"
	UserSecretLabelValue = "k8s-service-auth-dashboard"
)

// UserSecretData represents the data stored in a user Secret.
type UserSecretData struct {
	PasswordHash string
	CreatedAt    time.Time
}

// GetUserSecret retrieves a user Secret by username.
func GetUserSecret(ctx context.Context, clientset *kubernetes.Clientset, namespace, username string) (*UserSecretData, error) {
	secretName := UserSecretName(username)
	secret, err := clientset.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user secret: %w", err)
	}

	return parseUserSecret(secret)
}

// CreateUserSecret creates a new user Secret.
func CreateUserSecret(ctx context.Context, clientset *kubernetes.Clientset, namespace, username, passwordHash string) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      UserSecretName(username),
			Namespace: namespace,
			Labels: map[string]string{
				UserSecretLabelKey: UserSecretLabelValue,
				"user":             username,
			},
		},
		StringData: map[string]string{
			"passwordHash": passwordHash,
			"createdAt":    time.Now().UTC().Format(time.RFC3339),
		},
	}

	_, err := clientset.CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create user secret: %w", err)
	}

	return nil
}

// DeleteUserSecret deletes a user Secret.
func DeleteUserSecret(ctx context.Context, clientset *kubernetes.Clientset, namespace, username string) error {
	secretName := UserSecretName(username)
	err := clientset.CoreV1().Secrets(namespace).Delete(ctx, secretName, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to delete user secret: %w", err)
	}
	return nil
}

// ListUserSecrets lists all user Secrets in the namespace.
func ListUserSecrets(ctx context.Context, clientset *kubernetes.Clientset, namespace string) ([]string, error) {
	secrets, err := clientset.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", UserSecretLabelKey, UserSecretLabelValue),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list user secrets: %w", err)
	}

	usernames := make([]string, 0, len(secrets.Items))
	for _, secret := range secrets.Items {
		username := secret.Labels["user"]
		if username != "" {
			usernames = append(usernames, username)
		}
	}

	return usernames, nil
}

// UserSecretName returns the Secret name for a given username.
func UserSecretName(username string) string {
	return fmt.Sprintf("dashboard-user-%s", username)
}

// UsernameFromSecretName extracts username from Secret name.
func UsernameFromSecretName(name string) (string, bool) {
	prefix := "dashboard-user-"
	if strings.HasPrefix(name, prefix) {
		return strings.TrimPrefix(name, prefix), true
	}
	return "", false
}

func parseUserSecret(secret *corev1.Secret) (*UserSecretData, error) {
	passwordHash := string(secret.Data["passwordHash"])
	if passwordHash == "" {
		return nil, fmt.Errorf("invalid user secret: missing passwordHash")
	}

	createdAtStr := string(secret.Data["createdAt"])
	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		createdAt = time.Now()
	}

	return &UserSecretData{
		PasswordHash: passwordHash,
		CreatedAt:    createdAt,
	}, nil
}

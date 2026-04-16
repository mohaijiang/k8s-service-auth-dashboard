package bootstrap

import (
	"context"
	"log"

	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/auth"
	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/k8s"
	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/validator"
	"k8s.io/client-go/kubernetes"
)

// InitializeAdmin creates the initial admin user from environment variables.
func InitializeAdmin(ctx context.Context, clientset *kubernetes.Clientset, namespace, username, password string) error {
	if password == "" {
		log.Println("INIT_ADMIN_PASSWORD not set, skipping admin initialization")
		return nil
	}

	if err := validator.ValidateUsername(username); err != nil {
		log.Printf("Invalid INIT_ADMIN_USERNAME: %v", err)
		return err
	}

	usernames, err := k8s.ListUserSecrets(ctx, clientset, namespace)
	if err != nil {
		log.Printf("Failed to list users: %v", err)
		return err
	}

	if len(usernames) > 0 {
		log.Printf("Users already exist (%d), skipping admin initialization", len(usernames))
		return nil
	}

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return err
	}

	if err := k8s.CreateUserSecret(ctx, clientset, namespace, username, passwordHash); err != nil {
		log.Printf("Failed to create admin user: %v", err)
		return err
	}

	log.Printf("Created initial admin user: %s", username)
	return nil
}

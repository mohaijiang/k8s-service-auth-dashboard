package k8s

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	// ErrHtpasswdNotFound is returned when an htpasswd secret is not found.
	ErrHtpasswdNotFound = errors.New("htpasswd secret not found")
	// ErrUserNotFound is returned when a user is not found in an htpasswd secret.
	ErrUserNotFound = errors.New("user not found in htpasswd")
)

const (
	HtpasswdTypeLabelKey   = "app.kubernetes.io/type"
	HtpasswdTypeLabelValue = "htpasswd"
)

// HtpasswdSecretData represents parsed data from an htpasswd Secret.
type HtpasswdSecretData struct {
	Name      string
	Namespace string
	Users     []string
	CreatedAt time.Time
}

// ListHtpasswdSecrets lists all htpasswd Secrets in the given namespace.
func ListHtpasswdSecrets(ctx context.Context, clientset *kubernetes.Clientset, namespace string) ([]HtpasswdSecretData, error) {
	labelSelector := fmt.Sprintf("%s=%s,%s=%s",
		UserSecretLabelKey, UserSecretLabelValue,
		HtpasswdTypeLabelKey, HtpasswdTypeLabelValue,
	)

	secrets, err := clientset.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list htpasswd secrets: %w", err)
	}

	result := make([]HtpasswdSecretData, 0, len(secrets.Items))
	for i := range secrets.Items {
		data := parseHtpasswdSecret(&secrets.Items[i])
		result = append(result, data)
	}
	return result, nil
}

// GetHtpasswdSecret retrieves a single htpasswd Secret by name.
func GetHtpasswdSecret(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string) (*HtpasswdSecretData, error) {
	secret, err := clientset.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, ErrHtpasswdNotFound
		}
		return nil, fmt.Errorf("failed to get htpasswd secret: %w", err)
	}

	data := parseHtpasswdSecret(secret)
	return &data, nil
}

// CreateHtpasswdSecret creates a new htpasswd Secret with the given users.
func CreateHtpasswdSecret(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string, users map[string]string) error {
	content := generateHtpasswdContent(users)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				UserSecretLabelKey:   UserSecretLabelValue,
				HtpasswdTypeLabelKey: HtpasswdTypeLabelValue,
			},
			Annotations: map[string]string{
				"createdAt": time.Now().UTC().Format(time.RFC3339),
			},
		},
		Data: map[string][]byte{
			".htpasswd": content,
		},
	}

	_, err := clientset.CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create htpasswd secret: %w", err)
	}
	return nil
}

// AddUserToHtpasswd adds or updates a user in an existing htpasswd Secret.
// Retries up to 3 times on conflict (409) to handle concurrent writes.
func AddUserToHtpasswd(ctx context.Context, clientset *kubernetes.Clientset, namespace, name, username, password string) error {
	const maxRetries = 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		secret, err := clientset.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			if k8serrors.IsNotFound(err) {
				return ErrHtpasswdNotFound
			}
			return fmt.Errorf("failed to get htpasswd secret: %w", err)
		}

		// Read existing htpasswd lines to preserve hashes
		userMap := make(map[string]string)
		lines := strings.Split(strings.TrimSpace(string(secret.Data[".htpasswd"])), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				userMap[parts[0]] = parts[1]
			}
		}

		userMap[username] = fmt.Sprintf("{SHA}%s", sha1Base64(password))
		secret.Data[".htpasswd"] = generateHtpasswdContentFromMap(userMap)

		_, err = clientset.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
		if err != nil {
			if k8serrors.IsConflict(err) && attempt < maxRetries-1 {
				continue // retry on conflict
			}
			return fmt.Errorf("failed to update htpasswd secret: %w", err)
		}
		return nil
	}
	return fmt.Errorf("failed to update htpasswd secret after %d retries", maxRetries)
}

// RemoveUserFromHtpasswd removes a user from an existing htpasswd Secret.
// Retries up to 3 times on conflict (409) to handle concurrent writes.
func RemoveUserFromHtpasswd(ctx context.Context, clientset *kubernetes.Clientset, namespace, name, username string) error {
	const maxRetries = 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		secret, err := clientset.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			if k8serrors.IsNotFound(err) {
				return ErrHtpasswdNotFound
			}
			return fmt.Errorf("failed to get htpasswd secret: %w", err)
		}

		existingUsers := parseHtpasswdContent(secret.Data[".htpasswd"])
		found := false
		for _, u := range existingUsers {
			if u == username {
				found = true
				break
			}
		}
		if !found {
			return ErrUserNotFound
		}

		// Build user map without the deleted user
		userMap := make(map[string]string)
		lines := strings.Split(strings.TrimSpace(string(secret.Data[".htpasswd"])), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 && parts[0] != username {
				userMap[parts[0]] = parts[1]
			}
		}

		secret.Data[".htpasswd"] = generateHtpasswdContentFromMap(userMap)

		_, err = clientset.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
		if err != nil {
			if k8serrors.IsConflict(err) && attempt < maxRetries-1 {
				continue // retry on conflict
			}
			return fmt.Errorf("failed to update htpasswd secret: %w", err)
		}
		return nil
	}
	return fmt.Errorf("failed to update htpasswd secret after %d retries", maxRetries)
}

// DeleteHtpasswdSecret deletes an htpasswd Secret.
func DeleteHtpasswdSecret(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string) error {
	err := clientset.CoreV1().Secrets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ErrHtpasswdNotFound
		}
		return fmt.Errorf("failed to delete htpasswd secret: %w", err)
	}
	return nil
}

// generateHtpasswdContent generates htpasswd file content from a username->password map.
func generateHtpasswdContent(users map[string]string) []byte {
	if len(users) == 0 {
		return nil
	}

	// Sort usernames for deterministic output
	usernames := make([]string, 0, len(users))
	for u := range users {
		usernames = append(usernames, u)
	}
	sort.Strings(usernames)

	lines := make([]string, 0, len(usernames))
	for _, u := range usernames {
		hash := sha1Base64(users[u])
		lines = append(lines, fmt.Sprintf("%s:{SHA}%s", u, hash))
	}

	return []byte(strings.Join(lines, "\n") + "\n")
}

// generateHtpasswdContentFromMap generates htpasswd content from a username->hash map.
func generateHtpasswdContentFromMap(users map[string]string) []byte {
	if len(users) == 0 {
		return nil
	}

	usernames := make([]string, 0, len(users))
	for u := range users {
		usernames = append(usernames, u)
	}
	sort.Strings(usernames)

	lines := make([]string, 0, len(usernames))
	for _, u := range usernames {
		lines = append(lines, fmt.Sprintf("%s:%s", u, users[u]))
	}

	return []byte(strings.Join(lines, "\n") + "\n")
}

// parseHtpasswdContent parses htpasswd file content and returns a list of usernames.
func parseHtpasswdContent(data []byte) []string {
	if len(data) == 0 {
		return nil
	}

	var users []string
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 1 {
			continue
		}
		users = append(users, parts[0])
	}
	return users
}

// parseHtpasswdSecret extracts data from a K8s Secret into HtpasswdSecretData.
func parseHtpasswdSecret(secret *corev1.Secret) HtpasswdSecretData {
	data := HtpasswdSecretData{
		Name:      secret.Name,
		Namespace: secret.Namespace,
		Users:     parseHtpasswdContent(secret.Data[".htpasswd"]),
	}

	createdAtStr := secret.Annotations["createdAt"]
	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		data.CreatedAt = t
	}

	return data
}

// sha1Base64 computes the SHA1 hash of a string and returns it as base64.
func sha1Base64(s string) string {
	h := sha1.Sum([]byte(s))
	return base64.StdEncoding.EncodeToString(h[:])
}

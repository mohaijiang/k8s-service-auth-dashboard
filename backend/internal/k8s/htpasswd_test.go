package k8s

import (
	"strings"
	"testing"
)

func TestGenerateHtpasswdContent(t *testing.T) {
	tests := []struct {
		name           string
		users          map[string]string // username -> password
		expectedLines  int
		expectedPrefix string
	}{
		{
			name:           "single user",
			users:          map[string]string{"admin": "password123"},
			expectedLines:  1,
			expectedPrefix: "admin:{SHA}",
		},
		{
			name: "multiple users",
			users: map[string]string{
				"admin": "pass1",
				"user1": "pass2",
				"user2": "pass3",
			},
			expectedLines:  3,
			expectedPrefix: "", // multiple users, no single prefix
		},
		{
			name:           "empty users",
			users:          map[string]string{},
			expectedLines:  0,
			expectedPrefix: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := generateHtpasswdContent(tt.users)

			if tt.expectedLines == 0 {
				if len(content) != 0 {
					t.Errorf("expected empty content, got %q", string(content))
				}
				return
			}

			lines := strings.Split(strings.TrimSpace(string(content)), "\n")
			if len(lines) != tt.expectedLines {
				t.Errorf("expected %d lines, got %d", tt.expectedLines, len(lines))
			}

			if tt.expectedPrefix != "" {
				if !strings.HasPrefix(lines[0], tt.expectedPrefix) {
					t.Errorf("expected line to start with %q, got %q", tt.expectedPrefix, lines[0])
				}
			}

			// Verify all lines have {SHA} format
			for _, line := range lines {
				if !strings.Contains(line, ":{SHA}") {
					t.Errorf("line %q does not contain :{{SHA}}", line)
				}
			}
		})
	}
}

func TestParseHtpasswdContent(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		expectedUsers []string
	}{
		{
			name:          "single user",
			content:       "admin:{SHA}0DPiKuNIrrVmD8IUCuw1hQxNqZc=\n",
			expectedUsers: []string{"admin"},
		},
		{
			name:          "multiple users",
			content:       "admin:{SHA}0DPiKuNIrrVmD8IUCuw1hQxNqZc=\nuser1:{SHA}xyz123abc\n",
			expectedUsers: []string{"admin", "user1"},
		},
		{
			name:          "empty content",
			content:       "",
			expectedUsers: []string{},
		},
		{
			name:          "whitespace only",
			content:       "  \n  \n",
			expectedUsers: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := parseHtpasswdContent([]byte(tt.content))

			if len(users) != len(tt.expectedUsers) {
				t.Errorf("expected %d users, got %d", len(tt.expectedUsers), len(users))
				return
			}

			for i, expected := range tt.expectedUsers {
				if users[i] != expected {
					t.Errorf("user[%d] = %q, want %q", i, users[i], expected)
				}
			}
		})
	}
}

func TestGenerateAndParseRoundTrip(t *testing.T) {
	users := map[string]string{
		"admin": "password123",
		"user1": "pass456",
	}

	content := generateHtpasswdContent(users)
	parsed := parseHtpasswdContent(content)

	if len(parsed) != len(users) {
		t.Fatalf("round trip: expected %d users, got %d", len(users), len(parsed))
	}

	parsedSet := make(map[string]bool)
	for _, u := range parsed {
		parsedSet[u] = true
	}

	for username := range users {
		if !parsedSet[username] {
			t.Errorf("round trip: user %q not found in parsed output", username)
		}
	}
}

func TestSHA1HashDeterministic(t *testing.T) {
	users := map[string]string{"admin": "password123"}

	content1 := generateHtpasswdContent(users)
	content2 := generateHtpasswdContent(users)

	if string(content1) != string(content2) {
		t.Errorf("SHA1 hash should be deterministic:\n  got1: %s\n  got2: %s", content1, content2)
	}
}

package validator

import "testing"

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		wantErr     bool
		errContains string
	}{
		{"valid simple", "admin", false, ""},
		{"valid with hyphens", "my-user-123", false, ""},
		{"valid min length 3", "abc", false, ""},
		{"valid max length 32", "abcdefghijklmnopqrstuvwxyz123456", false, ""},
		{"too short", "ab", true, "at least 3 characters"},
		{"too long", "abcdefghijklmnopqrstuvwxyz1234567", true, "at most 32 characters"},
		{"uppercase letters", "Admin", true, "lowercase letters"},
		{"starts with hyphen", "-admin", true, "start and end with alphanumeric"},
		{"ends with hyphen", "admin-", true, "start and end with alphanumeric"},
		{"contains spaces", "my user", true, "lowercase letters"},
		{"contains underscores", "my_user", true, "lowercase letters"},
		{"empty string", "", true, "at least 3 characters"},
		{"single char", "a", true, "at least 3 characters"},
		{"numbers only", "123", false, ""},
		{"consecutive hyphens", "my--user", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername(%q) error = %v, wantErr %v", tt.username, err, tt.wantErr)
			}
			if err != nil && tt.errContains != "" {
				if !containsString(err.Error(), tt.errContains) {
					t.Errorf("ValidateUsername(%q) error = %q, want to contain %q", tt.username, err.Error(), tt.errContains)
				}
			}
		})
	}
}

func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

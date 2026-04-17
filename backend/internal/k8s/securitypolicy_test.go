package k8s

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestParseSecurityPolicy(t *testing.T) {
	tests := []struct {
		name                string
		obj                 *unstructured.Unstructured
		wantName            string
		wantNamespace       string
		wantBasicAuth       bool
		wantTLS             bool
		wantTargetName      string
		wantTargetKind      string
		wantBasicAuthSecret string
	}{
		{
			name: "policy with basic auth",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name":      "my-app-auth",
						"namespace": "production",
					},
					"spec": map[string]interface{}{
						"targetRef": map[string]interface{}{
							"name":      "my-app-route",
							"namespace": "production",
							"kind":      "HTTPRoute",
						},
						"basicAuth": map[string]interface{}{
							"users": map[string]interface{}{
								"name": "htpasswd-secret",
							},
						},
					},
				},
			},
			wantName:            "my-app-auth",
			wantNamespace:       "production",
			wantBasicAuth:       true,
			wantTLS:             false,
			wantTargetName:      "my-app-route",
			wantTargetKind:      "HTTPRoute",
			wantBasicAuthSecret: "htpasswd-secret",
		},
		{
			name: "policy with tls",
			wantBasicAuthSecret: "",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name":      "my-app-tls",
						"namespace": "staging",
					},
					"spec": map[string]interface{}{
						"targetRef": map[string]interface{}{
							"name": "my-app-route",
							"kind": "HTTPRoute",
						},
						"tls": map[string]interface{}{
							"certificateRef": map[string]interface{}{
								"name": "my-cert",
							},
						},
					},
				},
			},
			wantName:       "my-app-tls",
			wantNamespace:  "staging",
			wantBasicAuth:  false,
			wantTLS:        true,
			wantTargetName: "my-app-route",
			wantTargetKind: "HTTPRoute",
		},
		{
			name: "policy with both basic auth and tls",
			wantBasicAuthSecret: "",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name":      "full-policy",
						"namespace": "default",
					},
					"spec": map[string]interface{}{
						"targetRef": map[string]interface{}{
							"name": "some-route",
							"kind": "HTTPRoute",
						},
						"basicAuth": map[string]interface{}{},
						"tls":       map[string]interface{}{},
					},
				},
			},
			wantName:       "full-policy",
			wantNamespace:  "default",
			wantBasicAuth:  true,
			wantTLS:        true,
			wantTargetName: "some-route",
			wantTargetKind: "HTTPRoute",
		},
		{
			name: "policy with minimal fields",
			wantBasicAuthSecret: "",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name":      "minimal-policy",
						"namespace": "test",
					},
					"spec": map[string]interface{}{
						"targetRef": map[string]interface{}{
							"name": "route",
							"kind": "HTTPRoute",
						},
					},
				},
			},
			wantName:       "minimal-policy",
			wantNamespace:  "test",
			wantBasicAuth:  false,
			wantTLS:        false,
			wantTargetName: "route",
			wantTargetKind: "HTTPRoute",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseSecurityPolicy(tt.obj)

			if result.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", result.Name, tt.wantName)
			}
			if result.Namespace != tt.wantNamespace {
				t.Errorf("Namespace = %q, want %q", result.Namespace, tt.wantNamespace)
			}
			if result.HasBasicAuth != tt.wantBasicAuth {
				t.Errorf("HasBasicAuth = %v, want %v", result.HasBasicAuth, tt.wantBasicAuth)
			}
			if result.HasTLS != tt.wantTLS {
				t.Errorf("HasTLS = %v, want %v", result.HasTLS, tt.wantTLS)
			}
			if result.TargetRef.Name != tt.wantTargetName {
				t.Errorf("TargetRef.Name = %q, want %q", result.TargetRef.Name, tt.wantTargetName)
			}
			if result.TargetRef.Kind != tt.wantTargetKind {
				t.Errorf("TargetRef.Kind = %q, want %q", result.TargetRef.Kind, tt.wantTargetKind)
			}
			if result.BasicAuthSecretName != tt.wantBasicAuthSecret {
				t.Errorf("BasicAuthSecretName = %q, want %q", result.BasicAuthSecretName, tt.wantBasicAuthSecret)
			}
		})
	}
}

package k8s

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestParseHTTPRoute(t *testing.T) {
	tests := []struct {
		name            string
		obj             *unstructured.Unstructured
		wantName        string
		wantNamespace   string
		wantHostnames   int
		wantBackendRefs int
		wantParentRefs  int
	}{
		{
			name: "full httproute with all fields",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name":      "my-app-route",
						"namespace": "production",
					},
					"spec": map[string]interface{}{
						"hostnames": []interface{}{"app.example.com", "app.example.org"},
						"rules": []interface{}{
							map[string]interface{}{
								"backendRefs": []interface{}{
									map[string]interface{}{
										"name":      "my-app",
										"namespace": "production",
									},
								},
							},
						},
						"parentRefs": []interface{}{
							map[string]interface{}{
								"name":      "envoy-gateway",
								"namespace": "gateway",
							},
						},
					},
				},
			},
			wantName:        "my-app-route",
			wantNamespace:   "production",
			wantHostnames:   2,
			wantBackendRefs: 1,
			wantParentRefs:  1,
		},
		{
			name: "httproute with minimal fields",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name":      "minimal-route",
						"namespace": "default",
					},
					"spec": map[string]interface{}{},
				},
			},
			wantName:        "minimal-route",
			wantNamespace:   "default",
			wantHostnames:   0,
			wantBackendRefs: 0,
			wantParentRefs:  0,
		},
		{
			name: "httproute with backendRef without namespace",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name":      "svc-route",
						"namespace": "staging",
					},
					"spec": map[string]interface{}{
						"rules": []interface{}{
							map[string]interface{}{
								"backendRefs": []interface{}{
									map[string]interface{}{
										"name": "svc",
									},
								},
							},
						},
					},
				},
			},
			wantName:        "svc-route",
			wantNamespace:   "staging",
			wantHostnames:   0,
			wantBackendRefs: 1,
			wantParentRefs:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseHTTPRoute(tt.obj)

			if result.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", result.Name, tt.wantName)
			}
			if result.Namespace != tt.wantNamespace {
				t.Errorf("Namespace = %q, want %q", result.Namespace, tt.wantNamespace)
			}
			if len(result.Hostnames) != tt.wantHostnames {
				t.Errorf("len(Hostnames) = %d, want %d", len(result.Hostnames), tt.wantHostnames)
			}
			if len(result.BackendRefs) != tt.wantBackendRefs {
				t.Errorf("len(BackendRefs) = %d, want %d", len(result.BackendRefs), tt.wantBackendRefs)
			}
			if len(result.ParentRefs) != tt.wantParentRefs {
				t.Errorf("len(ParentRefs) = %d, want %d", len(result.ParentRefs), tt.wantParentRefs)
			}
		})
	}
}

func TestParseHTTPRouteFieldValues(t *testing.T) {
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name":      "my-route",
				"namespace": "prod",
			},
			"spec": map[string]interface{}{
				"hostnames": []interface{}{"app.example.com"},
				"rules": []interface{}{
					map[string]interface{}{
						"backendRefs": []interface{}{
							map[string]interface{}{"name": "svc-a", "namespace": "prod"},
							map[string]interface{}{"name": "svc-b"},
						},
					},
				},
				"parentRefs": []interface{}{
					map[string]interface{}{"name": "my-gateway", "namespace": "gw"},
				},
			},
		},
	}

	result := parseHTTPRoute(obj)

	if result.Hostnames[0] != "app.example.com" {
		t.Errorf("Hostname[0] = %q, want %q", result.Hostnames[0], "app.example.com")
	}
	if result.BackendRefs[0].Name != "svc-a" {
		t.Errorf("BackendRefs[0].Name = %q, want %q", result.BackendRefs[0].Name, "svc-a")
	}
	if result.BackendRefs[0].Namespace != "prod" {
		t.Errorf("BackendRefs[0].Namespace = %q, want %q", result.BackendRefs[0].Namespace, "prod")
	}
	if result.BackendRefs[1].Name != "svc-b" {
		t.Errorf("BackendRefs[1].Name = %q, want %q", result.BackendRefs[1].Name, "svc-b")
	}
	if result.BackendRefs[1].Namespace != "" {
		t.Errorf("BackendRefs[1].Namespace = %q, want empty", result.BackendRefs[1].Namespace)
	}
	if result.ParentRefs[0].Name != "my-gateway" {
		t.Errorf("ParentRefs[0].Name = %q, want %q", result.ParentRefs[0].Name, "my-gateway")
	}
}

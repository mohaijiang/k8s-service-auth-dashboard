package handler

import (
	"testing"

	corev1 "k8s.io/api/core/v1"

	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/k8s"
)

func TestBuildServiceOverviews(t *testing.T) {
	tests := []struct {
		name           string
		services       []k8s.ServiceInfo
		routes         []k8s.HTTPRouteData
		policies       []k8s.SecurityPolicyData
		wantLen        int
		wantRouteName  map[string]string // key: "ns/svc" -> route name
		wantPolicyName map[string]string // key: "ns/svc" -> policy name
		wantNoRoute    []string          // key: "ns/svc" -> expect no route
		wantNoPolicy   []string          // key: "ns/svc" -> expect no policy
	}{
		{
			name: "service with route and policy",
			services: []k8s.ServiceInfo{
				{Name: "my-app", Namespace: "prod", ClusterIP: "10.0.0.1"},
			},
			routes: []k8s.HTTPRouteData{
				{
					Name: "my-app-route", Namespace: "prod",
					Hostnames:   []string{"app.example.com"},
					BackendRefs: []k8s.BackendRef{{Name: "my-app", Namespace: "prod"}},
					ParentRefs:  []k8s.ParentRefData{{Name: "envoy-gateway", Namespace: "gateway"}},
				},
			},
			policies: []k8s.SecurityPolicyData{
				{
					Name: "my-app-auth", Namespace: "prod",
					TargetRef:    k8s.TargetRefData{Name: "my-app-route", Kind: "HTTPRoute"},
					HasBasicAuth: true, HasTLS: false,
				},
			},
			wantLen:        1,
			wantRouteName:  map[string]string{"prod/my-app": "my-app-route"},
			wantPolicyName: map[string]string{"prod/my-app": "my-app-auth"},
		},
		{
			name: "service with no route",
			services: []k8s.ServiceInfo{
				{Name: "orphan-svc", Namespace: "default", ClusterIP: "10.0.0.2"},
			},
			routes:       []k8s.HTTPRouteData{},
			policies:     []k8s.SecurityPolicyData{},
			wantLen:      1,
			wantNoRoute:  []string{"default/orphan-svc"},
			wantNoPolicy: []string{"default/orphan-svc"},
		},
		{
			name: "route backendRef without namespace defaults to route namespace",
			services: []k8s.ServiceInfo{
				{Name: "svc", Namespace: "staging", ClusterIP: "10.0.0.3"},
			},
			routes: []k8s.HTTPRouteData{
				{
					Name: "svc-route", Namespace: "staging",
					BackendRefs: []k8s.BackendRef{{Name: "svc"}},
				},
			},
			policies:      []k8s.SecurityPolicyData{},
			wantLen:       1,
			wantRouteName: map[string]string{"staging/svc": "svc-route"},
			wantNoPolicy:  []string{"staging/svc"},
		},
		{
			name: "route found but no matching policy",
			services: []k8s.ServiceInfo{
				{Name: "svc-a", Namespace: "default", ClusterIP: "10.0.0.4"},
			},
			routes: []k8s.HTTPRouteData{
				{
					Name: "svc-a-route", Namespace: "default",
					BackendRefs: []k8s.BackendRef{{Name: "svc-a", Namespace: "default"}},
				},
			},
			policies:      []k8s.SecurityPolicyData{},
			wantLen:       1,
			wantRouteName: map[string]string{"default/svc-a": "svc-a-route"},
			wantNoPolicy:  []string{"default/svc-a"},
		},
		{
			name:     "empty lists",
			services: []k8s.ServiceInfo{},
			routes:   []k8s.HTTPRouteData{},
			policies: []k8s.SecurityPolicyData{},
			wantLen:  0,
		},
		{
			name: "multiple services with mixed associations",
			services: []k8s.ServiceInfo{
				{Name: "app-a", Namespace: "prod", ClusterIP: "10.0.0.1"},
				{Name: "app-b", Namespace: "prod", ClusterIP: "10.0.0.2"},
				{Name: "app-c", Namespace: "staging", ClusterIP: "10.0.0.3"},
			},
			routes: []k8s.HTTPRouteData{
				{
					Name: "app-a-route", Namespace: "prod",
					BackendRefs: []k8s.BackendRef{{Name: "app-a", Namespace: "prod"}},
				},
			},
			policies: []k8s.SecurityPolicyData{
				{
					Name: "app-a-auth", Namespace: "prod",
					TargetRef:    k8s.TargetRefData{Name: "app-a-route", Kind: "HTTPRoute"},
					HasBasicAuth: true,
				},
			},
			wantLen:        3,
			wantRouteName:  map[string]string{"prod/app-a": "app-a-route"},
			wantPolicyName: map[string]string{"prod/app-a": "app-a-auth"},
			wantNoRoute:    []string{"prod/app-b", "staging/app-c"},
			wantNoPolicy:   []string{"prod/app-b", "staging/app-c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildServiceOverviews(tt.services, tt.routes, tt.policies)

			if len(result) != tt.wantLen {
				t.Fatalf("got %d overviews, want %d", len(result), tt.wantLen)
			}

			// Check route and policy associations
			for _, svc := range tt.services {
				key := svc.Namespace + "/" + svc.Name
				for _, ov := range result {
					if ov.Name != svc.Name || ov.Namespace != svc.Namespace {
						continue
					}
					ovKey := ov.Namespace + "/" + ov.Name

					if wantName, ok := tt.wantRouteName[ovKey]; ok {
						if ov.HTTPRoute == nil {
							t.Errorf("svc %q: expected route %q, got nil", key, wantName)
						} else if ov.HTTPRoute.Name != wantName {
							t.Errorf("svc %q: route name = %q, want %q", key, ov.HTTPRoute.Name, wantName)
						}
					}

					if wantName, ok := tt.wantPolicyName[ovKey]; ok {
						if ov.SecurityPolicy == nil {
							t.Errorf("svc %q: expected policy %q, got nil", key, wantName)
						} else if ov.SecurityPolicy.Name != wantName {
							t.Errorf("svc %q: policy name = %q, want %q", key, ov.SecurityPolicy.Name, wantName)
						}
					}
				}
			}

			// Check no-route expectations
			for _, key := range tt.wantNoRoute {
				for _, ov := range result {
					ovKey := ov.Namespace + "/" + ov.Name
					if ovKey == key && ov.HTTPRoute != nil {
						t.Errorf("svc %q: expected no route, got %q", key, ov.HTTPRoute.Name)
					}
				}
			}

			// Check no-policy expectations
			for _, key := range tt.wantNoPolicy {
				for _, ov := range result {
					ovKey := ov.Namespace + "/" + ov.Name
					if ovKey == key && ov.SecurityPolicy != nil {
						t.Errorf("svc %q: expected no policy, got %q", key, ov.SecurityPolicy.Name)
					}
				}
			}
		})
	}
}

func TestBuildServiceOverviewsPorts(t *testing.T) {
	services := []k8s.ServiceInfo{
		{
			Name: "svc", Namespace: "default", ClusterIP: "10.0.0.1",
			Ports: []corev1.ServicePort{
				{Name: "http", Port: 80, Protocol: corev1.ProtocolTCP},
			},
		},
	}

	result := BuildServiceOverviews(services, nil, nil)
	if len(result) != 1 {
		t.Fatalf("got %d overviews, want 1", len(result))
	}
	if len(result[0].Ports) != 1 {
		t.Errorf("got %d ports, want 1", len(result[0].Ports))
	}
	if result[0].Ports[0].Port != 80 {
		t.Errorf("Port = %d, want 80", result[0].Ports[0].Port)
	}
	if result[0].Ports[0].Protocol != "TCP" {
		t.Errorf("Protocol = %q, want %q", result[0].Ports[0].Protocol, "TCP")
	}
}

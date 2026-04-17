package model

// ServiceOverview represents a K8s Service with its associated routing and security info.
type ServiceOverview struct {
	Name           string             `json:"name"`
	Namespace      string             `json:"namespace"`
	ClusterIP      string             `json:"clusterIP"`
	Ports          []ServicePort      `json:"ports"`
	Selector       map[string]string  `json:"selector"`
	HTTPRoute      *HTTPRouteInfo     `json:"httpRoute,omitempty"`
	SecurityPolicy *SecurityPolicyInfo `json:"securityPolicy,omitempty"`
}

// ServicePort represents a single port on a K8s Service.
type ServicePort struct {
	Name       string `json:"name"`
	Port       int32  `json:"port"`
	TargetPort string `json:"targetPort"`
	Protocol   string `json:"protocol"`
}

// HTTPRouteInfo is a summary of the HTTPRoute associated with a Service.
type HTTPRouteInfo struct {
	Name       string      `json:"name"`
	Namespace  string      `json:"namespace"`
	Hostnames  []string    `json:"hostnames"`
	ParentRefs []ParentRef `json:"parentRefs"`
}

// ParentRef represents a Gateway that an HTTPRoute is attached to.
type ParentRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

// SecurityPolicyInfo is a summary of the SecurityPolicy targeting an HTTPRoute.
type SecurityPolicyInfo struct {
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	HasBasicAuth bool   `json:"hasBasicAuth"`
	HasTLS       bool   `json:"hasTLS"`
}

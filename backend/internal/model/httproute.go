package model

// HTTPRouteCreateRequest defines the request body for creating an HTTPRoute.
type HTTPRouteCreateRequest struct {
	Name             string                 `json:"name" binding:"required"`
	Namespace        string                 `json:"namespace" binding:"required"`
	Hostnames        []string               `json:"hostnames" binding:"required,min=1"`
	ServiceName      string                 `json:"serviceName" binding:"required"`
	ServicePort      int32                  `json:"servicePort" binding:"required,min=1"`
	ParentRefs       []ParentRef             `json:"parentRefs" binding:"required,min=1"`
	SecurityPolicy   *SecurityPolicyConfig    `json:"securityPolicy,omitempty"`
}

// SecurityPolicyConfig defines the configuration for creating a SecurityPolicy alongside the HTTPRoute.
type SecurityPolicyConfig struct {
	BasicAuthSecretName string `json:"basicAuthSecretName" binding:"required"`
}

// HTTPRouteDetail represents the response for HTTPRoute details.
type HTTPRouteDetail struct {
	Name       string     `json:"name"`
	Namespace  string     `json:"namespace"`
	Hostnames  []string   `json:"hostnames"`
	ParentRefs []ParentRef `json:"parentRefs"`
	ServiceName string    `json:"serviceName"`
}

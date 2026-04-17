package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// SecurityPolicyData holds parsed SecurityPolicy information.
type SecurityPolicyData struct {
	Name                string
	Namespace           string
	TargetRef           TargetRefData
	HasBasicAuth        bool
	HasTLS              bool
	BasicAuthSecretName string
}

// TargetRefData represents the target reference of a SecurityPolicy.
type TargetRefData struct {
	Name      string
	Namespace string
	Kind      string
}

var securityPolicyGVR = schema.GroupVersionResource{
	Group:    "gateway.envoyproxy.io",
	Version:  "v1alpha1",
	Resource: "securitypolicies",
}

// ListSecurityPolicies lists SecurityPolicy resources, optionally filtered by namespace.
func ListSecurityPolicies(ctx context.Context, dynClient dynamic.Interface, namespace string) ([]SecurityPolicyData, error) {
	list, err := dynClient.Resource(securityPolicyGVR).Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list securitypolicies: %w", err)
	}

	result := make([]SecurityPolicyData, 0, len(list.Items))
	for i := range list.Items {
		policy := parseSecurityPolicy(&list.Items[i])
		result = append(result, policy)
	}
	return result, nil
}

// parseSecurityPolicy extracts structured data from an unstructured SecurityPolicy.
func parseSecurityPolicy(obj *unstructured.Unstructured) SecurityPolicyData {
	policy := SecurityPolicyData{
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
	}

	targetRefName, _, _ := unstructured.NestedString(obj.Object, "spec", "targetRef", "name")
	targetRefNS, _, _ := unstructured.NestedString(obj.Object, "spec", "targetRef", "namespace")
	targetRefKind, _, _ := unstructured.NestedString(obj.Object, "spec", "targetRef", "kind")
	policy.TargetRef = TargetRefData{
		Name:      targetRefName,
		Namespace: targetRefNS,
		Kind:      targetRefKind,
	}

	_, ok, _ := unstructured.NestedMap(obj.Object, "spec", "basicAuth")
	policy.HasBasicAuth = ok

	if ok {
		secretName, _, _ := unstructured.NestedString(obj.Object, "spec", "basicAuth", "users", "name")
		policy.BasicAuthSecretName = secretName
	}

	_, hasTLS, _ := unstructured.NestedMap(obj.Object, "spec", "tls")
	policy.HasTLS = hasTLS

	return policy
}

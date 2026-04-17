package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// HTTPRouteData holds parsed HTTPRoute information.
type HTTPRouteData struct {
	Name        string
	Namespace   string
	Hostnames   []string
	BackendRefs []BackendRef
	ParentRefs  []ParentRefData
}

// BackendRef represents a backend reference in an HTTPRoute rule.
type BackendRef struct {
	Name      string
	Namespace string
}

// ParentRefData represents a parent gateway reference.
type ParentRefData struct {
	Name      string
	Namespace string
}

var httpRouteGVR = schema.GroupVersionResource{
	Group:    "gateway.networking.k8s.io",
	Version:  "v1",
	Resource: "httproutes",
}

// ListHTTPRoutes lists HTTPRoute resources, optionally filtered by namespace.
func ListHTTPRoutes(ctx context.Context, dynClient dynamic.Interface, namespace string) ([]HTTPRouteData, error) {
	list, err := dynClient.Resource(httpRouteGVR).Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list httproutes: %w", err)
	}

	result := make([]HTTPRouteData, 0, len(list.Items))
	for i := range list.Items {
		route := parseHTTPRoute(&list.Items[i])
		result = append(result, route)
	}
	return result, nil
}

// parseHTTPRoute extracts structured data from an unstructured HTTPRoute.
func parseHTTPRoute(obj *unstructured.Unstructured) HTTPRouteData {
	route := HTTPRouteData{
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
	}

	if hostnames, ok, _ := unstructured.NestedStringSlice(obj.Object, "spec", "hostnames"); ok {
		route.Hostnames = hostnames
	}

	rules, ok, _ := unstructured.NestedSlice(obj.Object, "spec", "rules")
	if ok {
		for _, rule := range rules {
			ruleMap, ok := rule.(map[string]interface{})
			if !ok {
				continue
			}
			backendRefs, ok, _ := unstructured.NestedSlice(ruleMap, "backendRefs")
			if ok {
				for _, ref := range backendRefs {
					refMap, ok := ref.(map[string]interface{})
					if !ok {
						continue
					}
					name, _, _ := unstructured.NestedString(refMap, "name")
					ns, _, _ := unstructured.NestedString(refMap, "namespace")
					route.BackendRefs = append(route.BackendRefs, BackendRef{
						Name:      name,
						Namespace: ns,
					})
				}
			}
		}
	}

	parentRefs, ok, _ := unstructured.NestedSlice(obj.Object, "spec", "parentRefs")
	if ok {
		for _, ref := range parentRefs {
			refMap, ok := ref.(map[string]interface{})
			if !ok {
				continue
			}
			name, _, _ := unstructured.NestedString(refMap, "name")
			ns, _, _ := unstructured.NestedString(refMap, "namespace")
			route.ParentRefs = append(route.ParentRefs, ParentRefData{
				Name:      name,
				Namespace: ns,
			})
		}
	}

	return route
}

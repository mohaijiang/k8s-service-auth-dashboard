package handler

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/k8s"
	"github.com/mohaijiang/k8s-service-auth-dashboard/backend/internal/model"
)

// BuildServiceOverviews builds the Service -> HTTPRoute -> SecurityPolicy association chain.
// This is a pure function with no external dependencies, making it easy to unit test.
func BuildServiceOverviews(
	services []k8s.ServiceInfo,
	httpRoutes []k8s.HTTPRouteData,
	policies []k8s.SecurityPolicyData,
) []model.ServiceOverview {
	// Index HTTPRoutes by (namespace, backendRefName)
	routeByBackend := make(map[string]*k8s.HTTPRouteData)
	for i := range httpRoutes {
		route := &httpRoutes[i]
		for _, backend := range route.BackendRefs {
			ns := backend.Namespace
			if ns == "" {
				ns = route.Namespace
			}
			key := ns + "/" + backend.Name
			routeByBackend[key] = route
		}
	}

	// Index SecurityPolicies by (namespace, targetRefName)
	policyByTarget := make(map[string]*k8s.SecurityPolicyData)
	for i := range policies {
		policy := &policies[i]
		ns := policy.TargetRef.Namespace
		if ns == "" {
			ns = policy.Namespace
		}
		key := ns + "/" + policy.TargetRef.Name
		policyByTarget[key] = policy
	}

	result := make([]model.ServiceOverview, 0, len(services))
	for _, svc := range services {
		overview := model.ServiceOverview{
			Name:      svc.Name,
			Namespace: svc.Namespace,
			ClusterIP: svc.ClusterIP,
			Ports:     convertPorts(svc.Ports),
			Selector:  svc.Selector,
		}

		// Find associated HTTPRoute
		svcKey := svc.Namespace + "/" + svc.Name
		if route, ok := routeByBackend[svcKey]; ok {
			overview.HTTPRoute = &model.HTTPRouteInfo{
				Name:       route.Name,
				Namespace:  route.Namespace,
				Hostnames:  route.Hostnames,
				ParentRefs: convertParentRefs(route.ParentRefs),
			}

			// Find SecurityPolicy targeting this HTTPRoute
			routeKey := route.Namespace + "/" + route.Name
			if policy, ok := policyByTarget[routeKey]; ok {
				overview.SecurityPolicy = &model.SecurityPolicyInfo{
					Name:         policy.Name,
					Namespace:    policy.Namespace,
					HasBasicAuth: policy.HasBasicAuth,
					HasTLS:       policy.HasTLS,
				}
			}
		}

		result = append(result, overview)
	}

	return result
}

func convertPorts(ports []corev1.ServicePort) []model.ServicePort {
	result := make([]model.ServicePort, 0, len(ports))
	for _, p := range ports {
		result = append(result, model.ServicePort{
			Name:       p.Name,
			Port:       p.Port,
			TargetPort: p.TargetPort.String(),
			Protocol:   string(p.Protocol),
		})
	}
	return result
}

func convertParentRefs(refs []k8s.ParentRefData) []model.ParentRef {
	result := make([]model.ParentRef, 0, len(refs))
	for _, r := range refs {
		result = append(result, model.ParentRef{
			Name:      r.Name,
			Namespace: r.Namespace,
		})
	}
	return result
}

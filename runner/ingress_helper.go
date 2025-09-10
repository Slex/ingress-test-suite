package runner

import (
	"fmt"
	"ingress-test-suite/pkg/messages"
	"ingress-test-suite/test_load"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ConvertPathType(s string) (networkingv1.PathType, error) {
	switch s {
	case "Exact":
		return networkingv1.PathTypeExact, nil
	case "Prefix":
		return networkingv1.PathTypePrefix, nil
	case "ImplementationSpecific":
		return networkingv1.PathTypeImplementationSpecific, nil
	default:
		return "", fmt.Errorf(messages.InvalidPathType, s)
	}
}

func createIngressRule(t test_load.IngressTestEntry, c test_load.IngressTestsFile) *networkingv1.Ingress {
	pathType, err := ConvertPathType(t.PathType)
	if err != nil {
		log.Fatalf(messages.FailedConvertPathType, err)
	}

	return &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("test-%s", t.Host),
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: &c.IngressClassName,
			Rules: []networkingv1.IngressRule{
				{
					Host: t.Host,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     t.Path,
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: t.Service,
											Port: networkingv1.ServiceBackendPort{
												Number: int32(t.Port),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

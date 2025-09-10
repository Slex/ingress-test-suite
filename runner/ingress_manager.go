package runner

import (
	"context"
	"fmt"
	"ingress-test-suite/pkg/messages"
	"ingress-test-suite/test_load"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type K8sIngressManager struct {
	client *kubernetes.Clientset
}

func NewK8sIngressManager(client *kubernetes.Clientset) *K8sIngressManager {
	return &K8sIngressManager{client: client}
}

func (m *K8sIngressManager) Create(entry test_load.IngressTestEntry, c test_load.IngressTestsFile) error {
	ingress := createIngressRule(entry, c)
	_, err := m.client.NetworkingV1().Ingresses(entry.Namespace).
		Create(context.TODO(), ingress, metav1.CreateOptions{})
	return err
}

func (m *K8sIngressManager) Delete(entry test_load.IngressTestEntry) error {
	return m.client.NetworkingV1().Ingresses(entry.Namespace).
		Delete(context.TODO(), fmt.Sprintf("test-%s", entry.Host), metav1.DeleteOptions{})
}

func (m *K8sIngressManager) CheckExist(entry test_load.IngressTestEntry) (bool, error) {
	ingressName := fmt.Sprintf("test-%s", entry.Host)

	_, err := m.client.NetworkingV1().Ingresses(entry.Namespace).Get(
		context.TODO(),
		ingressName,
		metav1.GetOptions{},
	)

	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf(messages.FailedCheckExistIngressRule, ingressName, err)
	}

	return true, nil
}

package updater

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func DeleteServices(namespace string, selector labels.Selector, clientset *kubernetes.Clientset) error {
	services, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		return err
	}

	var lastErr error
	deletePolicy := metav1.DeletePropagationForeground
	for _, svc := range services.Items {
		err := clientset.CoreV1().Services(namespace).Delete(
			context.TODO(),
			svc.Name,
			metav1.DeleteOptions{
				PropagationPolicy: &deletePolicy,
			},
		)
		if err != nil {
		}
		lastErr = err
	}

	return lastErr
}

func CreateOrPatchService(s *corev1.Service, cl client.Client) error {
	return createOrPatchObject(s, cl)
}

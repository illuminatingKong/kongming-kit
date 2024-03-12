package updater

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func DeleteNamespace(name string, clientset *kubernetes.Clientset) error {
	deletePolicy := metav1.DeletePropagationForeground
	return clientset.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

func CreateNamespace(ns *corev1.Namespace, cl client.Client) error {
	return createObject(ns, cl)
}

func CreateNamespaceByName(ns string, labels map[string]string, cl client.Client) error {
	n := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "",
			Name:      ns,
			Labels:    labels,
		},
	}
	return CreateNamespace(n, cl)
}

func CreateOrPatchNamespace(n *corev1.Namespace, cl client.Client) error {
	return createOrPatchObject(n, cl)
}

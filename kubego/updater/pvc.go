package updater

import (
	"context"
	"github.com/illuminatingKong/kongming-kit/kubego/util"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func DeletePersistentVolumeClaims(namespace string, selector labels.Selector, clientset *kubernetes.Clientset) error {
	deletePolicy := metav1.DeletePropagationForeground
	err := clientset.CoreV1().PersistentVolumeClaims(namespace).DeleteCollection(
		context.TODO(),
		metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		},
		metav1.ListOptions{
			LabelSelector: selector.String(),
		},
	)

	return util.IgnoreNotFoundError(err)
}

func DeletePvcWithName(namespace, name string, clientset *kubernetes.Clientset) error {
	deletePolicy := metav1.DeletePropagationForeground
	return clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(
		context.TODO(), name,
		metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		},
	)
}

func CreatePvc(namespace string, pvc *corev1.PersistentVolumeClaim, clientset *kubernetes.Clientset) error {
	_, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Create(context.TODO(), pvc,
		metav1.CreateOptions{})
	return err
}

func UpdatePvc(namespace string, pvc *corev1.PersistentVolumeClaim, clientset *kubernetes.Clientset) error {
	_, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Update(context.TODO(), pvc,
		metav1.UpdateOptions{})
	return err
}

package updater

import (
	"context"
	"github.com/illuminatingKong/kongming-kit/kubego/cluster"
	"github.com/illuminatingKong/kongming-kit/kubego/util"

	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func DeleteIngresses(namespace string, selector labels.Selector, clientset *kubernetes.Clientset) error {
	version, err := clientset.Discovery().ServerVersion()
	if err != nil {
		return err
	}

	deletePolicy := metav1.DeletePropagationForeground

	if cluster.VersionLessThan122(version) {
		err = clientset.ExtensionsV1beta1().Ingresses(namespace).DeleteCollection(
			context.TODO(),
			metav1.DeleteOptions{
				PropagationPolicy: &deletePolicy,
			}, metav1.ListOptions{
				LabelSelector: selector.String(),
			})
	} else {
		err = clientset.NetworkingV1().Ingresses(namespace).DeleteCollection(
			context.TODO(),
			metav1.DeleteOptions{
				PropagationPolicy: &deletePolicy,
			},
			metav1.ListOptions{
				LabelSelector: selector.String(),
			},
		)
	}

	return util.IgnoreNotFoundError(err)
}

func DeleteIngresseWithName(namespace, name string, clientset *kubernetes.Clientset) error {
	version, err := clientset.Discovery().ServerVersion()
	if err != nil {
		return err
	}

	deletePolicy := metav1.DeletePropagationForeground
	if cluster.VersionLessThan122(version) {
		return clientset.ExtensionsV1beta1().Ingresses(namespace).Delete(
			context.TODO(), name,
			metav1.DeleteOptions{
				PropagationPolicy: &deletePolicy,
			})
	}

	return clientset.NetworkingV1().Ingresses(namespace).Delete(
		context.TODO(), name,
		metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		},
	)
}

func CreateIngress(namespace string, ingress *v1.Ingress, clientset *kubernetes.Clientset) error {
	_, err := clientset.NetworkingV1().Ingresses(namespace).Create(
		context.TODO(), ingress,
		metav1.CreateOptions{},
	)
	return err
}

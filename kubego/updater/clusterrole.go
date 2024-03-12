package updater

import (
	"context"
	"github.com/illuminatingKong/kongming-kit/kubego/cluster"
	"github.com/illuminatingKong/kongming-kit/kubego/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func DeleteClusterRoles(selector labels.Selector, clientset *kubernetes.Clientset) error {
	version, err := clientset.Discovery().ServerVersion()
	if err != nil {
		return err
	}

	deletePolicy := metav1.DeletePropagationForeground

	if cluster.VersionLessThan122(version) {
		err = clientset.RbacV1beta1().ClusterRoles().DeleteCollection(context.TODO(), metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}, metav1.ListOptions{
			LabelSelector: selector.String(),
		})
	} else {
		err = clientset.RbacV1().ClusterRoles().DeleteCollection(context.TODO(), metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}, metav1.ListOptions{
			LabelSelector: selector.String(),
		})
	}

	return util.IgnoreNotFoundError(err)
}

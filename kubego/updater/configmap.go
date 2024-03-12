package updater

import (
	"context"
	"github.com/illuminatingKong/kongming-kit/kubego/util"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func DeleteConfigMaps(namespace string, selector labels.Selector, clientset *kubernetes.Clientset) error {
	deletePolicy := metav1.DeletePropagationForeground
	err := clientset.CoreV1().ConfigMaps(namespace).DeleteCollection(
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

func UpdateConfigMap(namespace string, cm *corev1.ConfigMap, clientset *kubernetes.Clientset) error {
	_, err := clientset.CoreV1().ConfigMaps(namespace).Update(
		context.TODO(),
		cm,
		metav1.UpdateOptions{},
	)
	return err
}

func DeleteConfigMap(ns, name string, cl client.Client) error {
	return deleteObjectWithDefaultOptions(&corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      name,
		},
	}, cl)
}

func CreateConfigMap(cm *corev1.ConfigMap, cl client.Client) error {
	return createObjectNeverAnnotation(cm, cl)
}

func DeleteConfigMapsAndWait(ns string, selector labels.Selector, cl client.Client) error {
	gvk := schema.GroupVersionKind{
		Group:   "",
		Kind:    "ConfigMap",
		Version: "v1",
	}
	return deleteObjectsAndWait(ns, selector, &corev1.ConfigMap{}, gvk, cl)
}

func CreateOrPatchConfigMap(cm *corev1.ConfigMap, cl client.Client) error {
	return createOrPatchObject(cm, cl)
}

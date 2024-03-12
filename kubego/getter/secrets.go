package getter

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetSecret(ns, name string, cl client.Client) (*corev1.Secret, bool, error) {
	svc := &corev1.Secret{}
	found, err := GetResourceInCache(ns, name, svc, cl)
	if err != nil || !found {
		svc = nil
	}

	return svc, found, err
}

func ListSecrets(ns string, cl client.Client) ([]*corev1.Secret, error) {
	l := &corev1.SecretList{}
	gvk := schema.GroupVersionKind{
		Group:   "core",
		Kind:    "Secret",
		Version: "v1",
	}
	l.SetGroupVersionKind(gvk)
	err := ListResourceInCache(ns, nil, nil, l, cl)
	if err != nil {
		return nil, err
	}

	var res []*corev1.Secret
	for i := range l.Items {
		res = append(res, &l.Items[i])
	}
	return res, err
}

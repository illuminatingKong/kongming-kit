package getter

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ListConfigMaps(ns string, selector labels.Selector, cl client.Client) ([]*corev1.ConfigMap, error) {
	l := &corev1.ConfigMapList{}
	gvk := schema.GroupVersionKind{
		Group:   "core",
		Kind:    "ConfigMap",
		Version: "v1",
	}
	l.SetGroupVersionKind(gvk)
	err := ListResourceInCache(ns, selector, nil, l, cl)
	if err != nil {
		return nil, err
	}

	var res []*corev1.ConfigMap
	for i := range l.Items {
		res = append(res, &l.Items[i])
	}
	return res, err
}

func GetConfigMap(ns, name string, cl client.Client) (*corev1.ConfigMap, bool, error) {
	g := &corev1.ConfigMap{}
	found, err := GetResourceInCache(ns, name, g, cl)
	if err != nil || !found {
		g = nil
	}

	return g, found, err
}

func ListConfigMapsYaml(ns string, selector labels.Selector, cl client.Client) ([][]byte, error) {
	gvk := schema.GroupVersionKind{
		Group:   "",
		Kind:    "ConfigMap",
		Version: "v1",
	}
	return ListResourceYamlInCache(ns, selector, nil, gvk, cl)
}

package getter

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetService(ns, name string, cl client.Client) (*corev1.Service, bool, error) {
	svc := &corev1.Service{}
	found, err := GetResourceInCache(ns, name, svc, cl)
	if err != nil || !found {
		svc = nil
	}

	return svc, found, err
}

func ListServices(ns string, selector labels.Selector, cl client.Client) ([]*corev1.Service, error) {
	ss := &corev1.ServiceList{}
	err := ListResourceInCache(ns, selector, nil, ss, cl)
	if err != nil {
		return nil, err
	}

	var res []*corev1.Service
	for i := range ss.Items {
		res = append(res, &ss.Items[i])
	}
	return res, err
}

func ListServicesWithCache(selector labels.Selector, lister informers.SharedInformerFactory) ([]*corev1.Service, error) {
	if selector == nil {
		selector = labels.NewSelector()
	}
	return lister.Core().V1().Services().Lister().List(selector)
}

func ListServicesYaml(ns string, selector labels.Selector, cl client.Client) ([][]byte, error) {
	gvk := schema.GroupVersionKind{
		Group:   "",
		Kind:    "Service",
		Version: "v1",
	}
	return ListResourceYamlInCache(ns, selector, nil, gvk, cl)
}

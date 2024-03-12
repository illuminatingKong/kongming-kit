package getter

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetStatefulSet(ns, name string, cl client.Client) (*appsv1.StatefulSet, bool, error) {
	ss := &appsv1.StatefulSet{}
	found, err := GetResourceInCache(ns, name, ss, cl)
	if err != nil || !found {
		ss = nil
	}

	return ss, found, err
}

func ListStatefulSets(ns string, selector labels.Selector, cl client.Client) ([]*appsv1.StatefulSet, error) {
	ss := &appsv1.StatefulSetList{}
	err := ListResourceInCache(ns, selector, nil, ss, cl)
	if err != nil {
		return nil, err
	}

	var res []*appsv1.StatefulSet
	for i := range ss.Items {
		res = append(res, &ss.Items[i])
	}
	return res, err
}

func ListStatefulSetsWithCache(selector labels.Selector, lister informers.SharedInformerFactory) ([]*appsv1.StatefulSet, error) {
	if selector == nil {
		selector = labels.NewSelector()
	}
	return lister.Apps().V1().StatefulSets().Lister().List(selector)
}

func ListStatefulSetsYaml(ns string, selector labels.Selector, cl client.Client) ([][]byte, error) {
	gvk := schema.GroupVersionKind{
		Group:   "apps",
		Kind:    "StatefulSet",
		Version: "v1",
	}
	return ListResourceYamlInCache(ns, selector, nil, gvk, cl)
}

func GetStatefulSetYaml(ns string, name string, cl client.Client) ([]byte, bool, error) {
	gvk := schema.GroupVersionKind{
		Group:   "apps",
		Kind:    "StatefulSet",
		Version: "v1",
	}
	return GetResourceYamlInCache(ns, name, gvk, cl)
}

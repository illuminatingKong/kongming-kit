package getter

import (
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetExtensionsV1Beta1Ingress(namespace, name string, lister informers.SharedInformerFactory) (*extensionsv1beta1.Ingress, bool, error) {
	ret, err := lister.Extensions().V1beta1().Ingresses().Lister().Ingresses(namespace).Get(name)
	if err == nil {
		return ret, true, nil
	}
	return nil, false, err
}

func GetNetworkingV1Ingress(namespace, name string, lister informers.SharedInformerFactory) (*v1.Ingress, error) {
	return lister.Networking().V1().Ingresses().Lister().Ingresses(namespace).Get(name)
}

// ListExtensionsV1Beta1Ingresses gets the ingress (extensions/v1beta1) from the informer
func ListExtensionsV1Beta1Ingresses(selector labels.Selector, lister informers.SharedInformerFactory) ([]*extensionsv1beta1.Ingress, error) {
	if selector == nil {
		selector = labels.NewSelector()
	}
	return lister.Extensions().V1beta1().Ingresses().Lister().List(selector)
}

func ListNetworkingV1Ingress(selector labels.Selector, lister informers.SharedInformerFactory) ([]*v1.Ingress, error) {
	if selector == nil {
		selector = labels.NewSelector()
	}
	return lister.Networking().V1().Ingresses().Lister().List(selector)
}

func ListIngressesYaml(ns string, selector labels.Selector, cl client.Client) ([][]byte, error) {
	gvk := schema.GroupVersionKind{
		Group:   "extensions",
		Kind:    "Ingress",
		Version: "v1beta1",
	}
	return ListResourceYamlInCache(ns, selector, nil, gvk, cl)
}

func ListIngresses(namespace string, cl client.Client, lessThan122 bool) (*unstructured.UnstructuredList, error) {
	gvk := schema.GroupVersionKind{
		Group:   "extensions",
		Kind:    "Ingress",
		Version: "v1beta1",
	}
	if !lessThan122 {
		gvk = schema.GroupVersionKind{
			Group:   "networking.k8s.io",
			Kind:    "Ingress",
			Version: "v1",
		}
	}
	u := &unstructured.UnstructuredList{}
	u.SetGroupVersionKind(gvk)

	err := ListResourceInCache(namespace, labels.Everything(), nil, u, cl)
	if err != nil {
		return u, err
	}
	return u, err
}

package getter

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetDeployment(ns, name string, cl client.Client) (*appsv1.Deployment, bool, error) {
	g := &appsv1.Deployment{}
	found, err := GetResourceInCache(ns, name, g, cl)
	if err != nil || !found {
		g = nil
	}

	return g, found, err
}

func ListDeployments(ns string, selector labels.Selector, cl client.Client) ([]*appsv1.Deployment, error) {
	ss := &appsv1.DeploymentList{}
	err := ListResourceInCache(ns, selector, nil, ss, cl)
	if err != nil {
		return nil, err
	}

	var res []*appsv1.Deployment
	for i := range ss.Items {
		res = append(res, &ss.Items[i])
	}
	return res, err
}

func ListDeploymentsWithCache(selector labels.Selector, lister informers.SharedInformerFactory) ([]*appsv1.Deployment, error) {
	if selector == nil {
		selector = labels.NewSelector()
	}
	return lister.Apps().V1().Deployments().Lister().List(selector)
}

func ListDeploymentsYaml(ns string, selector labels.Selector, cl client.Client) ([][]byte, error) {
	gvk := schema.GroupVersionKind{
		Group:   "apps",
		Kind:    "Deployment",
		Version: "v1",
	}
	return ListResourceYamlInCache(ns, selector, nil, gvk, cl)
}

func GetDeploymentYaml(ns string, name string, cl client.Client) ([]byte, bool, error) {
	gvk := schema.GroupVersionKind{
		Group:   "apps",
		Kind:    "Deployment",
		Version: "v1",
	}
	return GetResourceYamlInCache(ns, name, gvk, cl)
}

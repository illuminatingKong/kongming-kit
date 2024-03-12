package getter

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetPod(ns, name string, cl client.Client) (*corev1.Pod, bool, error) {
	p := &corev1.Pod{}
	found, err := GetResourceInCache(ns, name, p, cl)
	if err != nil || !found {
		p = nil
	}

	return p, found, err
}

func ListPods(ns string, selector labels.Selector, cl client.Client) ([]*corev1.Pod, error) {
	ps := &corev1.PodList{}
	err := ListResourceInCache(ns, selector, nil, ps, cl)
	if err != nil {
		return nil, err
	}

	var res []*corev1.Pod
	for i := range ps.Items {
		res = append(res, &ps.Items[i])
	}
	return res, err
}

func ListPodsWithCache(selector labels.Selector, informer informers.SharedInformerFactory) ([]*corev1.Pod, error) {
	if selector == nil {
		selector = labels.NewSelector()
	}

	return informer.Core().V1().Pods().Lister().List(selector)
}

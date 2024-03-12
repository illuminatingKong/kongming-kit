package getter

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ListEvents(ns string, selector fields.Selector, cl client.Reader) ([]*corev1.Event, error) {
	es := &corev1.EventList{}
	err := ListResourceInCache(ns, nil, selector, es, cl)
	if err != nil {
		return nil, err
	}

	var res []*corev1.Event
	for i := range es.Items {
		res = append(res, &es.Items[i])
	}
	return res, err
}

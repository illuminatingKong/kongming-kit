package getter

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ListPvcs(ns string, selector fields.Selector, cl client.Reader) ([]*corev1.PersistentVolumeClaim, error) {
	pvcList := &corev1.PersistentVolumeClaimList{}
	gvk := schema.GroupVersionKind{
		Group:   "core",
		Kind:    "PersistentVolumeClaim",
		Version: "v1",
	}
	pvcList.SetGroupVersionKind(gvk)
	err := ListResourceInCache(ns, nil, selector, pvcList, cl)
	if err != nil {
		return nil, err
	}

	var res []*corev1.PersistentVolumeClaim
	for i := range pvcList.Items {
		res = append(res, &pvcList.Items[i])
	}
	return res, err
}

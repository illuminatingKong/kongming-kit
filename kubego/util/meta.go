package util

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GetGroupVersionKind(obj runtime.Object) (*schema.GroupVersionKind, error) {
	t, err := meta.TypeAccessor(obj)
	if err != nil {
		return nil, err
	}

	gvk := schema.FromAPIVersionAndKind(t.GetAPIVersion(), t.GetKind())

	return &gvk, nil
}

func SetGroupVersionKind(obj runtime.Object, gvk *schema.GroupVersionKind) (runtime.Object, error) {
	t, err := meta.TypeAccessor(obj)
	if err != nil {
		return nil, err
	}

	t.SetAPIVersion(gvk.GroupVersion().String())
	t.SetKind(gvk.Kind)

	return obj, nil
}

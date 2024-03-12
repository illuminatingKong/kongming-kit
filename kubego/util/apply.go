package util

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubectl/pkg/util"
)

var metadataAccessor = meta.NewAccessor()

func CreateApplyAnnotation(obj runtime.Object) error {
	return util.CreateApplyAnnotation(obj, unstructured.UnstructuredJSONScheme)
}

func GetModifiedConfiguration(obj runtime.Object) ([]byte, error) {
	return util.GetModifiedConfiguration(obj, true, unstructured.UnstructuredJSONScheme)
}

// GetOriginalConfiguration retrieves the original configuration of the object
// from the annotation, or nil if no annotation was found.
func GetOriginalConfiguration(obj runtime.Object) ([]byte, error) {
	annots, err := metadataAccessor.Annotations(obj)
	if err != nil {
		return nil, err
	}

	if annots == nil {
		return nil, nil
	}

	original, ok := annots[corev1.LastAppliedConfigAnnotation]
	if !ok {
		return nil, nil
	}

	return []byte(original), nil
}

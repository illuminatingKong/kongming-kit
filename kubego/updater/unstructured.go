package updater

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateOrPatchUnstructured(u *unstructured.Unstructured, cl client.Client) error {
	return createOrPatchObject(u, cl)
}

func CreateOrPatchUnstructuredNeverAnnotation(u *unstructured.Unstructured, cl client.Client) error {
	return createOrPatchObjectNeverAnnotation(u, cl)
}

func UpdateOrCreateUnstructured(u *unstructured.Unstructured, cl client.Client) error {
	return updateOrCreateObject(u, cl)
}

func DeleteUnstructured(u *unstructured.Unstructured, cl client.Client) error {
	return deleteObjectWithDefaultOptions(u, cl)
}

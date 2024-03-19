package util

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"helm.sh/helm/v3/pkg/releaseutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func NewResources(parsedYaml string, rs *runtime.Scheme) ([]*unstructured.Unstructured, error) {
	var err *multierror.Error
	resources, errs := Verify(parsedYaml)
	if len(errs) > 0 {
		err.Errors = errs
		return resources, err.ErrorOrNil()
	}
	return resources, err.ErrorOrNil()
}

func Verify(yamlStr string) ([]*unstructured.Unstructured, []error) {
	manifests := releaseutil.SplitManifests(yamlStr)
	resources := make([]*unstructured.Unstructured, 0)
	scheme := runtime.NewScheme()
	var Err []error
	for _, item := range manifests {
		u, err := NewDecoder(scheme).YamlToUnstructured([]byte(item))
		if err != nil {
			e := fmt.Sprintf("failed to convert yaml to Unstructured, manifest is %v , error: %v", item, err)
			ne := errors.New(e)
			Err = append(Err, ne)
			continue
		}
		resources = append(resources, u)
	}
	return resources, Err
}

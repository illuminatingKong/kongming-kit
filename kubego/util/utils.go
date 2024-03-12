package util

import (
	"bytes"
	"fmt"
	"io"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
)

type ObjGVK struct {
	Object runtime.Object
	GVK    *schema.GroupVersionKind
}

func IgnoreNotFoundError(err error) error {
	if apierrors.IsNotFound(err) {
		return nil
	}

	return err
}

func ParseManifest(manifest string) ([]ObjGVK, error) {
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(manifest)), 5*1024*1024)

	res := []ObjGVK{}
	for {
		var rawObj runtime.RawExtension

		err := decoder.Decode(&rawObj)
		if err == io.EOF {
			break
		}
		if err != nil {
			return res, fmt.Errorf("failed to decode raw extension: %s", err)
		}

		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if err != nil {
			return res, fmt.Errorf("failed to decode RawExtension to runtime.Object and GVK: %s", err)
		}

		res = append(res, ObjGVK{
			Object: obj,
			GVK:    gvk,
		})
	}

	return res, nil
}

func GetSvcNamesFromManifest(manifest string) ([]string, error) {
	objgvks, err := ParseManifest(manifest)
	if err != nil {
		return nil, err
	}

	svcNames := []string{}
	for _, objgvk := range objgvks {
		if objgvk.GVK.GroupKind().String() != "Service" {
			continue
		}

		objMeta, err := meta.Accessor(objgvk.Object)
		if err != nil {
			return nil, fmt.Errorf("failed to construct object meta: %s", err)
		}

		svcNames = append(svcNames, objMeta.GetName())
	}

	return svcNames, nil
}

package util

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

type decoder struct {
	typedDeserializer      runtime.Decoder
	unstructuredSerializer runtime.Serializer
}

func NewDecoder(rs *runtime.Scheme) *decoder {
	d := &decoder{
		typedDeserializer:      serializer.NewCodecFactory(rs).UniversalDeserializer(),
		unstructuredSerializer: yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme),
	}
	return d
}

func (d *decoder) YamlToUnstructured(manifest []byte) (*unstructured.Unstructured, error) {
	obj := &unstructured.Unstructured{}
	decoded, _, err := d.unstructuredSerializer.Decode(manifest, nil, obj)

	u, ok := decoded.(*unstructured.Unstructured)
	if !ok {
		return nil, fmt.Errorf("object is not an Unstructured")
	}

	return u, err
}

func (d *decoder) JSONToRuntimeObject(manifest []byte) (runtime.Object, error) {
	return d.YamlToRuntimeObject(manifest)
}

func (d *decoder) YamlToRuntimeObject(manifest []byte) (runtime.Object, error) {
	obj, _, err := d.typedDeserializer.Decode(manifest, nil, nil)

	return obj, err
}

func (d *decoder) YamlToDeployment(manifest []byte) (*appsv1.Deployment, error) {
	obj, err := d.YamlToRuntimeObject(manifest)
	if err != nil {
		return nil, err
	}

	res, ok := obj.(*appsv1.Deployment)
	if !ok {
		return nil, fmt.Errorf("object is not a Deployment")
	}

	return res, err
}

func (d *decoder) YamlToStatefulSet(manifest []byte) (*appsv1.StatefulSet, error) {
	obj, err := d.YamlToRuntimeObject(manifest)
	if err != nil {
		return nil, err
	}

	res, ok := obj.(*appsv1.StatefulSet)
	if !ok {
		return nil, fmt.Errorf("object is not a StatefulSet")
	}

	return res, err
}

func (d *decoder) YamlToIngress(manifest []byte) (*extensionsv1beta1.Ingress, error) {
	obj, err := d.YamlToRuntimeObject(manifest)
	if err != nil {
		return nil, err
	}

	res, ok := obj.(*extensionsv1beta1.Ingress)
	if !ok {
		return nil, fmt.Errorf("object is not a Ingress")
	}

	return res, err
}

func (d *decoder) YamlToJob(manifest []byte) (*batchv1.Job, error) {
	obj, err := d.YamlToRuntimeObject(manifest)
	if err != nil {
		return nil, err
	}

	res, ok := obj.(*batchv1.Job)
	if !ok {
		return nil, fmt.Errorf("object is not a Job")
	}

	return res, err
}

func (d *decoder) YamlToCronJob(manifest []byte) (*batchv1beta1.CronJob, error) {
	obj, err := d.YamlToRuntimeObject(manifest)
	if err != nil {
		return nil, err
	}

	res, ok := obj.(*batchv1beta1.CronJob)
	if !ok {
		return nil, fmt.Errorf("object is not a CronJob")
	}

	return res, err
}

func (d *decoder) YamlToConfigMap(manifest []byte) (*corev1.ConfigMap, error) {
	obj, err := d.YamlToRuntimeObject(manifest)
	if err != nil {
		return nil, err
	}

	res, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return nil, fmt.Errorf("object is not a ConfigMap")
	}

	return res, err
}

func (d *decoder) JSONToDeployment(manifest []byte) (*appsv1.Deployment, error) {
	return d.YamlToDeployment(manifest)
}

func (d *decoder) JSONToStatefulSet(manifest []byte) (*appsv1.StatefulSet, error) {
	return d.YamlToStatefulSet(manifest)
}

func (d *decoder) JSONToJob(manifest []byte) (*batchv1.Job, error) {
	return d.YamlToJob(manifest)
}

func (d *decoder) JSONToCronJob(manifest []byte) (*batchv1beta1.CronJob, error) {
	return d.YamlToCronJob(manifest)
}

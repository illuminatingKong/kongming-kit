package yaml

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/illuminatingKong/kongming-kit/base/logx"
	"github.com/illuminatingKong/kongming-kit/kubego/util"
	"github.com/illuminatingKong/kongming-kit/rollout"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

func NewDeploy(name string) *Deploy {
	gr := schema.GroupResource{
		Group:    "apps",
		Resource: "deployments",
	}
	return &Deploy{
		Name:          name,
		groupResource: gr,
	}
}

func (d *Deploy) Log() logx.Logger {
	return rollout.InitLog().Sub("Deployment")
}

func (d *Deploy) Verify(yamlText string) []error {
	_, err := util.Verify(yamlText)
	return err
}

func (d *Deploy) UpdatedAnnotations(annotations map[string]string) map[string]string {
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations["updated-by-devops"] = fmt.Sprintf("%d", time.Now().Unix())
	return annotations
}

func Scheme(cl client.Client) *runtime.Scheme {
	return cl.Scheme()
}

// ReNew  can use a fake scheme runtime.NewScheme()
func (d *Deploy) ReNew(yamlText string, rs *runtime.Scheme, opts Option) ([]runtime.Object, error) {
	resources, err := util.NewResources(yamlText, rs)
	if err != nil {
		d.Log().Errorf("failed to create deployment resources, error: %v", err)
		return nil, err
	}
	var multiResource = make([]runtime.Object, 0)
	var errs *multierror.Error
	for _, u := range resources {
		switch u.GetKind() {
		case "Deployment":
			if opts.NameSpace != "" {
				u.SetNamespace(opts.NameSpace)
			}

			if len(opts.Labels) > 0 {
				labels := u.GetLabels()
				if labels == nil {
					labels = make(map[string]string)
				}
				for k, v := range opts.Labels {
					labels[k] = v
				}
				u.SetLabels(labels)
				podLabels, _, err := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "labels")
				if err != nil {
					podLabels = nil
				}
				err = unstructured.SetNestedStringMap(u.Object, util.MergeLabels(labels, podLabels), "spec", "template", "metadata", "labels")
				if err != nil {
					d.Log().Errorf("merge label failed err:%s", err)
					u.Object = util.SetFieldValueIsNotExist(u.Object, util.MergeLabels(labels, podLabels), "spec", "template", "metadata", "labels")
				}

				podAnnotations, _, err := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "annotations")
				if err != nil {
					podAnnotations = nil
				}
				err = unstructured.SetNestedStringMap(u.Object, d.UpdatedAnnotations(podAnnotations), "spec", "template", "metadata", "annotations")
				if err != nil {
					d.Log().Errorf("merge annotation failed err:%s", err)
					u.Object = util.SetFieldValueIsNotExist(u.Object, d.UpdatedAnnotations(podAnnotations), "spec", "template", "metadata", "annotations")
				}
			}
			jsonData, err := u.MarshalJSON()
			if err != nil {
				d.Log().Errorf("Failed to marshal JSON, manifest is\n%v\n, error: %v", u, err)
				errs = multierror.Append(errs, err)
				continue
			}
			obj, err := util.NewDecoder(rs).JSONToRuntimeObject(jsonData)
			if err != nil {
				d.Log().Errorf("Failed to convert JSON to Object, manifest is\n%v\n, error: %v", u, err)
				errs = multierror.Append(errs, err)
				continue
			}
			//return obj, errs.ErrorOrNil()
			errs = multierror.Append(errs, err)

			multiResource = append(multiResource, obj)
		}
	}
	if errs.Len() == 0 {
		return multiResource, errs.ErrorOrNil()
	}
	return nil, errs.ErrorOrNil()
}

package yaml

import (
	"github.com/hashicorp/go-multierror"
	"github.com/illuminatingKong/kongming-kit/kubego/updater"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (d *Deploy) Apply(objs []runtime.Object, cl client.Client) error {

	var errs *multierror.Error
	for _, obj := range objs {
		switch res := obj.(type) {
		case *appsv1.Deployment:
			err := updater.CreateOrPatchDeployment(res, cl)
			if err != nil {
				errs = multierror.Append(errs, err)
			}
		case *corev1.ConfigMap:
			err := updater.CreateOrPatchConfigMap(res, cl)
			if err != nil {
				errs = multierror.Append(errs, err)
			}
		case *corev1.Namespace:
			err := updater.CreateOrPatchNamespace(res, cl)
			if err != nil {
				errs = multierror.Append(errs, err)
			}
		case *corev1.Service:
			err := updater.CreateOrPatchService(res, cl)
			if err != nil {
				errs = multierror.Append(errs, err)
			}
		}
	}

	return errs.ErrorOrNil()
}

func (d *Deploy) CreateOrPatch(yamlText string, cl client.Client, opts Option) error {
	var errs *multierror.Error
	verifyErr := d.Verify(yamlText)
	if verifyErr != nil {
		errs = multierror.Append(errs, verifyErr...)
		return errs
	}
	rs := cl.Scheme()
	obj, err := d.ReNew(yamlText, rs, opts)
	if err != nil {
		return err
	}
	return d.Apply(obj, cl)
}

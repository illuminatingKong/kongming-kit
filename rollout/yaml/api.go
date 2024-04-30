package yaml

import (
	"github.com/hashicorp/go-multierror"
	"github.com/illuminatingKong/kongming-kit/kubego/updater"
	"github.com/illuminatingKong/kongming-kit/kubego/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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

func (d *Deploy) CreateOrPatchCustomResource(yamlText string, cl client.Client) error {
	var errs *multierror.Error
	multiunstructured, verifyErr := util.Verify(yamlText)
	if verifyErr != nil {
		errs = multierror.Append(errs, verifyErr...)
		return errs
	}
	for _, unstructured := range multiunstructured {
		err := updater.CreateOrPatchUnstructured(unstructured, cl)
		errs = multierror.Append(errs, err)
	}

	return errs.ErrorOrNil()
}

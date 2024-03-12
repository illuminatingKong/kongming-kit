package updater

import (
	"context"
	"github.com/illuminatingKong/kongming-kit/kubego/getter"
	"github.com/illuminatingKong/kongming-kit/kubego/patcher"
	"github.com/illuminatingKong/kongming-kit/kubego/util"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

func patchObject(obj client.Object, patchBytes []byte, cl client.Client) error {
	return cl.Patch(context.TODO(), obj, client.RawPatch(types.StrategicMergePatchType, patchBytes))
}

func updateObject(obj client.Object, cl client.Client) error {
	// always add the last-applied-configuration annotation
	err := util.CreateApplyAnnotation(obj)
	if err != nil {
		return err
	}
	return cl.Update(context.TODO(), obj)
}

func deleteObject(obj client.Object, cl client.Client, opts ...client.DeleteOption) error {
	return cl.Delete(context.TODO(), obj, opts...)
}

func deleteObjects(obj client.Object, cl client.Client, opts ...client.DeleteAllOfOption) error {
	return cl.DeleteAllOf(context.TODO(), obj, opts...)
}

func deleteObjectWithDefaultOptions(obj client.Object, cl client.Client) error {
	deletePolicy := metav1.DeletePropagationForeground
	return deleteObject(obj, cl, &client.DeleteOptions{PropagationPolicy: &deletePolicy})
}

func deleteObjectsWithDefaultOptions(ns string, selector labels.Selector, obj client.Object, cl client.Client) error {
	deletePolicy := metav1.DeletePropagationForeground

	delOpt := &client.DeleteAllOfOptions{
		DeleteOptions: client.DeleteOptions{PropagationPolicy: &deletePolicy},
		ListOptions:   client.ListOptions{LabelSelector: selector, Namespace: ns},
	}

	return deleteObjects(obj, cl, delOpt)
}

func createObject(obj client.Object, cl client.Client) error {
	// always add the last-applied-configuration annotation
	err := util.CreateApplyAnnotation(obj)
	if err != nil {
		return err
	}
	return cl.Create(context.TODO(), obj)
}

func createObjectNeverAnnotation(obj client.Object, cl client.Client) error {
	return cl.Create(context.TODO(), obj)
}

func updateOrCreateObject(obj client.Object, cl client.Client) error {
	err := updateObject(obj, cl)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return createObject(obj, cl)
		}

		return err
	}

	return nil
}

func createOrPatchObject(modified client.Object, cl client.Client) error {
	c := modified.DeepCopyObject().(client.Object)
	found, err := getter.GetResourceInCache(modified.GetNamespace(), modified.GetName(), c, cl)
	if err != nil {
		return err
	} else if !found {
		return createObject(modified, cl)
	}

	// fill gvk in case it is missing.
	// for objects retrieved from APIServer, gvk is not set, so we need to set it
	// for objects in cache, gvk is set, nothing will be changed here
	gvk := modified.GetObjectKind().GroupVersionKind()
	current, err := util.SetGroupVersionKind(c, &gvk)
	if err != nil {
		return err
	}

	patchBytes, patchType, err := patcher.GeneratePatchBytes(current, modified, cl)
	if err != nil {
		return err
	}

	if len(patchBytes) == 0 || patchType == "" || string(patchBytes) == "{}" {
		return nil
	}

	return patchObjectWithType(current.(client.Object), patchBytes, patchType, cl)
}

func createOrPatchObjectNeverAnnotation(modified client.Object, cl client.Client) error {
	c := modified.DeepCopyObject().(client.Object)
	found, err := getter.GetResourceInCache(modified.GetNamespace(), modified.GetName(), c, cl)
	if err != nil {
		return err
	} else if !found {
		return createObjectNeverAnnotation(modified, cl)
	}

	gvk := modified.GetObjectKind().GroupVersionKind()
	current, err := util.SetGroupVersionKind(c, &gvk)
	if err != nil {
		return err
	}

	patchBytes, patchType, err := patcher.GeneratePatchBytes(current, modified, cl)
	if err != nil {
		return err
	}

	if len(patchBytes) == 0 || patchType == "" || string(patchBytes) == "{}" {
		return nil
	}

	return patchObjectWithType(current.(client.Object), patchBytes, patchType, cl)
}

func patchObjectWithType(obj client.Object, patchBytes []byte, patchType types.PatchType, cl client.Client) error {
	return cl.Patch(context.TODO(), obj, client.RawPatch(patchType, patchBytes))
}

// TODO: LOU（LOU is from zhangmen）: improve it
// deleteObjectAndWait delete the object and wait till it is gone.
func deleteObjectAndWait(obj client.Object, cl client.Client) error {
	err := deleteObjectWithDefaultOptions(obj, cl)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	return wait.PollImmediate(time.Second, 60*time.Second, func() (done bool, err error) {
		found, err := getter.GetResourceInCache(obj.GetNamespace(), obj.GetName(), obj, cl)
		if err != nil {
			return false, err
		}

		return !found, nil
	})
}

// TODO: LOU（LOU is form zhangmen）: improve it
// deleteObjectsAndWait deletes the objects and wait till all of them are gone.
func deleteObjectsAndWait(ns string, selector labels.Selector, obj client.Object, gvk schema.GroupVersionKind, cl client.Client) error {
	err := deleteObjectsWithDefaultOptions(ns, selector, obj, cl)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	return wait.PollImmediate(time.Second, 60*time.Second, func() (done bool, err error) {
		us, err := getter.ListUnstructuredResourceInCache(ns, selector, nil, gvk, cl)
		if err != nil {
			return false, err
		}

		return len(us) == 0, nil
	})
}

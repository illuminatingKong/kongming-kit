package updater

import (
	"context"
	"github.com/illuminatingKong/kongming-kit/kubego/util"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func DeleteJobs(namespace string, selector labels.Selector, clientset *kubernetes.Clientset) error {
	deletePolicy := metav1.DeletePropagationForeground
	err := clientset.BatchV1().Jobs(namespace).DeleteCollection(
		context.TODO(),
		metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		},
		metav1.ListOptions{
			LabelSelector: selector.String(),
		},
	)

	return util.IgnoreNotFoundError(err)
}

func CreateJob(job *batchv1.Job, cl client.Client) error {
	return createObject(job, cl)
}

func DeleteJob(ns, name string, cl client.Client) error {
	return deleteObjectWithDefaultOptions(&batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      name,
		},
	}, cl)
}

func DeleteJobAndWait(ns, name string, cl client.Client) error {
	return deleteObjectAndWait(&batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      name,
		},
	}, cl)
}

func DeleteJobsAndWait(ns string, selector labels.Selector, cl client.Client) error {
	gvk := schema.GroupVersionKind{
		Group:   "batch",
		Kind:    "Job",
		Version: "v1",
	}
	return deleteObjectsAndWait(ns, selector, &batchv1.Job{}, gvk, cl)
}

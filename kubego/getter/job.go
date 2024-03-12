package getter

import (
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ListJobs(ns string, selector labels.Selector, cl client.Client) ([]*batchv1.Job, error) {
	jobs := &batchv1.JobList{}
	err := ListResourceInCache(ns, selector, nil, jobs, cl)
	if err != nil {
		return nil, err
	}

	var res []*batchv1.Job
	for i := range jobs.Items {
		res = append(res, &jobs.Items[i])
	}
	return res, err
}

func GetJob(ns, name string, cl client.Client) (*batchv1.Job, bool, error) {
	g := &batchv1.Job{}
	found, err := GetResourceInCache(ns, name, g, cl)
	if err != nil || !found {
		g = nil
	}

	return g, found, err
}

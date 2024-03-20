package query

import (
	"github.com/illuminatingKong/kongming-kit/kubego/getter"
	Pod "github.com/illuminatingKong/kongming-kit/rollout/pod"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/informers"
)

var (
	PodReady      = "ready"
	JobReady      = "Completed"
	PodNotReady   = "NotReady"
	PodRunning    = "Running"
	PodError      = "Error"
	PodUnstable   = "Unstable"
	PodNonStarted = "UnStart"
)

func QueryPodsStatus(informer informers.SharedInformerFactory, label map[string]string) (string, string, []string, []Pod.SimplePod) {
	ls := labels.Set{}
	if len(label) > 0 {
		for k, v := range label {
			ls[k] = v
		}
	}
	return GetSelectedPodsInfo(ls.AsSelector(), informer)
}

func GetSelectedPodsInfo(selector labels.Selector, informer informers.SharedInformerFactory) (string, string, []string, []Pod.SimplePod) {
	pods, err := getter.ListPodsWithCache(selector, informer)
	if err != nil {
		return PodError, PodNotReady, nil, nil
	}
	if len(pods) == 0 {
		return PodNonStarted, PodNotReady, nil, nil
	}
	imageSet := sets.String{}
	var sPods []Pod.SimplePod
	for _, pod := range pods {
		ipod := Pod.Pod(pod)
		s := ipod.Transform()
		sPods = append(sPods, s)
		imageSet.Insert(ipod.Containers()...)
	}
	images := imageSet.List()
	ready := PodReady
	succeededPods := 0
	for _, pod := range pods {
		iPod := Pod.Pod(pod)

		if iPod.Succeeded() {
			succeededPods++
			continue
		}

		if !iPod.Ready() {
			return PodUnstable, PodNotReady, images, sPods
		}
	}

	if len(pods) == succeededPods {
		return string(corev1.PodSucceeded), JobReady, images, sPods
	}

	return PodRunning, ready, images, sPods
}

package pod

import corev1 "k8s.io/api/core/v1"

type pod struct {
	*corev1.Pod
}

func Pod(w *corev1.Pod) *pod {
	if w == nil {
		return nil
	}

	return &pod{
		Pod: w,
	}
}

// Unwrap returns the corev1.Pod object.
func (w *pod) Unwrap() *corev1.Pod {
	return w.Pod
}

// Finished means the pod is finished and closed, usually it is a job pod
func (w *pod) Finished() bool {
	p := w.phase()
	return p == corev1.PodSucceeded || p == corev1.PodFailed
}

func (w *pod) Pending() bool {
	return w.phase() == corev1.PodPending
}

// Succeeded means the pod is succeeded and closed, usually it is a job pod
func (w *pod) Succeeded() bool {
	return w.phase() == corev1.PodSucceeded
}

// Failed means the pod is failed and closed, usually it is a job pod
func (w *pod) Failed() bool {
	return w.phase() == corev1.PodFailed
}

func (w *pod) phase() corev1.PodPhase {
	return w.Pod.Status.Phase
}

func (w *pod) Phase() string {
	return string(w.phase())
}

// Ready indicates that the pod is ready for traffic.
func (w *pod) Ready() bool {
	cs := w.Status.Conditions
	for _, c := range cs {
		if c.Type == corev1.PodReady {
			return c.Status == corev1.ConditionTrue
		}
	}
	return false
}

// For a Pod that uses custom conditions, that Pod is evaluated to be ready only when both the following statements apply:
// All containers in the Pod are ready.
// All conditions specified in readinessGates are True.
// When a Pod's containers are Ready but at least one custom condition is missing or False, the kubelet sets the Pod's condition to ContainersReady.
// https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-readiness-status
func (w *pod) ContainersReady() bool {
	cs := w.Status.Conditions
	for _, c := range cs {
		if c.Type == corev1.ContainersReady {
			return c.Status == corev1.ConditionTrue
		}
	}
	return false
}

func (w *pod) Containers() []string {
	var cs []string
	for _, c := range w.Spec.Containers {
		cs = append(cs, c.Image)
	}

	return cs
}

func (w *pod) ContainerNames() []string {
	var cs []string
	for _, c := range w.Spec.Containers {
		cs = append(cs, c.Name)
	}

	return cs
}

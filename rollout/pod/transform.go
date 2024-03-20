package pod

import (
	corev1 "k8s.io/api/core/v1"
	"time"
)

type SimplePod struct {
	Name      string
	Phase     corev1.PodPhase
	HostIP    string
	PodIP     string
	StartTime time.Time
	Container []simpleContainer
}

type simpleContainer struct {
	Name  string
	Ready bool
	Image string
}

func (w *pod) Transform() SimplePod {

	phase := w.Status.Phase
	hostIP := w.Status.HostIP
	podName := w.ObjectMeta.Name
	podIP := w.Status.PodIP
	startTime := w.Status.StartTime.Time

	s := SimplePod{
		Name:      podName,
		Phase:     phase,
		HostIP:    hostIP,
		PodIP:     podIP,
		StartTime: startTime,
	}
	var cs []simpleContainer
	for _, c := range w.Status.ContainerStatuses {
		s := simpleContainer{
			Name:  c.Name,
			Ready: c.Ready,
			Image: c.Image,
		}
		cs = append(cs, s)
	}
	s.Container = cs
	return s

}

package getter

import (
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ListNodes(cl client.Client) ([]*corev1.Node, error) {
	nodes := &corev1.NodeList{}
	err := ListResourceInCache("", nil, nil, nodes, cl)
	if err != nil {
		return nil, err
	}

	var res []*corev1.Node
	for i := range nodes.Items {
		res = append(res, &nodes.Items[i])
	}
	return res, err
}

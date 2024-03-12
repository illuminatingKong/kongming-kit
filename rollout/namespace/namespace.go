package namespace

import (
	"github.com/illuminatingKong/kongming-kit/base/logx"
	"github.com/illuminatingKong/kongming-kit/kubego/getter"
	"github.com/illuminatingKong/kongming-kit/kubego/updater"
	"github.com/illuminatingKong/kongming-kit/rollout"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Namespacer interface {
	Log() logx.Logger
}

type Namespace struct {
	Name          string
	Labels        map[string]string
	groupResource schema.GroupResource
}

func NewNamespace(name string, labels map[string]string) *Namespace {
	gr := schema.GroupResource{
		Group:    "",
		Resource: "namespaces",
	}
	return &Namespace{
		Name:          name,
		Labels:        labels,
		groupResource: gr,
	}
}

func (n *Namespace) Log() logx.Logger {
	return rollout.InitLog().Sub("Namespace")
}

func (n *Namespace) GetNamespace(cl client.Client) error {
	_, found, err := getter.GetNamespace(n.Name, cl)
	if err != nil {
		return err
	}
	if !found {
		return errors.NewNotFound(n.groupResource, n.Name)
	}
	return nil
}

func (n *Namespace) CreateNamespace(cl client.Client) error {
	err := updater.CreateNamespaceByName(n.Name, n.Labels, cl)
	if err != nil {
		return errors.NewNotFound(n.groupResource, n.Name)
	}
	return nil
}

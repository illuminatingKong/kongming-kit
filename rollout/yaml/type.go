package yaml

import (
	"github.com/illuminatingKong/kongming-kit/base/logx"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Yamler interface {
	Log() logx.Logger
	Verify(yamlText string) []error
	ReNew(yamlText string, rs *runtime.Scheme, opts Option) ([]runtime.Object, error)
}

type Deploy struct {
	Name          string
	groupResource schema.GroupResource
}

type Option struct {
	NameSpace string
	Labels    map[string]string
}

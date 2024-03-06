package example_configx

import (
	"fmt"
	"github.com/illuminatingKong/kongming-kit/base/configx/config"
	"testing"
)

func TestLoadYamlConfig(t *testing.T) {
	var dir config.ConfigsDir = "$HOME/workspace/illuminatingKong/kongming-kit/examples/example-configx"
	conf := config.NewFile(config.WithConfigsDir(dir),
		config.WithConfigType("yaml"), config.WithConfigName("test"))

	err := conf.Load()
	if err != nil {
		panic(err)
	}

	fmt.Println(conf.GetString("core.addr"))
}

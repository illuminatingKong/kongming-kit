package example_configx

import (
	"bytes"
	"github.com/illuminatingKong/kongming-kit/base/configx/config"
	"testing"
)

func TestLoadYamlConfig(t *testing.T) {
	var dir config.ConfigsDir = "$HOME/workspace/illuminatingKong/kongming-kit/examples/example-configx"
	conf := config.NewFile(
		config.WithConfigsDir(dir),
		config.WithConfigType("yaml"), config.WithConfigName("test"),
		config.WithReadType(config.FileType),
	)

	err := conf.Load()
	if err != nil {
		panic(err)
	}

	addr := conf.GetString("core.addr")
	if addr != "0.0.0.0:80" {
		t.Fatal("addr is not 0.0.0.0:80")
	}
	return
}

func TestLoadIOConfig(t *testing.T) {
	var yamlExample = []byte(`
core:
  addr: 0.0.0.0:80

`)
	conf := config.NewFile(
		config.WithConfigType("yaml"),
		config.WithReadType(config.IOType),
		config.WithConfigIO(bytes.NewBuffer(yamlExample)),
	)
	err := conf.Load()
	if err != nil {
		panic(err)
	}

	addr := conf.GetString("core.addr")
	if addr != "0.0.0.0:80" {
		t.Fatal("addr is not 0.0.0.0:80")
	}
	return

}

package example_newcontainer

import (
	"context"
	"github.com/illuminatingKong/kongming-kit/runner"
	"sync"
	"testing"
)

func TestNewContainer(t *testing.T) {
	appName := "test"
	addr := "127.0.0.1:8080"
	configDir := "$HOME/workspace/illuminatingKong/kongming-kit/examples/example-configx"
	container := runner.NewContainer(appName, addr).NewConfig(
		configDir, "yaml", appName)
	var once sync.Once
	err := container.InitBase(context.Background(), &once)
	if err != nil {
		panic(err)
	}
	container.Logger.Info("new container")
	conf := runner.GetConf()
	configAddr := conf.Get("core.addr")
	log := runner.GetLogger()
	log.Info(configAddr)
}

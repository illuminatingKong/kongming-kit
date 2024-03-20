package example_newcontainer

import (
	"context"
	"fmt"
	"github.com/illuminatingKong/kongming-kit/runner"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"sync"
	"testing"
)

func TestNewContainer(t *testing.T) {
	appName := "test"
	addr := "127.0.0.1:8080"
	configDir := "$HOME/workspace/illuminatingKong/kongming-kit/examples/example-configx"
	logPath := "your log path"
	container := runner.NewContainer(appName, addr).NewConfig(
		configDir, "yaml", appName)

	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  fmt.Sprintf("%s/your.log", logPath),
		logrus.ErrorLevel: fmt.Sprintf("%s/your.log", logPath),
	}
	//container.Logger.
	container.Logger.AddHook(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))

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

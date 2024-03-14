package example_new_bind

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/illuminatingKong/kongming-kit/http/corehandler"
	"github.com/illuminatingKong/kongming-kit/http/middleware"
	"github.com/illuminatingKong/kongming-kit/http/service"
	"github.com/illuminatingKong/kongming-kit/http/webapi"
	"github.com/illuminatingKong/kongming-kit/runner"
	"github.com/oklog/run"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"
)

var (
	httpServiceMode = "debug"
	projectName     = "test"
	addr            = "0.0.0.0:8081"
	configDir       = "$HOME/workspace/illuminatingKong/kongming-kit/examples/example-new-bind"
)

type Router struct{}

var respSuccess = webapi.RespDefault
var respError = webapi.RespDefault

func (*Router) Inject(router *gin.RouterGroup) {

	version := router.Group("/core")
	{
		version.GET("/version", showVersion)
		version.GET("/hello", hello)
	}

}

func showVersion(c *gin.Context) {
	ctx := corehandler.NewContext(c)
	defer func() { corehandler.WebHttpApiResponse(c, ctx) }()
	now := time.Now()
	conf := runner.GetConf()
	coreversion := conf.GetString("core.version")

	ctx.Resp, ctx.Err = respSuccess().SetCode(201).SetData(map[string]interface{}{"version": coreversion,
		"now": now}).SetMessage("get version").SetHttpCode(200).SetPage(1), nil

}

func hello(c *gin.Context) {
	ctx := corehandler.NewContext(c)
	defer func() { corehandler.WebHttpApiResponse(c, ctx) }()
	ctx.Resp, ctx.Err = respSuccess().SetData(map[string]interface{}{"hello": "kongming"}).SetMessage("say hello"), nil

}

func TestStartProject(t *testing.T) {
	o := runner.NewContainer(projectName, addr).NewConfig(
		configDir, "yaml", projectName)
	var once sync.Once
	err := o.InitBase(context.Background(), &once)
	if err != nil {
		panic(err)
	}

	project := NewBind(projectName, o)
	project.Start()
}

type Bootstrap struct {
	HttpServiceEngine service.HttpServiceEngine
	HttpServer        *http.Server
	Project           runner.Options
}

// LoadHttpService is a load http service function
func LoadHttpService(serverAddr string) (service.HttpServiceEngine, *http.Server) {
	log := runner.GetLogger()
	log.Infof("http server read loading")
	e := service.HttpServiceEngine{
		Addr: serverAddr,
		Middlewares: []gin.HandlerFunc{
			middleware.Response(),
			middleware.RequestID(runner.GetLogger()),
			gin.Recovery(),
		},
		Mode: httpServiceMode,
	}
	e.RouterGroup = map[string]service.RouterInjector{
		// web api
		"/health": new(Router),
	}
	s := e.Init()
	server := &http.Server{Addr: e.Addr, Handler: s}
	return e, server
}

func NewBind(appName string, o *runner.Options) *Bootstrap {
	b := &Bootstrap{}
	b.Project = *o

	b.HttpServiceEngine, b.HttpServer = LoadHttpService(o.Addr)
	return b
}

func (b *Bootstrap) Start() {
	o := b.Project
	var g run.Group
	var once sync.Once
	type closeOnce struct {
		C     chan struct{}
		once  sync.Once
		Close func()
	}
	reloadReady := &closeOnce{
		C: make(chan struct{}),
	}
	reloadReady.Close = func() {
		reloadReady.once.Do(func() {
			close(reloadReady.C)
		})
	}

	{
		// termination handler.
		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)
		cancel := make(chan struct{})
		g.Add(
			func() error {
				select {
				case <-term:
					o.Logger.Warn("received sigterm, exiting gracefully...")
					reloadReady.Close()
				case <-cancel:
					reloadReady.Close()
				}
				return nil
			},
			func(err error) {
				close(cancel)

			},
		)
	}

	{
		// http server handler.
		g.Add(
			func() error {
				defer o.Logger.Info("http server is returned")
				err := b.HttpServiceEngine.Start(context.TODO(), b.HttpServer)
				return err
			},
			func(err error) {
				err = b.HttpServer.Close()
				if err != nil {
					o.Logger.Info("http server force exit")
					os.Exit(127)
				}
			},
		)
	}

	{
		// watch config handler.
		cancel := make(chan struct{})
		g.Add(
			func() error {
				o.WithConfWatch(o.OptionsCtx, &once)
				<-cancel
				return nil
			},
			func(err error) {
				close(cancel)
			},
		)
	}

	{
		// project first log handler.
		cancel := make(chan struct{})
		g.Add(
			func() error {
				o.Logger.Infof("project start  in: %s", time.Now().String())
				<-cancel
				return nil
			},
			func(err error) {
				//o.Logger.Error("showdown app: %s, error: %s", o.Name, err.Error())
				close(cancel)
			},
		)
	}

	gErr := g.Run()
	if gErr != nil {
		panic(gErr)
	}

	o.Logger.Info("project exit gracefully")
}

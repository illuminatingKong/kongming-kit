# kongming-kit


## 项目介绍

对于尝试用golang进行开发，会遇到一些问题在开发过程中，浪费很多时间，kongming-kit 就是为了解决这些问题而生。

kongming-kit 的定位是一个快速开发工具库，提供了一些常用的工具类和基础功能，方便快速开发。


## 功能

- [x] 日志 
- [x] 配置文件
- [x] 容器
- [x] 启动服务bind
- [x] kubego
- [x] rollout
- [x] http client
- [ ] sample httpclient 
- [ ] uwin cmdb client
- [ ] grpc
- [ ] exec command

## kit包介绍

### 日志

kongming-kit 使用logrus作为日志库，提供了一个日志的封装，方便使用。

```
package main
func TestLogFiled(t *testing.T) {
	var Logger logx.Logger
	var formatter logrusx.JsonFormatter
	Logger = logrusx.New(logrusx.WithFormatter(formatter))
	today := time.Now().String()
	Logger.WithFieldsX(Logger.Info, "your_date_module", "today", today)
}
```

输出
```shell
{"@level":"info","msg":"your_date_module: today=2024-03-20 19:35:56.118606 +0800 CST m=+0.000195126","timestamp":"2024-03-20 19:35:56.118"}
```

Logger.WithFieldsX 是一个封装的方法，可以传入多个字段，方便输出日志, 不用研发再去拼接字符串或fmt.Sprintf 格式化。


### 配置文件

configx 读取配置文件，使用viper库，支持yaml、json、toml等格式的配置文件。

```
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
```

配额文件 test.yaml
```yaml
core:
  addr: 0.0.0.0:80

```

### 容器

new-container 创建一个容器，并且导入配置文件和日志.
new-container 是希望通过容器的方式，来管理配置文件和日志，方便在项目中使用，不需要每次都传递配置文件和日志对象。

与传统 IoC container 不同的是，new-container 是一个简单的容器，只提供了配置文件和日志的导入，方便在项目中使用。并且没有DI的功能。


使用场景上偏向于单元测试，或者是一些简单启动项目。

启动项目中对日志路径的要求，可以用logrus hook的方式，来实现日志的输出。

```shell
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

```



### bind

kongming-kit 不愿意做框架，但是bind 提供创建一个新的项目时，要启用多个功能的情况下提供了一个简单的方式。

1 创建一个容器，导入配置文件和日志

2 新建bind服务，绑定路由、业务方法和中间件

3 启动bind服务，并且通过启动顺序完成组件初始化。例如 配置文件刷新，依赖的中间件，优雅下线等，使用先进后出加载顺序。

4 定义webapi.RespDefault 结构体，用于返回数据，方便统一处理返回数据

5 定义一个业务方法，用于处理请求，返回数据

```shell
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

func (*Router) Inject(router *gin.RouterGroup) {

	version := router.Group("/core")
	{
		version.GET("/version", showVersion)
		version.GET("/hello", hello)
	}

}

func customResponseSerializer(c *gin.Context) (*corehandler.Context, *webapi.WebHTTPApi) {
	return func(c *gin.Context) (*corehandler.Context, *webapi.WebHTTPApi) {
		ctx := corehandler.NewContext(c)
		resp := respSuccess().SetOmitEmptyKeys("extra", "page", "total", "limit")
		return ctx, resp
	}(c)
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

// hello is a hello function and use customResponseSerializer will return a custom response serializer
func hello(c *gin.Context) {
	ctx, resp := customResponseSerializer(c)
	defer func() { corehandler.WebHttpApiResponse(c, ctx) }()
	ctx.Resp, ctx.Err = resp.SetData(map[string]interface{}{"hello": "kongming"}).SetMessage("say hello"), nil
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

```


### kubego

kubego 是一个kubernetes的client-go的封装，方便使用kubernetes的api,支持多种资源的操作，例如pod、deployment、service等。

其中 getter 是对资源的获取，updater 是对资源的更新，删除和创建

informer 是对资源的监听，watcher 是对资源的watch。


### rollout

rollout 是一个面向yaml文件的kubernetes资源的操作，例如kubectl apply 能力。

```
import (
	"github.com/illuminatingKong/kongming-kit/kubego/client"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

func TestDemo(t *testing.T) {
	yamlContent := `apiVersion: v1
kind: Namespace
metadata:
  name: ops
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment-takumi
  namespace: ops
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx
        ports:
        - containerPort: 80`
	d := &Deploy{}
	kubeconfig := []byte(``)
	c, err := client.NewRestClient(kubeconfig)
	if err != nil {
		t.Fatal(err)
	}

	cl, err := client.NewClent(c, runtimeclient.Options{})
	if err != nil {
		t.Fatal(err)
	}
	createErr := d.CreateOrPatch(yamlContent, cl, Option{NameSpace: "ops"})
	if createErr != nil {
		panic(createErr)
	}
}


```

监听Pod的状态

```shell
func QueryPodsStatus(informer informers.SharedInformerFactory, label map[string]string) 
```


### 生成UUID

生成24位订单号，uuid 包下Generate 闭包函数可以生成。 


 

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
- [x] sample httpclient 
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

查看例子
```
代码在 examples/example-new-bind/example-new-bing_test.go
```


### kubego

kubego 是一个kubernetes的client-go的封装，方便使用kubernetes的api,支持多种资源的操作，例如pod、deployment、service等。

其中 getter 是对资源的获取，updater 是对资源的更新，删除和创建

informer 是对资源的监听，watcher 是对资源的watch。


### rollout

rollout 是一个面向yaml文件的kubernetes资源的操作，例如kubectl apply 能力。

同时提供一个 informer 的能力，对资源的监听，watch。

发布apply能力，参考  rollout/yaml/apply_test.go

监听Pod的状态

```shell
func QueryPodsStatus(informer informers.SharedInformerFactory, label map[string]string) 
```


### 生成UUID

生成24位订单号，uuid 包下Generate 闭包函数可以生成。 


### http client && sample httpclient

sample httpclient 是基于guzzle.Client的能力,简化操作。支持get、post、put、delete等方法，方便使用。

参考 examples/examplehttpguzzle/examplehttpguzzle_test.go



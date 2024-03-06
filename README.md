# kongming-kit


## 项目介绍

kongming-kit 是一个快速开发工具库，提供了一些常用的工具类和基础功能，方便快速开发。

## 项目结构

`base` 基础的日志、配置、配置目录功能

`http` http 请求客户端、服务端、路由、中间件、响应


## example 案例

代码更目录 example 下面有一些案例，可以参考使用
### 配置文件

example-configx 读取配置文件的

### 创建一个project并启动

example-new-bind 中 TestStartProject 方法创建一个新的项目，并启动

1 创建一个容器，导入配置文件和日志

2 新建bind服务，绑定路由、业务方法和中间件

3 启动bind服务，并且通过启动顺序完成组件初始化。例如 配置文件刷新，依赖的中间件，优雅下线等，使用先进后出加载顺序。

### 创建一个容器

example-new-container 创建一个容器，并且导入配置文件和日志



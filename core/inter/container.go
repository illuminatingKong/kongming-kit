package inter

type Bindings map[string]interface{}

type Container interface {
	BaseContainer
}

type BaseContainer interface {
	AppReader

	//Singleton 方法确保每次从容器中解析时都返回同一个实例，实例只会被创建一次
	Singleton(abstract interface{}, concrete interface{})

	//Bind 用于将一个接口或抽象类绑定到一个具体的实现类。
	//每次从容器中解析该接口或抽象类时，都会创建一个新的实例。
	//这意味着，每次解析都会创建一个新的实例。
	Bind(abstract interface{}, concrete interface{})

	//Instance 绑定到一个已经实现的类或抽象对象，容器会返回你注册的那个具体实例，而不是创建一个新的实例
	Instance(concrete interface{}) interface{}

	//Bindings 是通过 bind、singleton、instance 等方法注册到容器中的服务绑定关系
	Bindings() Bindings

	//Bound 返回bind 方法实现，表示一种服务绑定关系
	Bound(abstract string) bool

	//Extend 在服务被解析之前动态地修改或装饰，使用 extend 来添加日志记录、缓存功能或修改服务的某些属性
	Extend(abstract interface{}, function func(service interface{}) interface{})
}

type AppReader interface {
	// Make the given type from the container.
	Make(abstract interface{}) interface{}

	// MakeE the given type from the container.
	MakeE(abstract interface{}) (interface{}, error)
}

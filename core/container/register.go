package container

type RegisterServiceProvider interface {
	Register(container Container) Container
}

package container

type BootServiceProvider interface {
	Boot(container Container) Container
}

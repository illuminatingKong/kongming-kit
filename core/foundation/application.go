package foundation

import "github.com/illuminatingKong/kongming-kit/core/inter"

type Application struct {
	// The service container
	container *inter.Container
}

func NewApp() *Application {
	app := Application{}

	container := NewContainer()
	app.SetContainer(container)

	return &app
}

func NewContainer() *Container {
	containerStruct := Container{}
	containerStruct.bindings = make(inter.Bindings)
	containerStruct.singletons = make(inter.Bindings)

	return &containerStruct
}

// SetContainer set the service container
func (a *Application) SetContainer(container inter.Container) {
	a.container = &container
}

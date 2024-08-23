package foundation

import (
	"fmt"
	"github.com/illuminatingKong/kongming-kit/core/inter"
	"github.com/illuminatingKong/kongming-kit/core/support"
	"github.com/pkg/errors"
	"reflect"
)

type Container struct {
	// The container created at boot time
	bootContainer inter.Container

	// A key value store. Callbacks are not automatically executed and stored.
	bindings inter.Bindings

	// A key value store. If the value is a callback, it will be executed
	// once per request and the result will be saved.
	singletons inter.Bindings
}

func (c *Container) getConcreteBinding(
	concrete interface{},
	object interface{},
	abstractName string,
) (interface{}, error) {
	// If abstract is bound, use that object.
	concrete = object
	value := reflect.ValueOf(concrete)

	// If concrete is a callback, run it and save the result.
	if value.Kind() == reflect.Func {
		if value.Type().NumIn() != 0 {
			return nil, support.WithStack(support.CanNotInstantiateCallbackWithParameters)
		}
		concrete = value.Call([]reflect.Value{})[0].Interface()
	}

	// Don't save result in bootContainer. We don't want to share the result across multiple requests
	if c.bootContainer != nil {
		c.bindings[abstractName] = concrete
	}

	return concrete, nil
}

func (c Container) Make(abstract interface{}) interface{} {
	//TODO implement me
	panic("implement me")
}

func (c Container) MakeE(abstract interface{}) (interface{}, error) {
	var concrete interface{}
	var err error = nil
	var abstractName = support.Name(abstract)
	kind := support.Kind(abstract)
	if kind == reflect.Ptr && abstract == nil {
		return nil, errors.New("can't resolve interface. To resolve an interface, " +
			"use the following syntax: (*interface)(nil), use a string or use the struct itself")
	}

	if object, present := c.bindings[abstractName]; present {
		concrete = object

	} else if object, present := c.singletons[abstractName]; present {
		concrete, err = c.getConcreteBinding(concrete, object, abstractName)

	} else if c.bootContainer != nil && c.bootContainer.Bound(abstractName) {
		// Check the container that was created at boot time
		concrete, err = c.bootContainer.MakeE(abstract)
		c.bindings[abstractName] = concrete

	} else if kind == reflect.Struct {
		// If struct cannot be found, we simply have to use the struct itself.
		concrete = abstract
	} else if kind == reflect.String {
		var instances support.Map
		instances, err = support.NewMapE(c.bindings)
		if err == nil {
			var value support.Value
			if c.bootContainer != nil {
				bootBindings := c.bootContainer.Bindings()
				bootInstances, err := support.NewMapE(bootBindings)
				if err != nil {
					return nil, err
				}
				instances.Merge(bootInstances)
			}
			value, err = instances.GetE(abstract.(string))
			//goland:noinspection GoNilness
			concrete = value.Raw()
		}
	}

	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("get instance '%s' from container", abstractName))
	}

	resolvePointerValue(abstract, concrete)

	return concrete, err
}

func (c Container) Singleton(abstract interface{}, concrete interface{}) {
	//TODO implement me
	panic("implement me")
}

func (c Container) Bind(abstract interface{}, concrete interface{}) {
	//TODO implement me
	panic("implement me")
}

func (c Container) Instance(concrete interface{}) interface{} {
	//TODO implement me
	panic("implement me")
}

func (c Container) Bindings() inter.Bindings {
	//TODO implement me
	panic("implement me")
}

func (c Container) Bound(abstract string) bool {
	//TODO implement me
	panic("implement me")
}

func (c Container) Extend(abstract interface{}, function func(service interface{}) interface{}) {
	//TODO implement me
	panic("implement me")
}

func resolvePointerValue(abstract interface{}, concrete interface{}) {
	if support.Kind(abstract) == reflect.Ptr {
		of := reflect.ValueOf(abstract)
		if !of.IsNil() {
			of.Elem().Set(reflect.ValueOf(concrete))
		}
	}
}

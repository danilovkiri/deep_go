package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type SingletonService struct {
	NotEmptyStruct bool
}

type Container struct {
	constructors map[string]reflect.Value
	singletons   map[string]reflect.Value
}

func NewContainer() *Container {
	return &Container{
		constructors: make(map[string]reflect.Value),
		singletons:   make(map[string]reflect.Value),
	}
}

// RegisterType resolves dependencies lazily upon calling Resolve().
func (c *Container) RegisterType(name string, constructor any) {
	value := reflect.ValueOf(constructor)

	// we would have to check that the constructor is a function with one return value, but we won't
	//valueType := value.Type()
	//if valueType.Kind() != reflect.Func || valueType.NumOut() != 1 {
	//	panic("constructor must be a function with one return value")
	//}

	c.constructors[name] = value
}

// RegisterSingletonType resolves dependencies immediately.
func (c *Container) RegisterSingletonType(id string, constructor any) error {
	value := reflect.ValueOf(constructor)

	// we would have to check that the constructor is a function with one return value, but we won't
	//valueType := value.Type()
	//if valueType.Kind() != reflect.Func || valueType.NumOut() != 1 {
	//	panic("constructor must be a function with one return value")
	//}

	// temporarily register constructor to enable dependency resolution
	c.constructors[id] = value

	inst, err := c.Resolve(id)
	if err != nil {
		return fmt.Errorf("failed to resolve singleton '%s': %w", id, err)
	}

	// save as singleton, remove from constructors
	c.singletons[id] = reflect.ValueOf(inst)
	delete(c.constructors, id)
	return nil
}

func (c *Container) Resolve(name string) (any, error) {
	if inst, ok := c.singletons[name]; ok {
		return inst.Interface(), nil
	}

	ctor, ok := c.constructors[name]
	if !ok {
		return nil, fmt.Errorf("no constructor was registered for id '%s'", name)
	}

	ctorType := ctor.Type()
	args := make([]reflect.Value, ctorType.NumIn())

	for i := 0; i < ctorType.NumIn(); i++ {
		inType := ctorType.In(i)
		inID := inType.String()

		dep, err := c.Resolve(inID)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve dependencies for '%s': %w", inID, err)
		}
		args[i] = reflect.ValueOf(dep)
	}

	results := ctor.Call(args)
	instance := results[0] // enforced by the registering methods

	return instance.Interface(), nil
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() any { return &UserService{} })
	container.RegisterType("MessageService", func() any { return &MessageService{} })
	container.RegisterSingletonType("SingletonService", func() any { return &SingletonService{} })

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)

	singletonService1, err := container.Resolve("SingletonService")
	assert.NoError(t, err)
	singletonService2, err := container.Resolve("SingletonService")
	assert.NoError(t, err)

	s1 := singletonService1.(*SingletonService)
	s2 := singletonService2.(*SingletonService)
	assert.True(t, s1 == s2)
}

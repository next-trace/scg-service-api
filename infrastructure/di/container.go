// Package di provides dependency injection utilities.
//
// Note: This package requires the following dependencies:
// - github.com/uber-go/dig
//
// See docs/dependencies.md for more information.
package di

import (
	"fmt"
	"reflect"
)

// Container is a dependency injection container.
// In a real implementation, this would use the dig library from Uber.
// For now, we'll provide a simple implementation that can be replaced later.
type Container struct {
	providers map[reflect.Type]provider
	instances map[reflect.Type]interface{}
}

type provider struct {
	constructor interface{}
	params      []reflect.Type
}

// NewContainer creates a new dependency injection container.
func NewContainer() *Container {
	return &Container{
		providers: make(map[reflect.Type]provider),
		instances: make(map[reflect.Type]interface{}),
	}
}

// Provide registers a constructor function with the container.
// The constructor function should return a value of the type to be provided.
func (c *Container) Provide(constructor interface{}) error {
	constructorType := reflect.TypeOf(constructor)
	if constructorType.Kind() != reflect.Func {
		return fmt.Errorf("constructor must be a function")
	}

	if constructorType.NumOut() == 0 {
		return fmt.Errorf("constructor must return at least one value")
	}

	// Get the type of the first return value
	returnType := constructorType.Out(0)

	// Get parameter types
	params := make([]reflect.Type, constructorType.NumIn())
	for i := range constructorType.NumIn() {
		params[i] = constructorType.In(i)
	}

	// Register the provider
	c.providers[returnType] = provider{
		constructor: constructor,
		params:      params,
	}

	return nil
}

// Resolve resolves a dependency from the container.
func (c *Container) Resolve(target interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}

	targetElem := targetValue.Elem()
	targetType := targetElem.Type()

	// Check if we already have an instance
	if instance, ok := c.instances[targetType]; ok {
		targetElem.Set(reflect.ValueOf(instance))
		return nil
	}

	// Find the provider
	provider, ok := c.providers[targetType]
	if !ok {
		return fmt.Errorf("no provider found for type %s", targetType)
	}

	// Resolve dependencies
	args := make([]reflect.Value, len(provider.params))
	for i, paramType := range provider.params {
		// Create a new instance of the parameter type
		paramInstance := reflect.New(paramType).Interface()

		// Resolve the parameter
		if err := c.Resolve(paramInstance); err != nil {
			return fmt.Errorf("failed to resolve parameter %d: %w", i, err)
		}

		// Get the value of the parameter
		args[i] = reflect.ValueOf(paramInstance).Elem()
	}

	// Call the constructor
	results := reflect.ValueOf(provider.constructor).Call(args)
	if len(results) == 0 {
		return fmt.Errorf("constructor returned no values")
	}

	// Store the instance
	instance := results[0].Interface()
	c.instances[targetType] = instance

	// Set the target value
	targetElem.Set(reflect.ValueOf(instance))

	return nil
}

// Invoke calls the given function with resolved dependencies.
func (c *Container) Invoke(function interface{}) error {
	functionType := reflect.TypeOf(function)
	if functionType.Kind() != reflect.Func {
		return fmt.Errorf("function must be a function")
	}

	// Resolve dependencies
	args := make([]reflect.Value, functionType.NumIn())
	for i := range functionType.NumIn() {
		paramType := functionType.In(i)
		paramInstance := reflect.New(paramType).Interface()

		// Resolve the parameter
		if err := c.Resolve(paramInstance); err != nil {
			return fmt.Errorf("failed to resolve parameter %d: %w", i, err)
		}

		// Get the value of the parameter
		args[i] = reflect.ValueOf(paramInstance).Elem()
	}

	// Call the function
	reflect.ValueOf(function).Call(args)

	return nil
}

// Reset clears all instances from the container.
func (c *Container) Reset() {
	c.instances = make(map[reflect.Type]interface{})
}

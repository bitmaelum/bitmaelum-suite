// Copyright (c) 2020 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package container

/*
 * This is a very basic container system. It is NOT directly an dependecy container, as it does not do any resolving
 * and dependencies. But it's here to easily change the functionality. This is needed when we want for instance to
 * mock a service. From the code point of view, we still can use container.Get("service"), while we have set this
 * to a mocked service.
 *
 * There is no functionality for dependencies etc
 */

// Container is the main container structure holding all service
type Container struct {
	definitions map[string]*ServiceDefinition
	resolved    map[string]interface{}
}

// The main container instance
var container = Container{
	definitions: make(map[string]*ServiceDefinition),
	resolved:    make(map[string]interface{}),
}

// Get returns a service from the container
func Get(key string) interface{} {
	return container.Get(key)
}

type ServiceFunc func() (interface{}, error)

// Set sets a service from the container as a singleton
func Set(key string, build ServiceFunc) {
	container.Set(key, ServiceTypeSingle, build)
}

// Set sets a service from the container. It will return a new instance on each call
func SetMulti(key string, build ServiceFunc) {
	container.Set(key, ServiceTypeMulti, build)
}

// ServiceType defines what kind of service it is (singleton, or new instance on each call)
type ServiceType int

const (
	ServiceTypeSingle ServiceType = iota
	ServiceTypeMulti
)

// Service is a single service
type ServiceDefinition struct {
	Func ServiceFunc
	Type ServiceType
}

// NewContainer will create a new container
func NewContainer() Container {
	c := Container{
		definitions: make(map[string]*ServiceDefinition),
		resolved:    make(map[string]interface{}),
	}

	return c
}

// Set will set the function for the given service
func (c Container) Set(key string, t ServiceType, f ServiceFunc) {
	c.definitions[key] = &ServiceDefinition{
		Func: f,
		Type: t,
	}

	// Delete existing resolved object if any
	delete(c.resolved, key)
}

// Get will retrieve the function for the given service, or nil when not found
func (c Container) Get(key string) interface{} {
	s, ok := c.definitions[key]
	if !ok {
		return nil
	}

	// Multi means we don't use a shared instance but instead instantiate a new object each time called
	if s.Type == ServiceTypeMulti {
		obj, err := s.Func()
		if err != nil {
			return nil
		}
		return obj
	}

	// Already resolved, return
	if c.resolved[key] != nil {
		return c.resolved[key]
	}

	// Create instance, save it and return
	obj, err := s.Func()
	if err != nil {
		return nil
	}
	c.resolved[key] = obj

	return obj
}

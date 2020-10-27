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
 * mock a service. From the code point of view, we still can use  GetContainer().Get("service"), while we have set this
 * to a mocked service.
 *
 * There is no functionality for dependencies, singletons etc
 */

// Container is the main container structure holding all service
type Container struct {
	services map[string]Service
}

// The main container instance
var container = Container{
	services: make(map[string]Service),
}

// Service is a single service
type Service struct {
	Func interface{} // Function that resolves for this service
}

// NewContainer will create a new container
func NewContainer() Container {
	c := Container{
		services: make(map[string]Service),
	}

	return c
}

// GetContainer will return the current initialized container
func GetContainer() Container {
	return container
}

// Set will set the function for the given service
func (c Container) Set(key string, f interface{}) {
	s := Service{
		Func: f,
	}
	c.services[key] = s
}

// Get will retrieve the function for the given service, or nil when not found
func (c Container) Get(key string) interface{} {
	s, ok := c.services[key]
	if !ok {
		return nil
	}

	return s.Func
}

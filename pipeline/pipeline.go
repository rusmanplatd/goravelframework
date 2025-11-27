package pipeline

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/goravel/framework/contracts/foundation"
	contractspipeline "github.com/goravel/framework/contracts/pipeline"
)

// Pipeline implements the pipeline pattern for passing data through stages.
type Pipeline struct {
	container foundation.Application
	passable  any
	pipes     []any
	method    string
}

// NewPipeline creates a new Pipeline instance.
func NewPipeline(container foundation.Application) contractspipeline.Pipeline {
	return &Pipeline{
		container: container,
		method:    "Handle", // Default method name
		pipes:     make([]any, 0),
	}
}

// Send sets the object being sent through the pipeline.
func (p *Pipeline) Send(passable any) contractspipeline.Pipeline {
	p.passable = passable
	return p
}

// Through sets the array of pipes.
func (p *Pipeline) Through(pipes ...any) contractspipeline.Pipeline {
	p.pipes = pipes
	return p
}

// Pipe pushes additional pipes onto the pipeline.
func (p *Pipeline) Pipe(pipes ...any) contractspipeline.Pipeline {
	p.pipes = append(p.pipes, pipes...)
	return p
}

// Via sets the method to call on the pipes.
func (p *Pipeline) Via(method string) contractspipeline.Pipeline {
	p.method = method
	return p
}

// Then runs the pipeline with a final destination callback.
func (p *Pipeline) Then(destination func(any) any) any {
	pipeline := p.prepareDestination(destination)

	// Build the pipeline by reducing pipes in reverse order
	for i := len(p.pipes) - 1; i >= 0; i-- {
		pipeline = p.carry(pipeline, p.pipes[i])
	}

	return pipeline(p.passable)
}

// ThenReturn runs the pipeline and returns the result.
func (p *Pipeline) ThenReturn() any {
	return p.Then(func(passable any) any {
		return passable
	})
}

// prepareDestination wraps the final destination in a closure.
func (p *Pipeline) prepareDestination(destination func(any) any) func(any) any {
	return func(passable any) any {
		return destination(passable)
	}
}

// carry creates a closure that represents a slice of the pipeline onion.
func (p *Pipeline) carry(stack func(any) any, pipe any) func(any) any {
	return func(passable any) any {
		// Handle callable pipes (functions)
		if fn, ok := pipe.(func(any, func(any) any) any); ok {
			return fn(passable, stack)
		}

		// Handle Pipe interface
		if pipeObj, ok := pipe.(contractspipeline.Pipe); ok {
			return pipeObj.Handle(passable, stack)
		}

		// Handle string pipes (class names to be resolved from container)
		if pipeStr, ok := pipe.(string); ok {
			return p.handleStringPipe(pipeStr, passable, stack)
		}

		// Handle object pipes with custom methods
		return p.handleObjectPipe(pipe, passable, stack)
	}
}

// handleStringPipe resolves a pipe from the container and executes it.
func (p *Pipeline) handleStringPipe(pipeStr string, passable any, stack func(any) any) any {
	name, parameters := p.parsePipeString(pipeStr)

	// Resolve the pipe from the container
	instance, err := p.container.Make(name)
	if err != nil {
		panic(fmt.Sprintf("Pipeline: unable to resolve pipe '%s': %v", name, err))
	}

	// Build parameters array
	params := []any{passable, stack}
	params = append(params, parameters...)

	return p.callPipeMethod(instance, params)
}

// handleObjectPipe calls the method on an object pipe.
func (p *Pipeline) handleObjectPipe(pipe any, passable any, stack func(any) any) any {
	params := []any{passable, stack}
	return p.callPipeMethod(pipe, params)
}

// callPipeMethod calls the configured method on a pipe object using reflection.
func (p *Pipeline) callPipeMethod(pipe any, params []any) any {
	pipeValue := reflect.ValueOf(pipe)
	pipeType := pipeValue.Type()

	// Check if the method exists
	method := pipeValue.MethodByName(p.method)
	if !method.IsValid() {
		// If method doesn't exist, try calling the pipe as a function
		if pipeType.Kind() == reflect.Func {
			return p.callFunction(pipeValue, params)
		}
		panic(fmt.Sprintf("Pipeline: method '%s' does not exist on pipe %T", p.method, pipe))
	}

	return p.callFunction(method, params)
}

// callFunction calls a function using reflection with the given parameters.
func (p *Pipeline) callFunction(fn reflect.Value, params []any) any {
	// Convert params to reflect.Value
	in := make([]reflect.Value, len(params))
	for i, param := range params {
		in[i] = reflect.ValueOf(param)
	}

	// Call the function
	results := fn.Call(in)

	if len(results) == 0 {
		return nil
	}

	return results[0].Interface()
}

// parsePipeString parses a pipe string in the format "ClassName:param1,param2".
func (p *Pipeline) parsePipeString(pipe string) (string, []any) {
	parts := strings.SplitN(pipe, ":", 2)
	name := parts[0]

	if len(parts) == 1 {
		return name, []any{}
	}

	// Parse parameters
	paramStrs := strings.Split(parts[1], ",")
	parameters := make([]any, len(paramStrs))
	for i, param := range paramStrs {
		parameters[i] = strings.TrimSpace(param)
	}

	return name, parameters
}

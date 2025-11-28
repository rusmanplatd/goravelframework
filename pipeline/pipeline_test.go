package pipeline

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rusmanplatd/goravelframework/contracts/foundation"
)

// TestPipeline_BasicFlow tests basic pipeline execution with a single pipe
func TestPipeline_BasicFlow(t *testing.T) {
	pipeline := NewPipeline(nil)

	result := pipeline.
		Send("hello").
		Through(func(passable any, next func(any) any) any {
			str := passable.(string)
			return next(str + " world")
		}).
		Then(func(passable any) any {
			return passable
		})

	assert.Equal(t, "hello world", result)
}

// TestPipeline_MultiplePipes tests pipeline with multiple pipes
func TestPipeline_MultiplePipes(t *testing.T) {
	pipeline := NewPipeline(nil)

	result := pipeline.
		Send("hello").
		Through(
			func(passable any, next func(any) any) any {
				str := passable.(string)
				return next(str + " world")
			},
			func(passable any, next func(any) any) any {
				str := passable.(string)
				return next(strings.ToUpper(str))
			},
			func(passable any, next func(any) any) any {
				str := passable.(string)
				return next(str + "!")
			},
		).
		Then(func(passable any) any {
			return passable
		})

	assert.Equal(t, "HELLO WORLD!", result)
}

// TestPipeline_ThenReturn tests the ThenReturn method
func TestPipeline_ThenReturn(t *testing.T) {
	pipeline := NewPipeline(nil)

	result := pipeline.
		Send(42).
		Through(
			func(passable any, next func(any) any) any {
				num := passable.(int)
				return next(num * 2)
			},
			func(passable any, next func(any) any) any {
				num := passable.(int)
				return next(num + 10)
			},
		).
		ThenReturn()

	assert.Equal(t, 94, result) // (42 * 2) + 10 = 94
}

// TestPipeline_EmptyPipes tests pipeline with no pipes
func TestPipeline_EmptyPipes(t *testing.T) {
	pipeline := NewPipeline(nil)

	result := pipeline.
		Send("unchanged").
		Through().
		ThenReturn()

	assert.Equal(t, "unchanged", result)
}

// TestPipeline_PipeMethod tests adding pipes with the Pipe method
func TestPipeline_PipeMethod(t *testing.T) {
	pipeline := NewPipeline(nil)

	result := pipeline.
		Send(10).
		Through(func(passable any, next func(any) any) any {
			num := passable.(int)
			return next(num + 5)
		}).
		Pipe(func(passable any, next func(any) any) any {
			num := passable.(int)
			return next(num * 2)
		}).
		ThenReturn()

	assert.Equal(t, 30, result) // (10 + 5) * 2 = 30
}

// TestPipeline_ViaMethod tests custom method names
func TestPipeline_ViaMethod(t *testing.T) {
	pipeline := NewPipeline(nil)

	pipe := &CustomPipe{}

	result := pipeline.
		Send("test").
		Via("Process").
		Through(pipe).
		ThenReturn()

	assert.Equal(t, "PROCESSED: test", result)
}

// TestPipeline_PipeInterface tests pipes implementing the Pipe interface
func TestPipeline_PipeInterface(t *testing.T) {
	pipeline := NewPipeline(nil)

	pipe1 := &UppercasePipe{}
	pipe2 := &ExclamationPipe{}

	result := pipeline.
		Send("hello").
		Through(pipe1, pipe2).
		ThenReturn()

	assert.Equal(t, "HELLO!", result)
}

// TestPipeline_MixedPipes tests mixing different pipe types
func TestPipeline_MixedPipes(t *testing.T) {
	pipeline := NewPipeline(nil)

	result := pipeline.
		Send("start").
		Through(
			func(passable any, next func(any) any) any {
				return next(passable.(string) + " ")
			},
			&UppercasePipe{},
			func(passable any, next func(any) any) any {
				return next(passable.(string) + "!")
			},
		).
		ThenReturn()

	assert.Equal(t, "START !", result)
}

// TestPipeline_ComplexData tests pipeline with complex data types
func TestPipeline_ComplexData(t *testing.T) {
	type Request struct {
		Value int
	}

	pipeline := NewPipeline(nil)

	result := pipeline.
		Send(&Request{Value: 10}).
		Through(
			func(passable any, next func(any) any) any {
				req := passable.(*Request)
				req.Value += 5
				return next(req)
			},
			func(passable any, next func(any) any) any {
				req := passable.(*Request)
				req.Value *= 2
				return next(req)
			},
		).
		ThenReturn()

	assert.Equal(t, 30, result.(*Request).Value)
}

// TestPipeline_Chaining tests method chaining
func TestPipeline_Chaining(t *testing.T) {
	pipeline := NewPipeline(nil)

	result := pipeline.
		Send(1).
		Through(func(passable any, next func(any) any) any {
			return next(passable.(int) + 1)
		}).
		Pipe(func(passable any, next func(any) any) any {
			return next(passable.(int) * 2)
		}).
		Pipe(func(passable any, next func(any) any) any {
			return next(passable.(int) + 10)
		}).
		ThenReturn()

	assert.Equal(t, 14, result) // ((1 + 1) * 2) + 10 = 14
}

// Helper pipes for testing

// UppercasePipe implements the Pipe interface
type UppercasePipe struct{}

func (p *UppercasePipe) Handle(passable any, next func(any) any) any {
	str := passable.(string)
	return next(strings.ToUpper(str))
}

// ExclamationPipe implements the Pipe interface
type ExclamationPipe struct{}

func (p *ExclamationPipe) Handle(passable any, next func(any) any) any {
	str := passable.(string)
	return next(str + "!")
}

// CustomPipe for testing Via method
type CustomPipe struct{}

func (p *CustomPipe) Process(passable any, next func(any) any) any {
	str := passable.(string)
	return next("PROCESSED: " + str)
}

// MockContainer for testing string-based pipe resolution
type MockContainer struct {
	foundation.Application
	bindings map[string]any
}

func NewMockContainer() *MockContainer {
	return &MockContainer{
		bindings: make(map[string]any),
	}
}

func (m *MockContainer) Make(key any) (any, error) {
	if binding, ok := m.bindings[key.(string)]; ok {
		return binding, nil
	}
	return nil, nil
}

func (m *MockContainer) Register(key string, binding any) {
	m.bindings[key] = binding
}

// TestPipeline_StringPipes tests string-based pipe resolution
func TestPipeline_StringPipes(t *testing.T) {
	container := NewMockContainer()
	container.Register("uppercase", &UppercasePipe{})

	pipeline := NewPipeline(container)

	result := pipeline.
		Send("hello").
		Through("uppercase").
		ThenReturn()

	assert.Equal(t, "HELLO", result)
}

// TestPipeline_StringPipesWithParameters tests string pipes with parameters
func TestPipeline_StringPipesWithParameters(t *testing.T) {
	container := NewMockContainer()

	// Create a pipe that uses parameters
	paramPipe := &ParameterizedPipe{}
	container.Register("param_pipe", paramPipe)

	pipeline := NewPipeline(container)

	result := pipeline.
		Send("hello").
		Through("param_pipe:prefix,suffix").
		ThenReturn()

	assert.Equal(t, "prefix-hello-suffix", result)
}

// ParameterizedPipe for testing parameters
type ParameterizedPipe struct{}

func (p *ParameterizedPipe) Handle(passable any, next func(any) any, params ...any) any {
	str := passable.(string)
	if len(params) >= 2 {
		prefix := params[0].(string)
		suffix := params[1].(string)
		str = prefix + "-" + str + "-" + suffix
	}
	return next(str)
}

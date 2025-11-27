package pipeline

// Pipeline represents a pipeline for passing data through stages.
// It implements the "onion" pattern where data flows through a series of pipes,
// each having the opportunity to inspect or modify the data before passing it to the next pipe.
type Pipeline interface {
	// Send sets the object being sent through the pipeline.
	Send(passable any) Pipeline

	// Through sets the array of pipes.
	// Pipes can be:
	// - Functions with signature: func(passable any, next func(any) any) any
	// - Objects implementing the Pipe interface
	// - Strings representing class names (resolved from container)
	Through(pipes ...any) Pipeline

	// Pipe pushes additional pipes onto the pipeline.
	Pipe(pipes ...any) Pipeline

	// Via sets the method to call on the pipes.
	// Default is "Handle". This allows using custom method names on pipe objects.
	Via(method string) Pipeline

	// Then runs the pipeline with a final destination callback.
	// The destination receives the final processed value.
	Then(destination func(any) any) any

	// ThenReturn runs the pipeline and returns the result without a destination callback.
	// Equivalent to Then(func(passable any) any { return passable })
	ThenReturn() any
}

// Pipe represents a single stage in the pipeline.
// Pipes should implement this interface to be used in the pipeline.
type Pipe interface {
	// Handle processes the passable and calls the next pipe in the chain.
	// The next parameter is a function that should be called to pass control to the next pipe.
	Handle(passable any, next func(any) any) any
}

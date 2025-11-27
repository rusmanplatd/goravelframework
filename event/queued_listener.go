package event

import (
	"fmt"
	"reflect"
	"time"

	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/contracts/queue"
)

// QueuedListenerJob represents a job that executes a queued listener.
type QueuedListenerJob struct {
	listener  any
	eventName string
	payload   []any
}

// NewQueuedListenerJob creates a new queued listener job.
func NewQueuedListenerJob(listener any, eventName string, payload []any) *QueuedListenerJob {
	return &QueuedListenerJob{
		listener:  listener,
		eventName: eventName,
		payload:   payload,
	}
}

// Signature returns the unique identifier for the job.
func (j *QueuedListenerJob) Signature() string {
	return fmt.Sprintf("queued_listener:%s", j.eventName)
}

// Handle executes the queued listener.
func (j *QueuedListenerJob) Handle(args ...any) error {
	// Invoke the listener with the event and payload
	_, err := invokeListener(j.listener, j.eventName, j.payload)
	return err
}

// shouldQueueListener checks if a listener should be queued.
// Supports both ShouldQueue interface and shouldQueue(event) method.
func shouldQueueListener(listener any, payload []any) bool {
	// First check ShouldQueue interface
	if queueable, ok := listener.(event.ShouldQueue); ok {
		return queueable.ShouldQueue()
	}

	// Then check for shouldQueue method with event parameter (Laravel compatibility)
	v := reflect.ValueOf(listener)
	method := v.MethodByName("ShouldQueue")
	if method.IsValid() && method.Kind() == reflect.Func {
		// Call method with first payload item as event
		var args []reflect.Value
		if len(payload) > 0 {
			args = append(args, reflect.ValueOf(payload[0]))
		}

		results := method.Call(args)
		if len(results) > 0 && results[0].Kind() == reflect.Bool {
			return results[0].Bool()
		}
	}

	return false
}

// getListenerQueueConfig extracts queue configuration from a listener.
// Checks interface methods first, then falls back to struct fields (matching Laravel behavior).
func getListenerQueueConfig(listener any) (connection, queueName string, delay time.Duration) {
	// First, try interface methods (preferred)
	if queueable, ok := listener.(event.QueueableListener); ok {
		connection = queueable.ViaConnection()
		queueName = queueable.ViaQueue()
		// Convert seconds to duration
		delay = time.Duration(queueable.WithDelay()) * time.Second
		return
	}

	// Fall back to struct fields (Laravel compatibility)
	v := reflect.ValueOf(listener)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		// Check for Connection field
		if field := v.FieldByName("Connection"); field.IsValid() && field.CanInterface() {
			if connStr, ok := field.Interface().(string); ok {
				connection = connStr
			}
		}

		// Check for Queue field
		if field := v.FieldByName("Queue"); field.IsValid() && field.CanInterface() {
			if queueStr, ok := field.Interface().(string); ok {
				queueName = queueStr
			}
		}

		// Check for Delay field (in seconds as int)
		if field := v.FieldByName("Delay"); field.IsValid() && field.CanInterface() {
			if delayInt, ok := field.Interface().(int); ok {
				delay = time.Duration(delayInt) * time.Second
			}
		}
	}

	return
}

// queueListener queues a listener for execution.
func queueListener(queueInstance queue.Queue, listener any, eventName string, payload []any) error {
	// Create the job
	job := NewQueuedListenerJob(listener, eventName, payload)

	// Get queue configuration
	connection, queueName, delay := getListenerQueueConfig(listener)

	// Build the queue job
	queueJob := queueInstance.Job(job, []queue.Arg{})

	// Apply connection if specified
	if connection != "" {
		queueJob = queueJob.OnConnection(connection)
	}

	// Apply queue name if specified
	if queueName != "" {
		queueJob = queueJob.OnQueue(queueName)
	}

	// Apply delay if specified
	if delay > 0 {
		queueJob = queueJob.Delay(time.Now().Add(delay))
	}

	// Dispatch the job
	return queueJob.Dispatch()
}

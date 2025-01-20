package eventually

import (
	"fmt"
	"reflect"
)

// EventHandle is a function for event handling.
//
// The form of function is:
//
//	func(event Event)
//
// where the [Event] represent the event (struct).
//
// Event handling never fails, no error is returned.
//
// Example:
//
//	 pubMux := &eventually.PubMux{}
//	 pubMux.React(func(event OrderCompleted) {
//		    // handle the event
//	 })
type EventHandler any

// PubMux is a multiplexer for events. It will route events to its handlers.
type PubMux struct {
	handlers map[reflect.Type][]EventHandler
}

// NewPubMux allocates and returns a new PubMux.
func NewPubMux() *PubMux {
	return &PubMux{}
}

// React will handle the published event.
//
// It will handle only the event type defined by the fn.
// The fn needs to be a valid [EventHandler], otherwise it will panic.
// Multiple handlers can be registered for the same event type.
func (pm *PubMux) React(fn EventHandler) {
	if err := validateHandler(fn); err != nil {
		panic(err)
	}

	if pm.handlers == nil {
		pm.handlers = make(map[reflect.Type][]EventHandler)
	}

	fnType := reflect.TypeOf(fn)
	fnTypeIn := fnType.In(0)
	pm.handlers[fnTypeIn] = append(pm.handlers[fnTypeIn], fn)
}

// Publish the event.
//
// The event needs to be a valid [Event], otherwise it will panic.
func (pm *PubMux) Publish(event Event) {
	eventType := reflect.TypeOf(event)
	if eventType.Kind() != reflect.Struct {
		panic(fmt.Sprintf("eventually: event should be a struct (got: %v)", eventType.Kind()))
	}

	if pm.handlers == nil {
		return
	}

	handlers := pm.handlers[eventType]
	for _, handler := range handlers {
		invokeHandler(handler, event)
	}
}

func invokeHandler(handler EventHandler, event Event) {
	fnValue := reflect.ValueOf(handler)
	eventValue := reflect.ValueOf(event)
	fnValue.Call([]reflect.Value{eventValue})
}

func validateHandler(fn EventHandler) error {
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("eventually: fn EventHandler is not a function (got: %v)", fnType.Kind())
	}

	if fnType.NumIn() != 1 {
		return fmt.Errorf("eventually: fn EventHandler should have 1 input parameter (got: %d)", fnType.NumIn())
	}

	if fnType.NumOut() != 0 {
		return fmt.Errorf("eventually: fn EventHandler should have 0 output parameter (got: %d)", fnType.NumOut())
	}

	if fnType.In(0).Kind() != reflect.Struct {
		return fmt.Errorf("eventually: fn EventHandler input parameter should be a struct (got: %v)", fnType.In(0).Kind())
	}

	return nil
}

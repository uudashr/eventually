package eventually

import (
	"fmt"
	"reflect"
)

// EventHandle is a function type that handle the event.
//
// The form of function is:
//
//	func(event Event)
//
// where the [Event] is the event type (struct) that will be handled.
//
// Example:
//
//	 pubMux := &eventually.PubMux{}
//	 pubMux.React(func(event OrderCompleted) {
//		    // handle the event
//	 })
type EventHandler any

// PubMux is a multiplexer for events. It will route events to it's handler.
type PubMux struct {
	handlers map[reflect.Type][]EventHandler
}

// React to an event with given fn as it's handler.
func (pm *PubMux) React(fn EventHandler) error {
	if err := validateHandler(fn); err != nil {
		return err
	}

	if pm.handlers == nil {
		pm.handlers = make(map[reflect.Type][]EventHandler)
	}

	fnType := reflect.TypeOf(fn)
	fnTypeIn := fnType.In(0)
	pm.handlers[fnTypeIn] = append(pm.handlers[fnTypeIn], fn)

	return nil
}

// Publish the event. Publish will only fail if the event is not a struct.
func (pm *PubMux) Publish(event Event) error {
	eventType := reflect.TypeOf(event)
	if eventType.Kind() != reflect.Struct {
		return fmt.Errorf("eventually: event should be a struct (got: %v)", eventType.Kind())
	}

	handlers := pm.handlers[eventType]
	for _, handler := range handlers {
		invokeHandler(handler, event)
	}

	return nil
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

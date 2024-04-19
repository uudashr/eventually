// Package eventually provides a flexible event handling framework that enables
// applications to react to various types of events dynamically. The core
// functionality revolves around defining, raising, and handling events that are
// structured as Go structs.
//
// The [Eventually] acts as the central hub for event management, allowing
// events to be raised, handlers to be registered or removed, and events to be
// dispatched to the appropriate handlers based on the event type. This package
// leverages Go's reflection capabilities to dynamically handle events according
// to their type at runtime.
//
// # Event
//
// An event in this context is defined as any Go struct. The name and fields of
// the struct describe the event. Events are used to carry data relevant to each
// occurrence that handlers might react to.
//
// Events are defined as structs. For example, an event to represent a completed
// order might look like:
//
//	type OrderCompleted struct {
//	  OrderID string
//	}
//
// # EventHandler
//
// An event handler is a function designed to respond to an [Event]. Handlers
// must accept a single parameter that is a struct type representing the event
// they intend to handle.
//
// Event handlers are functions that take a specific event struct as an
// argument. Here is how you might handle the OrderCompleted event:
//
//	func handleOrderCompleted(event OrderCompleted) {
//	  fmt.Println("Order completed:", event.OrderID)
//	}
//
// # Eventually
//
// The central struct in the package, responsible for managing the lifecycle of
// events and their handlers. It provides methods to raise events, add or remove
// event handlers, and process events through their respective handlers.
//
// Event handlers are registered to an instance of [Eventually] using the
// [Eventually.React] method. This setup allows the [Eventually] instance to
// call the handler when the corresponding event type is raised.
//
// Events are raised using the [Eventually.RaiseEvent] method on an [Eventually]
// instance. When an event is raised, all handlers registered for that event
// type are invoked with the event as a parameter.
//
// The package also provides functionality to store and retrieve an [Eventually]
// instance from a context.Context object, facilitating the passing of the event
// handling framework through request or operation contexts in more complex
// applications.
//
// This package is ideal for applications that require a modular and dynamic
// approach to handling events where the types of events are not known at
// compile time.
//
// Example:
//
//	func main() {
//	  evtl := &eventually.Eventually{}
//
//	  // Register the event handler
//	  evtl.React(handleOrderCompleted)
//
//	  // Raise an event
//	  evtl.RaiseEvent(OrderCompleted{OrderID: "1234"})
//	}
package eventually

// Package eventually provides a flexible event handling library that enables
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
// # PubMux
//
// The central struct in the package, responsible for managing the lifecycle of
// events and their handlers. It provides methods to raise events, add or remove
// event handlers, and process events through their respective handlers.
//
// Event handlers are registered to an instance of [PubMux] using the
// [PubMux.React] method. This setup allows the [PubMux] instance to
// call the handler when the corresponding event type is raised.
//
// Events are publish using the [PubMux.Publish]. When an event is publish,
// all registered handlers for that event will be invoked with the event as the
// parameter.
//
// Example:
//
//	func main() {
//	  pubMux := &eventually.PubMux{}
//
//	  // Register the event handler
//	  pubMux.React(handleOrderCompleted)
//
//	  // Raise an event
//	  evtl.RaiseEvent(OrderCompleted{OrderID: "1234"})
//	}
package eventually

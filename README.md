[![Go Reference](https://pkg.go.dev/badge/github.com/uudashr/eventually.svg)](https://pkg.go.dev/github.com/uudashr/eventually)

# Eventually
Eventually provides a flexible event handling library that enables applications to react to various types of events dynamically. The core functionality revolves around defining, raising, and handling events that are structured as Go structs.

## Usage

### Setup Publisher

```go
var pub eventually.Publisher // either eventually.PubMux or eventually.Recorder
ctx = eventually.ContextWithPub(context.Background(), pub)
```

We put `Publisher` inside the `context.Context` so it can be `Publish` anywhere as long as the context pass through the caller.

### Define Event
```go
type OrderCompleted struct {
    OrderId string
}
```

`Event` was defined as struct. The struct can have any fields that represent the event data.

### Publish Event

```go
eventually.Publish(ctx, OrderCompleted{
    OrderID: "123",
})
```

Publish event using `Publisher` available in the context. If there is no `Publisher` in the context, it will do nothing.

### Handling Events

There are 2 way to hanlding event, is by using `PubMux` or `Recorder`.

#### PubMux
`PubMux` is a multiplexer that can handle event based on it's type.

```go
// Setup PubMux
mux := eventually.NewPubMux()
ctx = eventually.ContextWithPub(context.Background(), mux)

// Listen and handle the event
mux.React(func(e OrderCompleted) {
    fmt.Printf("Order completed: %q\n", e.OrderID)
})

// Publish event
eventually.Publish(ctx, OrderCompleted{
    OrderID: "123",
})
```

It will handle all events that match the `OrderCompleted` type.

#### Recorder
`Recorder` records the published events.

```go
// Setup Recorder
var rec := eventually.NewRecorder()
ctx := eventually.ContextWithPub(context.Background(), rec)

// Publish event
eventually.Publish(ctx, OrderCompleted{
    OrderID: "123",
})

// Print the recorded events
fmt.Println(rec.Events)
```

## Related Links
- [A better domain events pattern](https://lostechies.com/jimmybogard/2014/05/13/a-better-domain-events-pattern/)
- [Domain events: Design and implementation](https://learn.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/domain-events-design-implementation)
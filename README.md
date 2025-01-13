[![Go Reference](https://pkg.go.dev/badge/github.com/uudashr/eventually.svg)](https://pkg.go.dev/github.com/uudashr/eventually)

# Eventually
Eventually provides a flexible event handling library that enables applications to react to various types of events dynamically. The core functionality revolves around defining, raising, and handling events that are structured as Go structs.

## Usage


### Setup Eventually

```go
ctx := context.TODO()

var pub eventually.Publisher // either eventually.PubMux or eventually.Recorder
ctx = eventually.PubWithContext(ctx, pub)
```

Eventually setup and put inside the context.Context.

### Define Event
```go
type OrderCompleted struct {
    OrderId string
}
```

Event was defined as struct. The struct can have any fields that represent the event data.

### Raising an Event

```go
eventually.Publish(ctx, OrderCompleted{
    OrderID: "123",
})
```

Raising event using eventually available in the context. It does nothing if the `Eventually` is not available in the context.

## Handling Events

```go
pubMux := &eventually.PubMux{}
pubMux.React(func(e OrderCompleted) {
    fmt.Printf("Order completed: %q\n", e.OrderID)
})
```

It will handle all events that match the type `OrderCompleted`.

## Related Links
- [A better domain events pattern](https://lostechies.com/jimmybogard/2014/05/13/a-better-domain-events-pattern/)
- [Domain events: Design and implementation](https://learn.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/domain-events-design-implementation)
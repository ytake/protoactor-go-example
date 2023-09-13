package message

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/protoactor-go-example/routing/slip/value"
)

type Car struct {
	Color             string
	HasNavigation     bool
	HasParkingSensors bool
}

type Order struct {
	Options []value.CarOptions
}

// RouteSlip is a message
type RouteSlip struct {
	RouteSlip []*actor.PID
	Message   interface{}
}

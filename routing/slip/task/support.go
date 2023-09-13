package task

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/protoactor-go-example/routing/slip/message"
)

type RouteSlip struct{}

// SendMessageToNextTask is a method
func (_ *RouteSlip) SendMessageToNextTask(context actor.Context, routeSlip []*actor.PID, msg interface{}) {
	nextTask := routeSlip[0]
	newSlip := routeSlip[1:]
	if len(newSlip) == 0 {
		context.Send(nextTask, msg)
		return
	}
	context.Send(nextTask, &message.RouteSlip{
		RouteSlip: newSlip,
		Message:   msg,
	})
}

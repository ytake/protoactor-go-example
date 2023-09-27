package license

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/protoactor-go-example/routing/prefrouting/message"
	"github.com/ytake/protoactor-go-example/routing/prefrouting/process"
)

type GetLicense struct {
	pipe    *actor.PID
	initial time.Duration
	id      *string
}

// NewGetLicense is
func NewGetLicense(initial time.Duration) actor.Actor {
	return &GetLicense{initial: initial}
}

func (state *GetLicense) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.SetService:
		state.id = &msg.ID
		state.initial = msg.ServiceTime
		time.Sleep(100 * time.Millisecond)
	case *message.PerformanceRoutingMessage:
		time.Sleep(state.initial)
		pi := &process.ImageProcessing{}
		li, err := pi.GetLicense(msg.Photo)
		if err != nil {
			context.Poison(context.Self())
		}
		context.Respond(&message.PerformanceRoutingMessage{
			Photo:       msg.Photo,
			License:     &li,
			ProcessedBy: state.id,
		})
	case string:
		if msg == "KillFirst" {
			context.Poison(context.Self())
		}
	}
}

type StopDecider struct{}

func NewStopDirective() *StopDecider {
	return &StopDecider{}
}

func (d *StopDecider) Decider() actor.DeciderFunc {
	return func(reason interface{}) actor.Directive {
		fmt.Println("handling failure for child")
		return actor.StopDirective
	}
}

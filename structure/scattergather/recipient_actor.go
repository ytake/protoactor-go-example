package scattergather

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/protoactor-go-example/structure/message"
)

type RecipientList struct {
	PIDs []*actor.PID
}

func (state *RecipientList) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.PhotoMessage:
		for _, recipient := range state.PIDs {
			context.Send(recipient, msg)
		}
	}
}

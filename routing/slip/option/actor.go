package option

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/protoactor-go-example/routing/slip/message"
	"github.com/ytake/protoactor-go-example/routing/slip/task"
)

type PaintCar struct {
	Color string
	*task.RouteSlip
}

func (state *PaintCar) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.RouteSlip:
		switch car := msg.Message.(type) {
		case *message.Car:
			nc := car
			nc.Color = state.Color
			msg.Message = nc
		}
		state.SendMessageToNextTask(context, msg.RouteSlip, msg.Message)
	}
}

type AddNavigation struct {
	*task.RouteSlip
}

func (state *AddNavigation) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.RouteSlip:
		switch car := msg.Message.(type) {
		case *message.Car:
			nc := car
			nc.HasNavigation = true
			msg.Message = nc
		}
		state.SendMessageToNextTask(context, msg.RouteSlip, msg.Message)
	}
}

type AddParkingSensor struct {
	*task.RouteSlip
}

func (state *AddParkingSensor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.RouteSlip:
		switch car := msg.Message.(type) {
		case *message.Car:
			nc := car
			nc.HasParkingSensors = true
			msg.Message = nc
		}
		state.SendMessageToNextTask(context, msg.RouteSlip, msg.Message)
	}
}

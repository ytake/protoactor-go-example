package router

import (
	"slices"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/protoactor-go-example/routing/slip/message"
	"github.com/ytake/protoactor-go-example/routing/slip/task"
	"github.com/ytake/protoactor-go-example/routing/slip/value"
)

// SlipActor is a sample
type SlipActor struct {
	*task.RouteSlip
}

// Receive is sent messages to be processed from the mailbox associated with the instance of the SlipActor
func (state *SlipActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.RouteSlip:
		state.SendMessageToNextTask(context, msg.RouteSlip, msg.Message)
	}
}

// SlipRouter is a sample
type SlipRouter struct {
	PaintBlack       *actor.PID
	PaintGray        *actor.PID
	AddNavigation    *actor.PID
	AddParkingSensor *actor.PID
	EndStep          *actor.PID
	*task.RouteSlip
}

func (state *SlipRouter) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.Order:
		// 回覧リストを作成する
		routeSlip := state.createRouteSlip(msg.Options)
		// 最初のタスクにメッセージを送信する
		state.SendMessageToNextTask(context, routeSlip, &message.Car{})
	}
}

// createRouteSlip is internal method
// このメソッドは、ルートスリップを作成するための内部メソッドです。
// 車のオプションには、色とナビゲーションと駐車に使うセンサーがオプションとしてあります。
func (state *SlipRouter) createRouteSlip(options []value.CarOptions) []*actor.PID {
	routeSlip := make([]*actor.PID, 0)
	// car needs a color
	containsGray := false
	// 灰色指定がある場合は灰色を塗る
	if slices.Contains(options, value.CarColorGray) {
		containsGray = true
	}
	// 指定がない場合は黒を塗る
	if !containsGray {
		routeSlip = append(routeSlip, state.PaintBlack)
	}
	for _, option := range options {
		switch option {
		case value.CarColorGray:
			routeSlip = append(routeSlip, state.PaintGray)
		case value.Navigation:
			routeSlip = append(routeSlip, state.AddNavigation)
		case value.ParkingSensor:
			routeSlip = append(routeSlip, state.AddParkingSensor)
		}
	}
	return append(routeSlip, state.EndStep)
}

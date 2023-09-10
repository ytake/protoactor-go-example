package scattergather

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/protoactor-go-example/structure/message"
	"github.com/ytake/protoactor-go-example/structure/process"
)

type GetSpeed struct {
	pipe *actor.PID
}

// NewGetSpeedActor is actor
func NewGetSpeedActor(pipe *actor.PID) actor.Actor {
	return &GetSpeed{
		pipe: pipe,
	}
}

func (g *GetSpeed) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.PhotoMessage:
		i := &process.ImageProcessing{}
		speed, _ := i.GetSpeed(msg.Photo)
		msg.Speed = &speed
		context.Send(g.pipe, msg)
	}
}

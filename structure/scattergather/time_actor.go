package scattergather

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/protoactor-go-http/structure/message"
	"github.com/ytake/protoactor-go-http/structure/process"
)

type GetTime struct {
	pipe *actor.PID
}

func NewGetTimeActor(pipe *actor.PID) actor.Actor {
	return &GetTime{
		pipe: pipe,
	}
}

func (g *GetTime) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.PhotoMessage:
		i := &process.ImageProcessing{}
		t, _ := i.GetTime(msg.Photo)
		msg.CreationTime = &t
		context.Send(g.pipe, msg)
	}
}

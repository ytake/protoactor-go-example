package scattergather

import (
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/protoactor-go-http/structure/message"
	"github.com/ytake/protoactor-go-http/structure/process"
)

type GetSpeed struct {
	pipe *actor.PID
}

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

type RecipientList struct {
	PIDs []*actor.PID
}

func (state *RecipientList) Receive(context actor.Context) {
	var msgs []*message.PhotoMessage
	var wg sync.WaitGroup
	switch msg := context.Message().(type) {
	case *message.PhotoMessage:
		for _, recipient := range state.PIDs {
			wg.Add(1)
			f := context.RequestFuture(recipient, msg, 1*time.Second)
			go func() {
				defer wg.Done()
				res, err := f.Result()
				if err != nil {
					return
				}
				switch msg := res.(type) {
				case *message.PhotoMessage:
					msgs = append(msgs, msg)
				}
			}()
		}
		wg.Wait()
		context.Respond(msgs)
	}
}

package scattergather

import (
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/stream"
	"github.com/ytake/protoactor-go-example/structure/message"
)

func stubActorProps(ref *actor.PID) *actor.Props {
	return actor.PropsFromFunc(func(ctx actor.Context) {
		switch msg := ctx.Message().(type) {
		case *message.PhotoMessage:
			ctx.Send(ref, msg)
		default:
		}
	})
}

func TestRecipientList_Receive(t *testing.T) {
	t.Run("scatter the message", func(t *testing.T) {

		system := actor.NewActorSystem()
		p := stream.NewTypedStream[*message.PhotoMessage](system)
		ti := time.Date(2001, 2, 3, 4, 5, 6, 0, time.UTC)
		expectMsg := &message.PhotoMessage{
			ID:    "id1",
			Photo: makeCreatePhotoString(ti, 60)}

		go func() {
			p1 := stubActorProps(p.PID())
			p2 := stubActorProps(p.PID())
			pa1 := system.Root.Spawn(p1)
			pa2 := system.Root.Spawn(p2)
			re := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return &RecipientList{
					PIDs: []*actor.PID{pa1, pa2}}
			}))
			system.Root.Send(re, expectMsg)
		}()
		p1msg := <-p.C()
		if p1msg != expectMsg {
			t.Errorf("expected %v, got %v", expectMsg, p1msg)
		}
		p2msg := <-p.C()
		if p2msg != expectMsg {
			t.Errorf("expected %v, got %v", expectMsg, p2msg)
		}
	})
}

package scattergather

import (
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/protoactor-go-http/structure/message"
)

func stubActorProps() *actor.Props {
	return actor.PropsFromFunc(func(ctx actor.Context) {
		switch msg := ctx.Message().(type) {
		case *message.PhotoMessage:
			ctx.Respond(msg)
		default:
		}
	})
}

func TestRecipientList_Receive(t *testing.T) {
	t.Run("scatter the message", func(t *testing.T) {
		system := actor.NewActorSystem()
		p1 := system.Root.Spawn(stubActorProps())
		p2 := system.Root.Spawn(stubActorProps())
		re := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &RecipientList{
				PIDs: []*actor.PID{p1, p2}}
		}))
		ti := time.Date(2001, 2, 3, 4, 5, 6, 0, time.UTC)
		f := system.Root.RequestFuture(re, &message.PhotoMessage{
			ID:    "id1",
			Photo: makeCreatePhotoString(ti, 60)},
			2*time.Second)

		r, err := f.Result()
		if err != nil {
			t.Errorf("expected %v, got %v", nil, err)
		}
		if len(r.([]*message.PhotoMessage)) != 2 {
			t.Errorf("expected %v, got %v", "testing", len(r.([]*message.PhotoMessage)))
		}
	})
}

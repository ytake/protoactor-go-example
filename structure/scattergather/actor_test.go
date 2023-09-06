package scattergather

import (
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/protoactor-go-http/structure/message"
	"github.com/ytake/protoactor-go-http/structure/process"
)

func TestRecipientList_Receive(t *testing.T) {
	system := actor.NewActorSystem()
	test := actor.PropsFromFunc(func(ctx actor.Context) {
		switch msg := ctx.Message().(type) {
		case *message.PhotoMessage:
			ctx.Respond(msg)
		default:
		}
	})
	p1 := system.Root.Spawn(test)
	p2 := system.Root.Spawn(test)
	re := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return &RecipientList{PIDs: []*actor.PID{p1, p2}}
	}))
	i := &process.ImageProcessing{}
	ti := time.Date(2001, 2, 3, 4, 5, 6, 0, time.UTC)
	f := system.Root.RequestFuture(re, &message.PhotoMessage{
		ID:    "testing",
		Photo: i.CreatePhotoString(ti, 60, "")},
		2*time.Second)
	r, err := f.Result()
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	if len(r.([]*message.PhotoMessage)) != 2 {
		t.Errorf("expected %v, got %v", "testing", len(r.([]*message.PhotoMessage)))
	}
}

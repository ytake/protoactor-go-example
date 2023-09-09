package scattergather

import (
	"reflect"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/stream"
	"github.com/ytake/protoactor-go-http/structure/message"
	"github.com/ytake/protoactor-go-http/structure/process"
)

func makeCreatePhotoString(ti time.Time, s int) string {
	i := &process.ImageProcessing{}
	return i.CreatePhotoString(ti, s, "")
}

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

func TestAggregator_Receive(t *testing.T) {

	t.Run("aggregate two messages", func(t *testing.T) {
		ti := time.Date(2001, 2, 3, 4, 5, 6, 0, time.UTC)
		create := time.Date(2023, 2, 3, 4, 5, 6, 0, time.UTC)
		ps := makeCreatePhotoString(ti, 60)

		msg1 := &message.PhotoMessage{ID: "id1", CreationTime: &create, Photo: ps}
		s := 60
		msg2 := &message.PhotoMessage{ID: "id1", Photo: ps, Speed: &s}
		expect := &message.PhotoMessage{ID: "id1", Photo: ps, CreationTime: msg1.CreationTime, Speed: msg2.Speed}

		system := actor.NewActorSystem()
		p := stream.NewTypedStream[*message.PhotoMessage](system)
		go func() {
			re := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return NewAggregator(1*time.Second, p.PID())
			}))
			system.Root.Send(re, msg1)
			system.Root.Send(re, msg2)
		}()
		res := <-p.C()
		if !reflect.DeepEqual(res, expect) {
			t.Errorf("expected %v, got %v", expect, res)
		}
	})

	t.Run("send message after timeout", func(t *testing.T) {
		ti := time.Date(2001, 2, 3, 4, 5, 6, 0, time.UTC)
		ps := makeCreatePhotoString(ti, 60)
		msg1 := &message.PhotoMessage{ID: "id1", CreationTime: &ti, Photo: ps}

		system := actor.NewActorSystem()
		p := stream.NewTypedStream[*message.PhotoMessage](system)
		go func() {
			re := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return NewAggregator(1*time.Second, p.PID())
			}))
			system.Root.Send(re, msg1)
		}()
		res := <-p.C()
		if !reflect.DeepEqual(res, msg1) {
			t.Errorf("expected %v, got %v", msg1, res)
		}
	})

	t.Run("aggregate two messages when restarting", func(t *testing.T) {
		ti := time.Date(2001, 2, 3, 4, 5, 6, 0, time.UTC)
		ps := makeCreatePhotoString(ti, 60)
		create := time.Date(2023, 2, 3, 4, 5, 6, 0, time.UTC)

		msg1 := &message.PhotoMessage{ID: "id1", CreationTime: &create, Photo: ps}
		s := 60
		msg2 := &message.PhotoMessage{ID: "id1", Photo: ps, Speed: &s}
		expect := &message.PhotoMessage{ID: "id1", Photo: ps, CreationTime: msg1.CreationTime, Speed: msg2.Speed}

		system := actor.NewActorSystem()
		p := stream.NewTypedStream[*message.PhotoMessage](system)
		go func() {
			re := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return NewAggregator(1*time.Second, p.PID())
			}))
			system.Root.Send(re, msg1)
			system.Root.Send(re, &message.IllegalStatePanicMessage{})
			system.Root.Send(re, msg2)
		}()
		res := <-p.C()
		if !reflect.DeepEqual(res, expect) {
			t.Errorf("expected %v, got %v", expect, res)
		}
	})
}

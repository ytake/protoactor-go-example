package scattergather

import (
	"reflect"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/stream"
	"github.com/ytake/protoactor-go-http/structure/message"
)

func TestScatterGather(t *testing.T) {
	timeOut := 2 * time.Second

	t.Run("scatter the message and gather them again", func(t *testing.T) {

		system := actor.NewActorSystem()
		p := stream.NewTypedStream[*message.PhotoMessage](system)

		d := time.Date(2023, 9, 10, 4, 5, 6, 0, time.UTC)
		speed := 60
		ps := makeCreatePhotoString(d, speed)

		msg := &message.PhotoMessage{
			ID:    "id1",
			Photo: ps,
		}
		go func() {
			ref := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return NewAggregator(timeOut, p.PID())
			}))
			speedRef := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return NewGetSpeedActor(ref)
			}))
			timeRef := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return NewGetTimeActor(ref)
			}))
			actorRef := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return &RecipientList{
					PIDs: []*actor.PID{speedRef, timeRef}}
			}))
			system.Root.Send(actorRef, msg)
		}()
		expect := &message.PhotoMessage{
			ID:           msg.ID,
			Photo:        msg.Photo,
			Speed:        &speed,
			CreationTime: &d,
		}
		res := <-p.C()
		if !reflect.DeepEqual(res, expect) {
			t.Errorf("expected %v, got %v", expect, res)
		}
	})
}

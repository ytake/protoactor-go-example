package filter

import (
	"errors"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/protoactor-go-http/structure/message"
)

func TestFilters(t *testing.T) {
	system := actor.NewActorSystem()
	license := system.Root.Spawn(actor.PropsFromProducer(NewLicense))
	speed := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return NewSpeed(50, license)
	}))
	msg := &message.Photo{License: "123xyz", Speed: 60}
	f := system.Root.RequestFuture(speed, msg, 2*time.Millisecond)
	res, err := f.Result()
	if err != nil {
		t.Error(err)
	}
	if res.(*message.Photo).License != msg.License {
		t.Errorf("expected %v, got %v", msg.License, res.(*message.Photo).License)
	}
	f = system.Root.RequestFuture(speed, &message.Photo{License: "", Speed: 60}, 10*time.Millisecond)
	_, err = f.Result()
	if !errors.Is(err, actor.ErrTimeout) {
		t.Errorf("expected %v, got %v", actor.ErrTimeout, err)
	}
	f = system.Root.RequestFuture(speed, &message.Photo{License: "123xyz", Speed: 49}, 10*time.Millisecond)
	_, err = f.Result()
	if !errors.Is(err, actor.ErrTimeout) {
		t.Errorf("expected %v, got %v", actor.ErrTimeout, err)
	}
}

package license

import (
	"errors"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/router"
	"github.com/ytake/protoactor-go-example/routing/prefrouting/message"
	"github.com/ytake/protoactor-go-example/routing/prefrouting/process"
)

func makeCreatePhotoString(ti time.Time, s int, license string) string {
	i := &process.ImageProcessing{}
	return i.CreatePhotoString(ti, s, license)
}

func TestGetLicense_Receive(t *testing.T) {
	t.Run("kill all routee", func(t *testing.T) {
		system := actor.NewActorSystem()
		ti := time.Date(2001, 2, 3, 4, 5, 6, 0, time.UTC)
		msg := &message.PerformanceRoutingMessage{
			Photo: makeCreatePhotoString(ti, 60, "123xyz")}
		var pids []*actor.PID
		for range [2]int{} {
			li := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return NewGetLicense(0)
			}, actor.WithSupervisor(
				actor.NewOneForOneStrategy(10, 1000, NewStopDirective().Decider()))))
			pids = append(pids, li)
		}
		ref := system.Root.Spawn(router.NewRoundRobinGroup(pids...))
		system.Root.Send(ref, "KillFirst")
		system.Root.Send(ref, "KillFirst")
		time.Sleep(1000 * time.Millisecond)
		f := system.Root.RequestFuture(ref, msg, 100*time.Millisecond)
		res, err := f.Result()
		if res != nil {
			t.Errorf("expected %v, got %v", nil, res)
		}
		if !errors.Is(err, actor.ErrDeadLetter) {
			t.Errorf("expected %v, got %v", actor.ErrDeadLetter, err)
		}
	})

	// 片方のアクターを落とした場合、もう片方のアクターが処理を行う
	t.Run("kill one routee", func(t *testing.T) {
		system := actor.NewActorSystem()
		ti := time.Date(2001, 2, 3, 4, 5, 6, 0, time.UTC)
		msg := &message.PerformanceRoutingMessage{
			Photo: makeCreatePhotoString(ti, 60, "123xyz")}
		var pids []*actor.PID
		for range [2]int{} {
			li := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return NewGetLicense(0)
			}, actor.WithSupervisor(
				actor.NewOneForOneStrategy(1, 1000, NewStopDirective().Decider()))))
			pids = append(pids, li)
		}
		ref := system.Root.Spawn(router.NewRoundRobinGroup(pids...))
		system.Root.Send(ref, "KillFirst")
		time.Sleep(1000 * time.Millisecond)
		f := system.Root.RequestFuture(ref, msg, 100*time.Millisecond)
		res, err := f.Result()
		if res == nil {
			t.Errorf("expected %v, got %v", nil, res)
		}
		if err != nil {
			t.Errorf("expected %v, got %v", nil, err)
		}
	})
}

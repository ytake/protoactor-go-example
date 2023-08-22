package ticket_actor

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/protoactor-go-http/ticket/order"
)

type Root struct {
	system *actor.ActorSystem
	pid    *actor.PID
}

// ActorSystem is a sample
func (ra *Root) ActorSystem() *actor.ActorSystem {
	return ra.system
}

// PID is a pid
func (ra *Root) PID() *actor.PID {
	return ra.pid
}

// NewBoxOfficeActorSystem is a sample
func NewBoxOfficeActorSystem() (*Root, error) {
	system := actor.NewActorSystem()
	supervisor := actor.NewOneForOneStrategy(10, 1000, NewStopDirective().Decider())
	props := actor.PropsFromProducer(
		order.NewBoxOfficeAPIActor,
		actor.WithSupervisor(supervisor),
	)
	pid, err := system.Root.SpawnNamed(props, order.BoxOfficeName)
	return &Root{
		system: system,
		pid:    pid,
	}, err
}

type StopDecider struct{}

func NewStopDirective() *StopDecider {
	return &StopDecider{}
}

func (d *StopDecider) Decider() actor.DeciderFunc {
	return func(reason interface{}) actor.Directive {
		fmt.Println("handling failure for child")
		return actor.StopDirective
	}
}

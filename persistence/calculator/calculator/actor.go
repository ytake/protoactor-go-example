package calculator

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"github.com/ytake/protoactor-go-example/persistence/calculator/command"
	"github.com/ytake/protoactor-go-example/persistence/calculator/event"
	"google.golang.org/protobuf/proto"
)

type PersistenceActor struct {
	persistence.Mixin
	state *CalculationResult
}

func NewPersistenceActor() actor.Actor {
	return &PersistenceActor{
		state: &CalculationResult{},
	}
}

// Receive for persistence actor
// Akka/Pekkoの receiveRecover, receiveCommandを一つにまとめた形で実装します
// persistenceのイベント名はcontext.Self().Idに紐づく形になっています
func (pa *PersistenceActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *persistence.RequestSnapshot:
		pa.PersistSnapshot(pa.state)
	case *command.Add:
		pa.persist(&event.Added{Result: msg.Value})
	case *command.Subtract:
		pa.persist(&event.Subtracted{Result: msg.Value})
	case *command.Divide:
		pa.persist(&event.Divided{Result: msg.Value})
	case *command.Multiply:
		pa.persist(&event.Multiplied{Result: msg.Value})
	case *command.Clear:
		pa.persist(&event.Reset{})
	case *command.GetResult:
		ctx.Respond(pa.state)
	case *command.PrintResult:
		fmt.Printf("Result is %v\n", pa.state.Result)
	}
}

func (pa *PersistenceActor) persist(msg proto.Message) {
	if !pa.Recovering() {
		pa.PersistReceive(msg)
	}
	switch ev := msg.(type) {
	case *event.Added:
		pa.state.Add(ev.Result)
	case *event.Subtracted:
		pa.state.Subtract(ev.Result)
	case *event.Divided:
		pa.state.Divide(ev.Result)
	case *event.Multiplied:
		pa.state.Multiply(ev.Result)
	case *event.Reset:
		pa.state.Reset()
	}
}

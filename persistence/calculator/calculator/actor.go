package calculator

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"github.com/ytake/protoactor-go-example/persistence/calculator/command"
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
		pa.persist(pa.state.Add(msg.Value))
	case *command.Subtract:
		pa.persist(pa.state.Subtract(msg.Value))
	case *command.Divide:
		pa.persist(pa.state.Divide(msg.Value))
	case *command.Multiply:
		pa.persist(pa.state.Multiply(msg.Value))
	case *command.Clear:
		pa.persist(pa.state.Reset())
	case *command.GetResult:
		ctx.Respond(pa.state)
	case *command.PrintResult:
		fmt.Printf("Result is %v\n", pa.state.Result)
	}
}

func (pa *PersistenceActor) persist(msg *CalculationResult) {
	if !pa.Recovering() {
		pa.PersistReceive(msg)
	}
	pa.state = msg
}

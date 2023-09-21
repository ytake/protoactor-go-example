package basket

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"github.com/ytake/protoactor-go-example/persistence/basket/command"
	"github.com/ytake/protoactor-go-example/persistence/basket/event"
	"github.com/ytake/protoactor-go-example/persistence/basket/protobuf"
	"github.com/ytake/protoactor-go-example/persistence/basket/value"
	"google.golang.org/protobuf/proto"
)

type Actor struct {
	ShopperID int
}

func (a Actor) Name() string {
	return fmt.Sprintf("basket_%d", a.ShopperID)
}

type PersistenceActor struct {
	persistence.Mixin
	nrEventsRecovered int
	state             *value.Items
}

func NewPersistenceActor() actor.Actor {
	return &PersistenceActor{
		nrEventsRecovered: 0,
		state:             value.NewItems(),
	}
}

func (pa *PersistenceActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case proto.Message:
		// 永続化されたイベントからのリカバリ時には、receiveRecoverが呼ばれます
		// この時に、proto.Messageを実装した構造体が送られてきます。
		//　proto actorではイベントやスナップショットを永続化するにはproto.Messageを実装した構造体である必要があります
		// この構造体を使って、指定のindexからアクターの状態を復元します
		// 内容を確認したい場合は下記をどうぞ
		// fmt.Println("Recovering basket from event", msg)
		pa.nrEventsRecovered++
		pa.updateState(msg)
	case *persistence.RequestSnapshot:
		pa.PersistSnapshot(pa.state)
	case *persistence.OfferSnapshot:
		fmt.Println("Recovering baskets from snapshot")
		pa.PersistSnapshot(pa.state)
	case *command.Add:
		pa.persist(&event.Added{Item: msg.Item.Item})
	case *command.GetItems:
		ctx.Respond(pa.state)
	case *command.RemoveItem:
		removed := &event.ItemRemoved{ProductID: msg.ProductID}
		ctx.Respond(removed)
		pa.persist(removed)
	case *command.Clear:
		pa.persist(&event.Cleared{Items: pa.state.Items})
	case *command.CountRecoveredEvents:
		ctx.Respond(command.RecoveredEventsCount{Count: pa.nrEventsRecovered})
	}
}

// receiveRecover is called when recovering from a snapshot or starting initially
// receiveRecoverは、スナップショットからのリカバリ時にproto.Messageを実装した構造体が送られてくるため、
// その構造体を使って再度アクターの状態を復元します
func (pa *PersistenceActor) receiveRecover(msg proto.Message) {
	pa.updateState(msg)
}

// persist persists the event and updates the state of the actor
func (pa *PersistenceActor) persist(msg proto.Message) {
	// Recovering()は、現在のアクターがリカバリ中かどうかを示します
	if !pa.Recovering() {
		pa.PersistReceive(msg)
	}
	pa.updateState(msg)
}

// updateState updates the state of the actor
func (pa *PersistenceActor) updateState(msg proto.Message) {
	switch ev := msg.(type) {
	case *event.Added:
		pa.state = pa.state.Add(&value.Item{Item: ev.Item})
	case *event.Cleared:
		pa.state = pa.state.Clear()
		//購入後にカートをクリアし、スナップショットを取得する
		pa.PersistSnapshot(pa.state)
	case *event.ItemRemoved:
		pa.state = pa.state.RemoveItem(ev.ProductID)
	case *protobuf.Items:
		pa.state = value.AggregateItems(ev.Items)
	}
}

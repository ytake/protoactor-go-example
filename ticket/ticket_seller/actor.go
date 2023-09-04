package ticket_seller

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/protoactor-go-http/ticket/message"
)

// TicketSellerActor is  Actor
type TicketSellerActor struct {
	Event
}

// NewTicketSellerActor is a filesystem based provider
func NewTicketSellerActor() actor.Actor {
	return &TicketSellerActor{}
}

// Receive is sent messages to be processed from the mailbox associated with the instance of the root
func (t *TicketSellerActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.Add:
		t.CreateEvent(msg.Name, msg.Tickets, context.Self().GetId())
	case *message.GetEvent:
		context.RequestWithCustomSender(
			context.Sender(),
			&message.Event{Name: t.Event.Name, Tickets: t.Event.Tickets},
			context.Parent())
	case *message.Cancel:
		// イベント自体のキャンセルとしてアクターを停止する
		// 再びアクターを利用する場合は、新たに作成する必要がある
		context.RequestWithCustomSender(
			context.Sender(),
			&message.Cancel{},
			context.Parent())
		context.Poison(context.Self())
	}
}

type Event struct {
	// tickets is the number of tickets
	Tickets int
	// name is the name of the box_office
	Name string
	id   string
}

// CreateEvent is a ticket_seller
func (t *TicketSellerActor) CreateEvent(name string, num int, id string) {
	t.Event = Event{Tickets: num, Name: name, id: id}
}

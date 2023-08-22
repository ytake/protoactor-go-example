package order

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/protoactor-go-http/ticket/ticket_seller"
)

// TicketSeller is a filesystem based provider
type TicketSeller struct{}

// NewTicketSellerActor is a filesystem based provider
func NewTicketSellerActor() actor.Actor {
	return &TicketSeller{}
}

// Receive is sent messages to be processed from the mailbox associated with the instance of the ticket_actor
func (r *TicketSeller) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *ticket_seller.Add:
		ticket_seller.CreateEvent(msg.Name, msg.Tickets, context.Self().GetId())
	case *ticket_seller.GetEvent:
		context.Respond(ticket_seller.GetEventByID(context.Self().GetId()))
	}
}

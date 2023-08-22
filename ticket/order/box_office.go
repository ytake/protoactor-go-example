package order

import (
	"errors"
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/protoactor-go-http/ticket/box_office"
	"github.com/ytake/protoactor-go-http/ticket/message"
	"github.com/ytake/protoactor-go-http/ticket/ticket_seller"
	"time"
)

const (
	BoxOfficeName = "BoxOffice"
)

type BoxOffice struct{}

func NewBoxOfficeAPIActor() actor.Actor {
	return &BoxOffice{}
}

// Receive is sent messages to be processed from the mailbox associated with the instance of the ticket_actor
func (r *BoxOffice) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.EventDescription:
		seller, err := context.SpawnNamed(actor.PropsFromProducer(NewTicketSellerActor), msg.Name)
		if errors.Is(err, actor.ErrNameExists) {
			// すでに存在する場合は、イベントが存在することを通知する
			context.Respond(&box_office.EventExists{})
			return
		}
		// イベントを作成する
		context.Send(seller, &ticket_seller.Add{Name: msg.Name, Tickets: msg.Tickets})
		// イベントが作成されたことを通知する
		context.Respond(&box_office.EventCreated{})
	case *box_office.GetEvents:
		var evs ticket_seller.Events
		for _, v := range context.Children() {
			f := context.RequestFuture(v, &ticket_seller.GetEvent{}, 2*time.Second)
			r, _ := f.Result()
			switch r.(type) {
			case *ticket_seller.Event:
				evs = append(evs, r.(*ticket_seller.Event))
			}
		}
		context.Send(context.Sender(), evs)
	case *box_office.GetEvent:
		// proto actorのforwardはメッセージをそのままに転送する
		// そのまま利用する場合、送信元がticket_seller.GetEventの扱いを知って送ることになるため、
		// ここではbox_officeで受け取ったメッセージをticket_sellerに転送し、
		// 送信先をエンドポイントに対応しているアクターに戻すよう指示している
		for _, v := range context.Children() {
			if v.GetId() == fmt.Sprintf("%s/%s", context.Self().GetId(), msg.Name) {
				// アクターはチケット名ごとに送信されるため、送信時にチケット名を指定する必要はない
				context.RequestWithCustomSender(v, &ticket_seller.GetEvent{}, context.Sender())
				return
			}
		}
		context.Respond(&box_office.EventNotFound{})
	}
}

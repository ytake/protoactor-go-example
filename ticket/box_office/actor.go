package box_office

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/protoactor-go-http/ticket/message"
	"github.com/ytake/protoactor-go-http/ticket/ticket_seller"
)

const (
	BoxOfficeName = "BoxOffice"
)

type BoxOfficeActor struct{}

func NewBoxOfficeAPIActor() actor.Actor {
	return &BoxOfficeActor{}
}

// Receive is sent messages to be processed from the mailbox associated with the instance of the root
func (r *BoxOfficeActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.EventDescription:
		seller, err := context.SpawnNamed(actor.PropsFromProducer(ticket_seller.NewTicketSellerActor), msg.Name)
		if errors.Is(err, actor.ErrNameExists) {
			// すでに存在する場合は、イベントが存在することを通知する
			context.Respond(&message.EventExists{})
			return
		}
		// イベントを作成する
		context.Send(seller, &message.Add{Name: msg.Name, Tickets: msg.Tickets})
		// イベントが作成されたことを通知する
		context.Respond(&message.EventCreated{})
	case *message.GetEvents:
		context.Send(context.Sender(), r.getEvents(context))
	case *message.GetEvent:
		// proto actorのforwardはメッセージをそのままに転送する
		// そのまま利用する場合、送信元がticket_seller.GetEventの扱いを知って送ることになるため、
		// ここではbox_officeで受け取ったメッセージをticket_sellerに転送し、
		// 送信先をエンドポイントに対応しているアクターに戻すよう指示している
		match := false
		for _, v := range context.Children() {
			if v.GetId() == fmt.Sprintf("%s/%s", context.Self().GetId(), msg.Name) {
				match = true
				// アクターはチケット名ごとに送信されるため、送信時にチケット名を指定する必要はない
				context.RequestWithCustomSender(v, &message.GetEvent{}, context.Sender())
			}
		}
		if !match {
			context.Respond(&message.EventNotFound{})
		}
	case *message.CancelEvent:
		for _, v := range context.Children() {
			if v.GetId() == fmt.Sprintf("%s/%s", context.Self().GetId(), msg.Name) {
				context.RequestWithCustomSender(v, &message.Cancel{}, context.Sender())
				return
			}
		}
		context.Respond(&message.EventNotFound{})
	}
}

// getEvents is a box_office
func (r *BoxOfficeActor) getEvents(context actor.Context) message.Events {
	var evs message.Events
	var wg sync.WaitGroup
	for _, v := range context.Children() {
		wg.Add(1)
		f := context.RequestFuture(v, &message.GetEvent{}, 2*time.Second)
		go func() {
			defer wg.Done()
			res, err := f.Result()
			if err != nil {
				return
			}
			switch msg := res.(type) {
			case *message.Event:
				evs = append(evs, msg)
			}
		}()
	}
	wg.Wait()
	return evs
}

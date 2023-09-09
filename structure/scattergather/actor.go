package scattergather

import (
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/protoactor-go-http/structure/message"
	"github.com/ytake/protoactor-go-http/structure/process"
)

type GetSpeed struct {
	pipe *actor.PID
}

func NewGetSpeedActor(pipe *actor.PID) actor.Actor {
	return &GetSpeed{
		pipe: pipe,
	}
}

func (g *GetSpeed) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.PhotoMessage:
		i := &process.ImageProcessing{}
		speed, _ := i.GetSpeed(msg.Photo)
		msg.Speed = &speed
		context.Send(g.pipe, msg)
	}
}

type GetTime struct {
	pipe *actor.PID
}

func NewGetTimeActor(pipe *actor.PID) actor.Actor {
	return &GetTime{
		pipe: pipe,
	}
}

func (g *GetTime) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.PhotoMessage:
		i := &process.ImageProcessing{}
		t, _ := i.GetTime(msg.Photo)
		msg.CreationTime = &t
		context.Send(g.pipe, msg)
	}
}

type RecipientList struct {
	PIDs []*actor.PID
}

func (state *RecipientList) Receive(context actor.Context) {
	var msgs []*message.PhotoMessage
	var wg sync.WaitGroup
	switch msg := context.Message().(type) {
	case *message.PhotoMessage:
		for _, recipient := range state.PIDs {
			wg.Add(1)
			f := context.RequestFuture(recipient, msg, 1*time.Second)
			go func() {
				defer wg.Done()
				res, err := f.Result()
				if err != nil {
					return
				}
				switch msg := res.(type) {
				case *message.PhotoMessage:
					msgs = append(msgs, msg)
				}
			}()
		}
		wg.Wait()
		context.Respond(msgs)
	}
}

// Aggregator is a aggregator
type Aggregator struct {
	timeout  time.Duration
	pipe     *actor.PID
	messages []*message.PhotoMessage
}

// NewAggregator is aggregator
func NewAggregator(timeout time.Duration, pipe *actor.PID) actor.Actor {
	return &Aggregator{
		timeout: timeout,
		pipe:    pipe,
	}
}

// Receive is called when a new message is sent to the actor
// 2つのメッセージを受け取る前提の例です
func (state *Aggregator) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Restarting:
		// panicを起こすとアクターの再起動を開始する
		// panicで停止した際に持っていたメッセージを再度自分自身に送り、
		// 現在メッセージを破棄する
		for _, v := range state.messages {
			ctx.Send(ctx.Self(), v)
		}
		state.messages = nil
	case *message.PhotoMessage:
		var found *message.PhotoMessage
		for _, m := range state.messages {
			if m.ID == msg.ID {
				found = m
				break
			}
		}
		if found != nil {
			newCombinedMsg := &message.PhotoMessage{
				ID:           msg.ID,
				Photo:        msg.Photo,
				CreationTime: msg.CreationTime,
				Speed:        msg.Speed,
			}
			newCombinedMsg.UpdatePhotoMessage(found)
			ctx.Send(state.pipe, newCombinedMsg)
			var newMessages []*message.PhotoMessage
			for _, m := range state.messages {
				if m.ID != found.ID {
					newMessages = append(newMessages, m)
				}
			}
			state.messages = newMessages
		} else {
			// メッセージが一件しかない場合はタイムアウトを設定した後に
			// アクターにメッセージを持たせている
			state.messages = append(state.messages, msg)
			ctx.SetReceiveTimeout(state.timeout)
		}
	case *actor.ReceiveTimeout:
		// タイムアウトが発生した場合（メッセージが一件しかない場合）
		// タイムアウトはproto actorではactor.ReceiveTimeoutというメッセージが送られる
		if len(state.messages) == 0 {
			return
		}
		ctx.Send(state.pipe, state.messages[0])
	case *message.IllegalStatePanicMessage:
		// 意図的にパニックを起こすメッセージを受信するとパニックを検知して、
		// Aggregatorのメッセージを持っている状態で停止する
		panic("this is a scheduled panic")
	}
}

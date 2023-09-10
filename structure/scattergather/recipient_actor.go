package scattergather

import (
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/protoactor-go-example/structure/message"
)

type RecipientList struct {
	PIDs []*actor.PID
}

func (state *RecipientList) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.PhotoMessage:
		var msgs []*message.PhotoMessage
		var wg sync.WaitGroup
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

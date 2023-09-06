package filter

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/protoactor-go-http/structure/message"
)

type Speed struct {
	minSpeed int
	license  *actor.PID
}

// NewSpeed is scattergather
func NewSpeed(minSpeed int, license *actor.PID) actor.Actor {
	return &Speed{minSpeed: minSpeed, license: license}
}

// Receive is sent messages to be processed from the mailbox associated with the instance of the Speed
func (r *Speed) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.Photo:
		if msg.Speed > r.minSpeed {
			f := context.RequestFuture(r.license, msg, 12*time.Millisecond)
			res, _ := f.Result()
			context.Respond(res)
		}
	}
}

type License struct {
	act *actor.PID
}

func NewLicense() actor.Actor {
	return &License{}
}

// Receive is sent messages to be processed from the mailbox associated with the instance of the License
func (l *License) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.Photo:
		if !msg.NoLicense() {
			fmt.Println("sending photo", msg.License)
			context.Respond(msg)
		}
	}
}

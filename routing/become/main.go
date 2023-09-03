package main

import (
	"log"

	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"
)

// SwitchRouter is actor
type SwitchRouter struct {
	behavior actor.Behavior
}

type (
	RouteStateOn  struct{}
	RouteStateOff struct{}
	Message       struct {
		Text string
	}
)

// Receive is sent messages to be processed from the mailbox associated with the instance of the SwitchRouter
func (s *SwitchRouter) Receive(context actor.Context) {
	s.behavior.Receive(context)
}

// On is state of SwitchRouter
// BecomeでOnになる
func (s *SwitchRouter) On(context actor.Context) {
	switch msg := context.Message().(type) {
	case *RouteStateOn:
		log.Println("Received on while already in on state")
	case *RouteStateOff:
		s.behavior.Become(s.Off)
	case *Message:
		log.Printf("Received message while in on state: %v", msg)
	}
}

// Off is state of SwitchRouter
// BecomeでOffになる
func (s *SwitchRouter) Off(context actor.Context) {
	switch msg := context.Message().(type) {
	case *RouteStateOn:
		s.behavior.Become(s.On)
	case *RouteStateOff:
		s.behavior.Become(s.Off)
		log.Println("Received off while already in off state")
	case *Message:
		log.Printf("Received message while in off state: %v", msg)
	}
}

func NewSetBehaviorActor() actor.Actor {
	act := &SwitchRouter{
		behavior: actor.NewBehavior(),
	}
	act.behavior.Become(act.Off)
	return act
}

func main() {
	system := actor.NewActorSystem()
	props := actor.PropsFromProducer(NewSetBehaviorActor)
	pid := system.Root.Spawn(props)
	// 初期はOff
	// RouteStateOffでメッセージを受けていることをログで確認できる
	// Received message while in off state: &{hello}
	system.Root.Send(pid, &Message{
		Text: "hello",
	})
	// Offの状態のままOffを指示すると下記が出力される
	// Received off while already in off state
	system.Root.Send(pid, &RouteStateOff{})
	// Onに変更指示
	system.Root.Send(pid, &RouteStateOn{})
	// Onの状態のままOnを指示すると下記が出力される
	// Received on while already in on state
	system.Root.Send(pid, &RouteStateOn{})
	// Onの状態でメッセージを送ると下記が出力される
	// Received message while in on state: &{hello}
	system.Root.Send(pid, &Message{
		Text: "hello",
	})
	// 再びOffに変更指示
	system.Root.Send(pid, &RouteStateOff{})
	// Offの状態でメッセージを送ると下記が出力される
	// Received message while in off state: &{done}
	system.Root.Send(pid, &Message{
		Text: "done",
	})
	_, _ = console.ReadLine()
}

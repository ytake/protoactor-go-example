package payload

import (
	"github.com/ytake/protoactor-go-example/ticket/message"
)

// Event is a payload
type Event struct {
	Tickets int    `json:"tickets,omitempty"`
	Name    string `json:"name,omitempty"`
}

type Events struct {
	E []Event `json:"events,omitempty"`
}

type CancelEvent struct{}

// NewEvents is a constructor for Event
func NewEvents(i interface{}) Events {
	var ev Events
	switch msg := i.(type) {
	case message.Events:
		for _, v := range msg {
			ev.E = append(ev.E, Event{
				Tickets: v.Tickets,
				Name:    v.Name,
			})
		}
		return ev
	}
	return ev
}

func NewEvent(i interface{}) Event {
	var ev Event
	switch msg := i.(type) {
	case *message.Event:
		ev = Event{
			Tickets: msg.Tickets,
			Name:    msg.Name,
		}
		return ev
	case *message.EventNotFound:
		return ev
	}
	return ev
}

func (e Event) IsValid() bool {
	return e.Tickets != 0 && e.Name != ""
}

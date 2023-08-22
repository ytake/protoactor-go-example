package message

// EventDescription is a message
type EventDescription struct {
	// Name is the name of the event
	Name string
	// Tickets is the number of tickets
	Tickets int
}

type GetEvents struct{}

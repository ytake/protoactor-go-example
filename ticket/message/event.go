package message

// EventDescription is a message
type EventDescription struct {
	// Name is the name of the event
	Name string
	// Tickets is the number of tickets
	Tickets int
}

// GetEvents is a box_office
// 全てのイベントを取得する
type GetEvents struct{}

type GetEvent struct {
	// Name is the name of the box_office
	Name string
}

type Events []*Event

// EventCreated is a box_office
type EventCreated struct{}
type EventExists struct{}
type EventNotFound struct{}

// CancelEvent is a box_office
type CancelEvent struct {
	// Name is the name of the box_office
	Name string
}

type Event struct {
	// tickets is the number of tickets
	Tickets int
	// name is the name of the box_office
	Name string
}

type Add struct {
	Tickets int
	Name    string
}

// Cancel is a ticket_seller
type Cancel struct{}

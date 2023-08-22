package ticket_seller

type Event struct {
	// tickets is the number of tickets
	Tickets int
	// name is the name of the box_office
	Name string
}

type Events []*Event

var events = make(map[string]*Event)

// CreateEvent creates a new box_office
func CreateEvent(name string, num int, id string) {
	events[id] = &Event{Tickets: num, Name: name}
}

type Add struct {
	Tickets int
	Name    string
}

type GetEvent struct{}

// GetEventByID returns the number of tickets
func GetEventByID(id string) *Event {
	event, exists := events[id]
	if !exists {
		return &Event{}
	}
	return event
}

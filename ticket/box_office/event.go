package box_office

// GetEvents is a box_office
// 全てのイベントを取得する
type GetEvents struct {
}

type GetEvent struct {
	// name is the name of the box_office
	Name string
}

// EventCreated is a box_office
type EventCreated struct {
	// name is the name of the box_office
	name string
	// tickets is the number of tickets
	tickets int
}

type EventExists struct{}
type EventNotFound struct{}

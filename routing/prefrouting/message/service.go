package message

import "time"

type SetService struct {
	ID          string
	ServiceTime time.Duration
}

type PerformanceRoutingMessage struct {
	Photo       string
	License     *string
	ProcessedBy *string
}

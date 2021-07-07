package event

import "github.com/AsynkronIT/protoactor-go/actor"

type (
	Credit struct {
		ReplyTo *actor.PID
		Amount  int
	}
	Debit struct {
		ReplyTo *actor.PID
		Amount  int
	}
	InternalServerError struct{}
	Refused             struct{}
	ServiceUnavailable  struct{}
)

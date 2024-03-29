package route

import (
	"net/http"
	"strconv"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/labstack/echo/v4"
	"github.com/ytake/protoactor-go-example/ticket/message"
	"github.com/ytake/protoactor-go-example/ticket/root"
)

type CreateEvent struct {
	actor *root.Root
}

func NewCreateEvent(actor *root.Root) *CreateEvent {
	return &CreateEvent{
		actor: actor,
	}
}

type ActorActionRouter interface {
	retrievePID() *actor.PID
	Handle(c echo.Context) error
}

func (ce *CreateEvent) retrievePID() *actor.PID {
	return ce.actor.PID()
}

func (ce *CreateEvent) Handle(c echo.Context) error {
	name := c.Param("name")
	tickets := c.FormValue("tickets")
	t, err := strconv.Atoi(tickets)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	ev := &message.EventDescription{Name: name, Tickets: t}
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	future := ce.actor.ActorSystem().Root.RequestFuture(ce.retrievePID(), ev, 2*time.Second)
	r, err := future.Result()
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	switch r.(type) {
	case *message.EventCreated:
		return c.String(http.StatusOK, "ok")
	case *message.EventExists:
		return c.String(http.StatusConflict, "already exists")
	}
	return c.NoContent(http.StatusNoContent)
}

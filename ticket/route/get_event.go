package route

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/labstack/echo/v4"
	"github.com/ytake/protoactor-go-http/ticket/box_office"
	"github.com/ytake/protoactor-go-http/ticket/payload"
	"github.com/ytake/protoactor-go-http/ticket/ticket_actor"
	"net/http"
	"time"
)

type GetEvent struct {
	actor *ticket_actor.Root
}

// NewGetEvent create new instance
func NewGetEvent(actor *ticket_actor.Root) *GetEvent {
	return &GetEvent{
		actor: actor,
	}
}

func (ce *GetEvent) retrievePID() *actor.PID {
	return ce.actor.PID()
}

func (ce *GetEvent) Handle(c echo.Context) error {
	ev := &box_office.GetEvent{Name: c.Param("name")}
	future := ce.actor.ActorSystem().Root.RequestFuture(ce.retrievePID(), ev, 2*time.Second)
	r, err := future.Result()
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	p := payload.NewEvent(r)
	if !p.IsValid() {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, p)
}

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

type GetEvents struct {
	actor *ticket_actor.Root
}

// NewGetEvents create new instance
func NewGetEvents(actor *ticket_actor.Root) *GetEvents {
	return &GetEvents{
		actor: actor,
	}
}

func (ce *GetEvents) retrievePID() *actor.PID {
	return ce.actor.PID()
}

func (ce *GetEvents) Handle(c echo.Context) error {
	ev := &box_office.GetEvents{}
	future := ce.actor.ActorSystem().Root.RequestFuture(ce.retrievePID(), ev, 2*time.Second)
	r, err := future.Result()
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, payload.NewEvents(r))
}

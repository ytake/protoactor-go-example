package route

import (
	"net/http"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/labstack/echo/v4"
	"github.com/ytake/protoactor-go-example/ticket/message"
	"github.com/ytake/protoactor-go-example/ticket/payload"
	"github.com/ytake/protoactor-go-example/ticket/root"
)

// GetEvents is a router
type GetEvents struct {
	actor *root.Root
}

// NewGetEvents create new instance
func NewGetEvents(actor *root.Root) *GetEvents {
	return &GetEvents{
		actor: actor,
	}
}

func (ce *GetEvents) retrievePID() *actor.PID {
	return ce.actor.PID()
}

func (ce *GetEvents) Handle(c echo.Context) error {
	ev := &message.GetEvents{}
	f := ce.actor.ActorSystem().Root.RequestFuture(ce.retrievePID(), ev, 2*time.Second)
	r, err := f.Result()
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, payload.NewEvents(r))
}

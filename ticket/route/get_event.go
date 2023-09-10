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

type GetEvent struct {
	actor *root.Root
}

// NewGetEvent create new instance
func NewGetEvent(actor *root.Root) *GetEvent {
	return &GetEvent{
		actor: actor,
	}
}

func (ce *GetEvent) retrievePID() *actor.PID {
	return ce.actor.PID()
}

func (ce *GetEvent) Handle(c echo.Context) error {
	ev := &message.GetEvent{Name: c.Param("name")}
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

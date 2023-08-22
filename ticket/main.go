package main

import (
	"github.com/labstack/echo/v4"
	"github.com/ytake/protoactor-go-http/ticket/route"
	"github.com/ytake/protoactor-go-http/ticket/ticket_actor"
	"log"
)

func main() {
	as, err := ticket_actor.NewBoxOfficeActorSystem()
	if err != nil {
		log.Fatal(err)
	}
	e := echo.New()
	e.POST("/events/:name", route.NewCreateEvent(as).Handle)
	e.GET("/events/:name", route.NewGetEvent(as).Handle)
	e.GET("/events", route.NewGetEvents(as).Handle)
	e.Logger.Fatal(e.Start(":8080"))
}

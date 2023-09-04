package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/ytake/protoactor-go-http/ticket/root"
	"github.com/ytake/protoactor-go-http/ticket/route"
)

func main() {
	as, err := root.NewBoxOfficeActorSystem()
	if err != nil {
		log.Fatal(err)
	}
	e := echo.New()
	e.POST("/events/:name", route.NewCreateEvent(as).Handle)
	e.GET("/events/:name", route.NewGetEvent(as).Handle)
	e.GET("/events", route.NewGetEvents(as).Handle)
	e.DELETE("/events/:name", route.NewCancelEvent(as).Handle)
	e.Logger.Fatal(e.Start(":8080"))
}

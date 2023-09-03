package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type MyActor struct{}

func (state *MyActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case string:
		ctx.Respond(fmt.Sprintf("Hello, %s", msg))
	}
}

func NewMyActor() actor.Actor {
	return &MyActor{}
}

func main() {
	system := actor.NewActorSystem()
	props := actor.PropsFromProducer(NewMyActor)
	myActor := system.Root.Spawn(props)

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "World"
		}
		future := system.Root.RequestFuture(myActor, name, 2*time.Second)
		result, err := future.Result()
		if err != nil {
			http.Error(w, "Error occurred", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, result.(string))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

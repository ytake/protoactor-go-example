package main

import (
	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ytake/protoactor-go-example/task"
	"log"
)

var (
	numberOfTransfers  = 1000
	uptime             = 99.99
	refusalProbability = 0.01
	busyProbability    = 0.01
	retryAttempts      = 0
	verbose            = false
	system = actor.NewActorSystem()
)

func main() {
	log.Println("booted")
	props := actor.PropsFromProducer(func() actor.Actor {return &task.Account{}})
	system.Root.Spawn(props)
	_, _ = console.ReadLine()
}

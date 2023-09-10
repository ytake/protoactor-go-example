package main

import (
	"fmt"
	"time"

	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/stream"
	"github.com/ytake/protoactor-go-http/structure/message"
	"github.com/ytake/protoactor-go-http/structure/process"
	"github.com/ytake/protoactor-go-http/structure/scattergather"
)

var timeOut = 2 * time.Second

// makePhoto is mane photo string
func makePhoto(t time.Time, speed int, license string) string {
	image := &process.ImageProcessing{}
	return image.CreatePhotoString(t, speed, license)
}

func main() {
	system := actor.NewActorSystem()
	p := stream.NewTypedStream[*message.PhotoMessage](system)

	msg := &message.PhotoMessage{
		ID:    "id1",
		Photo: makePhoto(time.Now(), 60, "scattering"),
	}
	go func() {
		ref := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return scattergather.NewAggregator(timeOut, p.PID())
		}))
		speedRef := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return scattergather.NewGetSpeedActor(ref)
		}))
		timeRef := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return scattergather.NewGetTimeActor(ref)
		}))
		actorRef := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &scattergather.RecipientList{
				PIDs: []*actor.PID{speedRef, timeRef}}
		}))
		system.Root.Send(actorRef, msg)
	}()
	res := <-p.C()
	// 写真名から作成日時とスピードを取得する個別のアクターが処理を行います
	// それらのアクターは、それぞれの処理が完了すると、
	// aggregatorアクターにメッセージを送信し、それぞれの結果を集約します
	// このように、複数のアクターにメッセージを送信し、
	// それらの処理結果を集約することをscatter-gatherパターンと呼びます
	// 結果は、streamを通じて受け取ります
	fmt.Println(res)
	_, _ = console.ReadLine()
}

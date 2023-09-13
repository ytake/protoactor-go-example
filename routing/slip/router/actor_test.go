package router

import (
	"reflect"
	"testing"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/stream"
	"github.com/ytake/protoactor-go-example/routing/slip/message"
	"github.com/ytake/protoactor-go-example/routing/slip/option"
	"github.com/ytake/protoactor-go-example/routing/slip/task"
	"github.com/ytake/protoactor-go-example/routing/slip/value"
)

func TestSlipRouter_Receive(t *testing.T) {
	t.Run("route messages correctly", func(t *testing.T) {
		system := actor.NewActorSystem()
		p := stream.NewTypedStream[*message.Car](system)
		go func() {
			pb, _ := system.Root.SpawnNamed(actor.PropsFromProducer(func() actor.Actor {
				return &option.PaintCar{Color: "black", RouteSlip: &task.RouteSlip{}}
			}), "paintBlack")
			sr, _ := system.Root.SpawnNamed(actor.PropsFromProducer(func() actor.Actor {
				return &SlipRouter{
					PaintBlack: pb,
					RouteSlip:  &task.RouteSlip{},
					EndStep:    p.PID()}
			}), "slipRouter")
			system.Root.Send(sr, &message.Order{})
		}()
		r := <-p.C()
		expect := &message.Car{Color: "black", HasNavigation: false, HasParkingSensors: false}
		if !reflect.DeepEqual(r, expect) {
			t.Errorf("expected %v, got %v", expect, r)
		}
	})

	t.Run("car full option", func(t *testing.T) {
		system := actor.NewActorSystem()
		order := message.Order{
			Options: []value.CarOptions{value.CarColorGray, value.Navigation, value.ParkingSensor}}
		p := stream.NewTypedStream[*message.Car](system)
		go func() {
			pb, _ := system.Root.SpawnNamed(actor.PropsFromProducer(func() actor.Actor {
				return &option.PaintCar{Color: "gray", RouteSlip: &task.RouteSlip{}}
			}), "paintGray")
			an, _ := system.Root.SpawnNamed(actor.PropsFromProducer(func() actor.Actor {
				return &option.AddNavigation{RouteSlip: &task.RouteSlip{}}
			}), "navigation")
			ap, _ := system.Root.SpawnNamed(actor.PropsFromProducer(func() actor.Actor {
				return &option.AddParkingSensor{RouteSlip: &task.RouteSlip{}}
			}), "parkingSensor")
			sr, _ := system.Root.SpawnNamed(actor.PropsFromProducer(func() actor.Actor {
				return &SlipRouter{
					PaintGray:        pb,
					AddNavigation:    an,
					AddParkingSensor: ap,
					RouteSlip:        &task.RouteSlip{},
					EndStep:          p.PID()}
			}), "slipRouter")
			system.Root.Send(sr, &order)
		}()
		r := <-p.C()
		expect := &message.Car{Color: "gray", HasNavigation: true, HasParkingSensors: true}
		if !reflect.DeepEqual(r, expect) {
			t.Errorf("expected %v, got %v", expect, r)
		}
	})

	t.Run("parking sensor option", func(t *testing.T) {
		system := actor.NewActorSystem()
		order := message.Order{
			Options: []value.CarOptions{value.ParkingSensor}}
		p := stream.NewTypedStream[*message.Car](system)
		go func() {
			pb, _ := system.Root.SpawnNamed(actor.PropsFromProducer(func() actor.Actor {
				return &option.PaintCar{Color: "black", RouteSlip: &task.RouteSlip{}}
			}), "paintBlack")
			pg, _ := system.Root.SpawnNamed(actor.PropsFromProducer(func() actor.Actor {
				return &option.PaintCar{Color: "gray", RouteSlip: &task.RouteSlip{}}
			}), "paintGray")
			ps, _ := system.Root.SpawnNamed(actor.PropsFromProducer(func() actor.Actor {
				return &option.AddParkingSensor{RouteSlip: &task.RouteSlip{}}
			}), "parkingSensor")
			an, _ := system.Root.SpawnNamed(actor.PropsFromProducer(func() actor.Actor {
				return &option.AddNavigation{RouteSlip: &task.RouteSlip{}}
			}), "navigation")
			sr, _ := system.Root.SpawnNamed(actor.PropsFromProducer(func() actor.Actor {
				return &SlipRouter{
					PaintBlack:       pb,
					PaintGray:        pg,
					AddNavigation:    an,
					AddParkingSensor: ps,
					RouteSlip:        &task.RouteSlip{},
					EndStep:          p.PID()}
			}), "slipRouter")
			system.Root.Send(sr, &order)
		}()
		r := <-p.C()
		expect := &message.Car{Color: "black", HasNavigation: false, HasParkingSensors: true}
		if !reflect.DeepEqual(r, expect) {
			t.Errorf("expected %v, got %v", expect, r)
		}
	})
}

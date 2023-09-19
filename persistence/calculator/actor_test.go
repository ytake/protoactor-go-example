package calculator

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"github.com/ytake/protoactor-go-example/persistence/command"
	"github.com/ytake/protoactor-go-example/persistence/protobuf"
	"github.com/ytake/protoactor-go-example/persistence/provider"
)

func TestPersistenceActor_Receive(t *testing.T) {
	t.Run("recover last known result after crash", func(t *testing.T) {
		system := actor.NewActorSystem()
		p := provider.NewInMemory(1)
		props := actor.PropsFromProducer(NewPersistenceActor,
			actor.WithReceiverMiddleware(persistence.Using(p)))
		pid, _ := system.Root.SpawnNamed(props, "calculator")
		system.Root.Send(pid, &command.Add{Value: 1})
		f := system.Root.RequestFuture(pid, &command.GetResult{}, 2*time.Second)
		res, _ := f.Result()
		calc := &protobuf.CalculationResult{Result: 1}
		if reflect.DeepEqual(res, calc) {
			t.Errorf("expected %v, got %v", calc, res)
		}
		system.Root.Send(pid, &command.Subtract{Value: 0.5})
		f = system.Root.RequestFuture(pid, &command.GetResult{}, 2*time.Second)
		res, _ = f.Result()
		calc = &protobuf.CalculationResult{Result: 0.5}
		if reflect.DeepEqual(res, calc) {
			t.Errorf("expected %v, got %v", calc, res)
		}

		// 生成したactorを意図的に落とします
		_ = system.Root.PoisonFuture(pid).Wait()
		fmt.Printf("*** actor restart ***\n")
		// 再度actorを生成します 落とす前と同じIDをアクターに割り当てます
		pid2, _ := system.Root.SpawnNamed(props, "calculator")
		// ここで前回の結果を復元します
		f = system.Root.RequestFuture(pid2, &command.GetResult{}, 2*time.Second)
		res, _ = f.Result()
		// 前回の結果が復元されていることを確認します
		calc = &protobuf.CalculationResult{Result: 0.5}
		if reflect.DeepEqual(res, calc) {
			t.Errorf("expected %v, got %v", calc, res)
		}
		// 新しい計算を行います（加算 1）
		system.Root.Send(pid, &command.Add{Value: 1})
		f = system.Root.RequestFuture(pid2, &command.GetResult{}, 2*time.Second)
		res, _ = f.Result()
		// 新しい計算結果が反映されていることを確認します
		calc = &protobuf.CalculationResult{Result: 1.5}
		if reflect.DeepEqual(res, calc) {
			t.Errorf("expected %v, got %v", calc, res)
		}
	})
}

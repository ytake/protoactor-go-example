package basket

import (
	"os"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"github.com/ytake/protoactor-go-example/persistence/basket/command"
	"github.com/ytake/protoactor-go-example/persistence/basket/event"
	"github.com/ytake/protoactor-go-example/persistence/basket/protobuf"
	"github.com/ytake/protoactor-go-example/persistence/basket/source"
	"github.com/ytake/protoactor-go-example/persistence/basket/value"
)

const dbDir = "../db/actor_test"

func setup() {}

func teardown() {
	_ = os.RemoveAll(dbDir)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

var (
	shopperID     = 2
	macbookPro    = &protobuf.Item{ProductID: "Apple Macbook Pro", Number: 1, UnitPrice: 2499.99}
	macPro        = &protobuf.Item{ProductID: "Apple Mac Pro", Number: 1, UnitPrice: 10499.9}
	displays      = &protobuf.Item{ProductID: "4K Display", Number: 3, UnitPrice: 2499.99}
	appleMouse    = &protobuf.Item{ProductID: "Apple Mouse", Number: 1, UnitPrice: 99.99}
	appleKeyboard = &protobuf.Item{ProductID: "Apple Keyboard", Number: 1, UnitPrice: 79.99}
	dWave         = &protobuf.Item{ProductID: "D-Wave One", Number: 1, UnitPrice: 14999999.99}
)

func TestPersistenceActor_Receive(t *testing.T) {
	// インメモリ永続化を使った復元テスト
	t.Run("in memory persistence", func(t *testing.T) {
		system := actor.NewActorSystem()
		p := source.NewInMemory(1)
		props := actor.PropsFromProducer(NewPersistenceActor,
			actor.WithReceiverMiddleware(persistence.Using(p)))
		pid, _ := system.Root.SpawnNamed(props, Actor{ShopperID: shopperID}.Name())
		system.Root.Send(pid, &command.Add{ShopperID: shopperID, Item: &value.Item{Item: macbookPro}})
		// 生成したactorを意図的に落とします
		_ = system.Root.PoisonFuture(pid).Wait()
		basketResurrected, _ := system.Root.SpawnNamed(props, Actor{ShopperID: shopperID}.Name())
		brf := system.Root.RequestFuture(basketResurrected, &command.GetItems{ShopperID: shopperID}, 2*time.Second)
		r, _ := brf.Result()
		items := r.(*value.Items)
		if !slices.Contains(items.Items.Items, macbookPro) {
			t.Errorf("invalid items")
		}
	})

	// GoLevelDB永続化を使った復元テスト
	// Akka in Actionの例と同じテストを実行します
	t.Run("skip basket events that occurred before Cleared during recovery", func(t *testing.T) {
		system := actor.NewActorSystem()
		p, err := source.NewGoLevelDBProvider(1, dbDir)
		if err != nil {
			t.Errorf("failed to create provider: %v", err)
		}
		props := actor.PropsFromProducer(NewPersistenceActor,
			actor.WithReceiverMiddleware(persistence.Using(p)))
		pid, _ := system.Root.SpawnNamed(props, Actor{ShopperID: shopperID}.Name())
		system.Root.Send(pid, &command.Add{ShopperID: shopperID, Item: &value.Item{Item: macbookPro}})
		system.Root.Send(pid, &command.Add{ShopperID: shopperID, Item: &value.Item{Item: displays}})
		f := system.Root.RequestFuture(pid, &command.GetItems{ShopperID: shopperID}, 2*time.Second)
		r, err := f.Result()
		if err != nil {
			t.Errorf("failed to get items: %v", err)
		}
		items := r.(*value.Items)
		if len(items.Items.Items) != 2 {
			t.Errorf("invalid items")
		}
		if !slices.Contains(items.Items.Items, macbookPro) {
			t.Errorf("invalid items")
		}
		if !slices.Contains(items.Items.Items, displays) {
			t.Errorf("invalid items")
		}
		// 一度カートを削除
		system.Root.Send(pid, &command.Clear{ShopperID: shopperID})
		system.Root.Send(pid, &command.Add{ShopperID: shopperID, Item: &value.Item{Item: macPro}})
		f2 := system.Root.RequestFuture(pid, &command.RemoveItem{ProductID: macPro.ProductID, ShopperID: shopperID}, 2*time.Second)
		r, err = f2.Result()
		if err != nil {
			t.Errorf("failed to get removed items: %v", err)
		}
		removed := r.(*event.ItemRemoved)
		if removed.ProductID != macPro.ProductID {
			t.Errorf("invalid removed item")
		}
		// 再度カートを削除
		system.Root.Send(pid, &command.Clear{ShopperID: shopperID})
		system.Root.Send(pid, &command.Add{ShopperID: shopperID, Item: &value.Item{Item: dWave}})
		system.Root.Send(pid, &command.Add{ShopperID: shopperID, Item: &value.Item{Item: displays}})
		f = system.Root.RequestFuture(pid, &command.GetItems{ShopperID: shopperID}, 2*time.Second)
		r, err = f.Result()
		if err != nil {
			t.Errorf("failed to get items: %v", err)
		}
		items = r.(*value.Items)
		if len(items.Items.Items) != 2 {
			t.Errorf("invalid items")
		}
		// 生成したactorを意図的に落とします
		_ = system.Root.PoisonFuture(pid).Wait()
		basketResurrected, _ := system.Root.SpawnNamed(props, Actor{ShopperID: shopperID}.Name())
		brf := system.Root.RequestFuture(basketResurrected, &command.GetItems{ShopperID: shopperID}, 2*time.Second)
		r, err = brf.Result()
		if err != nil {
			t.Errorf("failed to get items: %v", err)
		}
		// 再度カートが復元されていることを確認します
		items = r.(*value.Items)
		if len(items.Items.Items) != 2 {
			t.Errorf("invalid items")
		}
		// アクターが復元された時点のカートの状態になっていることを確認します
		if reflect.DeepEqual(items.Items.Items[0], dWave) {
			t.Errorf("invalid dWave want: %v, got: %v", dWave, items.Items.Items[0])
		}
		if reflect.DeepEqual(items.Items.Items[1], displays) {
			t.Errorf("invalid displays want: %v, got: %v", displays, items.Items.Items[0])
		}
		// 復元された回数を取得します
		// snapshot作成指示からの回数であることを確認します
		brf = system.Root.RequestFuture(basketResurrected, &command.CountRecoveredEvents{ShopperID: shopperID}, 2*time.Second)
		r, err = brf.Result()
		if err != nil {
			t.Errorf("failed to get items: %v", err)
		}
		count := r.(command.RecoveredEventsCount)
		if !reflect.DeepEqual(count, command.RecoveredEventsCount{Count: 2}) {
			t.Errorf("invalid RecoveredEventsCount want: %v, got: %v", 2, count)
		}
	})
}

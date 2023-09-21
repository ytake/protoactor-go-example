package source

import (
	"os"
	"reflect"
	"testing"

	"github.com/ytake/protoactor-go-example/persistence/basket/protobuf"
	"github.com/ytake/protoactor-go-example/persistence/basket/value"
)

const dbDir = "../db/testing"

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

func TestGoLevelDBProvider_PersistEvent(t *testing.T) {
	p, err := NewGoLevelDBProvider(3, dbDir)
	defer p.db.Close()
	if err != nil {
		t.Errorf("failed to create provider: %v", err)
	}
	item := &value.Item{Item: &protobuf.Item{ProductID: "1234", Number: 1, UnitPrice: 1000}}
	p.PersistEvent("test", 1, item)
	p.GetEvents("test", 1, 1, func(e interface{}) {
		if !reflect.DeepEqual(e, item.Item) {
			t.Errorf("invalid item")
		}
	})
}

func TestGoLevelDBProvider_PersistSnapshot(t *testing.T) {
	p, err := NewGoLevelDBProvider(3, dbDir)
	defer p.db.Close()
	if err != nil {
		t.Errorf("failed to create provider: %v", err)
	}
	item := &value.Item{Item: &protobuf.Item{ProductID: "1234", Number: 1, UnitPrice: 1000}}
	p.PersistSnapshot("test", 1, item)
	snapshot, _, ok := p.GetSnapshot("test")
	if !ok {
		t.Errorf("snapshot not found")
	}
	if !reflect.DeepEqual(snapshot, item.Item) {
		t.Errorf("invalid snapshot")
	}
}

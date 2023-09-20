package source

import (
	"reflect"
	"testing"

	"github.com/ytake/protoactor-go-example/persistence/basket/protobuf"
)

func TestEnvelope_Message(t *testing.T) {
	item := &protobuf.Item{ProductID: "1234", Number: 1, UnitPrice: 1000}
	env := newEnvelope(item, "event", 1)
	if !reflect.DeepEqual(env.message(), item) {
		t.Errorf("expected %v, got %v", item, env.message())
	}
}

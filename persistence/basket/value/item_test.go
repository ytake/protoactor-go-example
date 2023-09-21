package value

import (
	"fmt"
	"testing"

	"github.com/ytake/protoactor-go-example/persistence/basket/protobuf"
)

func TestAggregateItems(t *testing.T) {
	var items []*protobuf.Item
	a := &protobuf.Item{ProductID: "1234", Number: 1, UnitPrice: 1000}
	items = append(items, a)
	agg := AggregateItems(items)
	if len(agg.Items.Items) != 1 {
		t.Errorf("invalid aggregate items")
	}
	b := &protobuf.Item{ProductID: "1234", Number: 1, UnitPrice: 1000}
	items = append(items, b)
	agg = AggregateItems(items)
	if len(agg.Items.Items) != 2 {
		t.Errorf("invalid aggregate items")
	}
}

func TestItem_Aggregate(t *testing.T) {
	item := &protobuf.Item{ProductID: "1234", Number: 1, UnitPrice: 1000}
	a := &Item{Item: item}
	b := &Item{Item: &protobuf.Item{ProductID: "1234", Number: 1, UnitPrice: 1000}}
	agg := a.Aggregate(b)
	fmt.Println(agg)
}

func TestItem_Update(t *testing.T) {
	item := &protobuf.Item{ProductID: "1234", Number: 1, UnitPrice: 1000}
	a := &Item{Item: item}
	if a.Update(2).Number != 2 {
		t.Errorf("invalid update")
	}
}

func TestNewItems(t *testing.T) {
	items := NewItems(
		&Item{Item: &protobuf.Item{ProductID: "1234", Number: 1, UnitPrice: 1000}},
		&Item{Item: &protobuf.Item{ProductID: "1", Number: 1, UnitPrice: 1000}},
	)
	if len(items.Items.Items) != 2 {
		t.Errorf("invalid items")
	}
}

func TestItems_Add(t *testing.T) {
	items := NewItems(
		&Item{Item: &protobuf.Item{ProductID: "1234", Number: 1, UnitPrice: 1000}},
		&Item{Item: &protobuf.Item{ProductID: "1", Number: 1, UnitPrice: 1000}},
	)
	items = items.Add(&Item{Item: &protobuf.Item{ProductID: "1234", Number: 1, UnitPrice: 1000}})
	if len(items.Items.Items) != 3 {
		t.Errorf("invalid items")
	}
}

func TestItems_AddItems(t *testing.T) {
	items := NewItems(
		&Item{Item: &protobuf.Item{ProductID: "1234", Number: 1, UnitPrice: 1000}},
		&Item{Item: &protobuf.Item{ProductID: "1", Number: 1, UnitPrice: 1000}},
	)
	items2 := NewItems(
		&Item{Item: &protobuf.Item{ProductID: "1234", Number: 1, UnitPrice: 1000}},
		&Item{Item: &protobuf.Item{ProductID: "1", Number: 1, UnitPrice: 1000}},
	)
	items = items.AddItems(items2)
	if len(items.Items.Items) != 4 {
		t.Errorf("invalid items")
	}
}

func TestItems_ContainsProduct(t *testing.T) {
	items := NewItems(
		&Item{Item: &protobuf.Item{ProductID: "1234", Number: 1, UnitPrice: 1000}},
		&Item{Item: &protobuf.Item{ProductID: "1", Number: 1, UnitPrice: 1000}},
	)
	if !items.ContainsProduct("1234") {
		t.Errorf("invalid items")
	}
	if items.ContainsProduct("2") {
		t.Errorf("invalid items")
	}
}

func TestItems_RemoveItem(t *testing.T) {
	items := NewItems(
		&Item{Item: &protobuf.Item{ProductID: "1234", Number: 1, UnitPrice: 1000}},
		&Item{Item: &protobuf.Item{ProductID: "1", Number: 1, UnitPrice: 1000}},
	)
	items = items.RemoveItem("1234")
	if len(items.Items.Items) != 1 {
		t.Errorf("invalid items")
	}
}

func TestItems_UpdateItem(t *testing.T) {
	items := NewItems(
		&Item{Item: &protobuf.Item{ProductID: "1234", Number: 1, UnitPrice: 1000}},
		&Item{Item: &protobuf.Item{ProductID: "1", Number: 1, UnitPrice: 1000}},
	)
	items = items.UpdateItem("1234", 2)
	if items.Items.Items[0].Number != 2 {
		t.Errorf("invalid items")
	}
}

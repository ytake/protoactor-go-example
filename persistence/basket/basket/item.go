package basket

import "github.com/ytake/protoactor-go-example/persistence/basket/protobuf"

// Item is a struct
type Item struct {
	*protobuf.Item
}

type Items struct {
	Items *protobuf.Items
}

// Aggregate aggregates two items.
// ProductIDが同じなら、Numberを加算する
func (i *Item) Aggregate(item *Item) *Item {
	if i.Item.ProductID == item.Item.ProductID {
		return &Item{
			&protobuf.Item{
				ProductID: i.ProductID,
				Number:    i.Number + item.Number,
				UnitPrice: i.UnitPrice,
			},
		}
	}
	return nil
}

func (i *Item) Update(number int32) *Item {
	return &Item{
		&protobuf.Item{
			ProductID: i.ProductID,
			Number:    number,
			UnitPrice: i.UnitPrice,
		},
	}
}

func AggregateItems(list []*protobuf.Item) *Items {
	return &Items{
		&protobuf.Items{Items: list},
	}
}

func NewItems(args ...*Item) *Items {
	var items []*protobuf.Item
	for _, item := range args {
		items = append(items, item.Item)
	}
	return AggregateItems(items)
}

func (items *Items) Add(newItem *Item) *Items {
	return AggregateItems(append(items.Items.Items, newItem.Item))
}

func (items *Items) AddItems(newItems *Items) *Items {
	combinedList := append(items.Items.Items, newItems.Items.Items...)
	return AggregateItems(combinedList)
}

func (items *Items) ContainsProduct(productID string) bool {
	for _, item := range items.Items.Items {
		if item.ProductID == productID {
			return true
		}
	}
	return false
}

func (items *Items) RemoveItem(productID string) *Items {
	var newList []*protobuf.Item
	for _, item := range items.Items.Items {
		if item.ProductID != productID {
			newList = append(newList, item)
		}
	}
	return AggregateItems(newList)
}

func (items *Items) UpdateItem(productID string, number int) *Items {
	var newList []*protobuf.Item
	for _, item := range items.Items.Items {
		if item.ProductID == productID {
			item.Number = int32(number)
			newList = append(newList, item)
		} else {
			newList = append(newList, item)
		}
	}
	return AggregateItems(newList)
}

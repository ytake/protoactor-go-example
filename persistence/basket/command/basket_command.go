package command

import (
	"github.com/ytake/protoactor-go-example/persistence/basket/value"
)

type Add struct {
	Item      *value.Item
	ShopperID int
}

type Replace struct {
	Items     *value.Items
	ShopperID int
}

type RemoveItem struct {
	ProductID string
	ShopperID int
}

type UpdateItem struct {
	ProductID string
	Number    int
	ShopperID int
}

type Clear struct {
	ShopperID int
}

type GetItems struct {
	ShopperID int
}

type CountRecoveredEvents struct {
	ShopperID int
}

type RecoveredEventsCount struct {
	Count int
}

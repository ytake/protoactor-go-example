package command

type Add struct {
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

package models

const (
	OrderStatusPending  = "pending"
	OrderStatusfilled   = "filled"
	OrderStatusCanceled = "canceled"
	BuySide             = "buy"
	SellSide            = "sell"
)

var PermittedPairs = map[string]bool{
	"EUR-USD": true,
}

type Order struct {
	ID        string
	OwnerID   string
	Side      string
	Price     float64
	Quantity  float64
	AssetPair string
	Status    string
	FilledBy  string
}

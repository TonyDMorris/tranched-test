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
	ID        string  `json:"id"`
	OwnerID   string  `json:"owner_id"`
	Side      string  `json:"side"`
	Price     float64 `json:"price"`
	Quantity  float64 `json:"quantity"`
	AssetPair string  `json:"asset_pair"`
	Status    string  `json:"status"`
	FilledBy  string  `json:"filled_by"`
}

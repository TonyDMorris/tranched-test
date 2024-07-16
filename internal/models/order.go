package models

type Order struct {
	ID        string
	OwnerID   string
	Side      string
	Price     float64
	Quantity  float64
	AssetPair string
	Status    string
}

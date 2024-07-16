package order

import (
	"errors"
	"sync"

	"github.com/tonydmorris/tranched/internal/models"
	"github.com/tonydmorris/tranched/pkg/id"
)

const (
	OrderStatusPending  = "pending"
	OrderStatusfilled   = "filled"
	BuySide             = "buy"
	SellSide            = "sell"
	ErrInvalidAssetPair = "invalid asset pair"
)

var permittedPairs = map[string]bool{
	"EUR-USD": true,
}

// ContrivedOrderStorage is a simple in-memory storage for orders
// storage is keyed by Sybol -> Price -> Quantity -> []Order
type contrivedOrderStorageBySymbolPriceQuantity map[string]map[float64]map[float64][]*models.Order
type contrivedOrderStorageByUserID map[string][]*models.Order

type Repository struct {
	ordersBySymbol contrivedOrderStorageBySymbolPriceQuantity
	ordersByUserID contrivedOrderStorageByUserID
	idGen          idGenerator
	mu             *sync.RWMutex
}

func NewRepository() *Repository {
	return &Repository{
		ordersBySymbol: make(contrivedOrderStorageBySymbolPriceQuantity),
		ordersByUserID: make(contrivedOrderStorageByUserID),
		idGen:          id.New(),
		mu:             &sync.RWMutex{},
	}
}

func (r *Repository) CreateOrder(order models.Order) (models.Order, error) {
	_, ok := permittedPairs[order.AssetPair]
	if !ok {
		return models.Order{}, errors.New(ErrInvalidAssetPair)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	order.ID = r.idGen.New()
	order.Status = OrderStatusPending
	if _, ok := r.ordersBySymbol[order.AssetPair]; !ok {
		r.ordersBySymbol[order.AssetPair] = make(map[float64]map[float64][]*models.Order)
	}

	if _, ok := r.ordersBySymbol[order.AssetPair][order.Price]; !ok {
		r.ordersBySymbol[order.AssetPair][order.Price] = make(map[float64][]*models.Order)
	}

	validOrders := r.ordersBySymbol[order.AssetPair][order.Price][order.Quantity]

	r.ordersBySymbol[order.AssetPair][order.Price][order.Quantity] = append(validOrders, &order)

	for _, matchingOrder := range validOrders {
		if matchingOrder.Side == opositeSide(order.Side) {
			// fill the order
			matchingOrder.Status = OrderStatusfilled
			order.Status = OrderStatusfilled
			matchingOrder.FilledBy = order.OwnerID
			order.FilledBy = matchingOrder.OwnerID
		}
	}

	r.ordersByUserID[order.OwnerID] = append(r.ordersByUserID[order.OwnerID], &order)

	return order, nil

}

func (r *Repository) FindOrderByOwnerID(ownerID string) ([]models.Order, error) {

	var orders []models.Order

	for _, order := range r.ordersByUserID[ownerID] {
		orders = append(orders, *order)
	}

	return orders, nil

}

func opositeSide(side string) string {
	if side == BuySide {
		return SellSide
	}
	return BuySide
}

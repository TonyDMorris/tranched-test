package order

import (
	"errors"
	"sync"

	"github.com/tonydmorris/tranched/internal/models"
	"github.com/tonydmorris/tranched/pkg/id"
)

const (
	ErrInvalidAssetPair = "invalid asset pair"
)

// contrivedOrderStorage is a simple in-memory storage for orders
// storage is keyed by Sybol -> Price -> Quantity -> []Order
type contrivedOrderStorageBySymbolPriceQuantity map[string]map[float64]map[float64][]*models.Order

// contrivedOrderStorageByUserID is a simple in-memory storage for orders
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

// CreateOrder uses contrived in-memory storage to create an order
func (r *Repository) CreateOrder(order models.Order) (models.Order, error) {
	_, ok := models.PermittedPairs[order.AssetPair]
	if !ok {
		return models.Order{}, errors.New(ErrInvalidAssetPair)
	}
	// lock the "database"
	r.mu.Lock()
	defer r.mu.Unlock()

	// gen new ID for order
	order.ID = r.idGen.New()
	// set order status to pending
	order.Status = models.OrderStatusPending

	// instantiate the map if it doesn't exist
	if _, ok := r.ordersBySymbol[order.AssetPair]; !ok {
		r.ordersBySymbol[order.AssetPair] = make(map[float64]map[float64][]*models.Order)
	}
	// instantiate the map if it doesn't exist
	if _, ok := r.ordersBySymbol[order.AssetPair][order.Price]; !ok {
		r.ordersBySymbol[order.AssetPair][order.Price] = make(map[float64][]*models.Order)
	}
	// get the valid orders to match against
	validOrders := r.ordersBySymbol[order.AssetPair][order.Price][order.Quantity]
	for _, matchingOrder := range validOrders {
		if matchingOrder.Side == opositeSide(order.Side) {
			// fill the order
			matchingOrder.Status = models.OrderStatusfilled
			order.Status = models.OrderStatusfilled
			matchingOrder.FilledBy = order.OwnerID
			order.FilledBy = matchingOrder.OwnerID
		}
	}
	// append the new order to the list of valid orders
	r.ordersBySymbol[order.AssetPair][order.Price][order.Quantity] = append(validOrders, &order)
	// append the new order to the list of orders by user
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
	if side == models.BuySide {
		return models.SellSide
	}
	return models.BuySide
}

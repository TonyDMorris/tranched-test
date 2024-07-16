package order

import (
	"sync"

	"github.com/tonydmorris/tranched/internal/models"
)

const OrderStatusPending = "pending"
const OrderStatusfilled = "filled"
const BuySide = "buy"
const SellSide = "sell"

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
	}
}

func (r *Repository) CreateOrder(order models.Order) (models.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	order.ID = r.idGen.New()
	order.Status = "pending"
	if _, ok := r.ordersBySymbol[order.AssetPair]; !ok {
		r.ordersBySymbol[order.AssetPair] = make(map[float64]map[float64][]*models.Order)
	}

	if _, ok := r.ordersBySymbol[order.AssetPair][order.Price]; !ok {
		r.ordersBySymbol[order.AssetPair][order.Price] = make(map[float64][]*models.Order)
	}

	validOrders := r.ordersBySymbol[order.AssetPair][order.Price][order.Quantity]
	if len(validOrders) == 0 {
		r.ordersBySymbol[order.AssetPair][order.Price][order.Quantity] = append(validOrders, &order)
		return order, nil
	}

	for _, matchingOrder := range validOrders {
		if matchingOrder.Side == opositeSide(order.Side) {
			// fill the order
			matchingOrder.Status = OrderStatusfilled
			order.Status = OrderStatusfilled
			return order, nil
		}
	}

	r.ordersByUserID[order.OwnerID] = append(r.ordersByUserID[order.OwnerID], &order)

	return order, nil

}

func (r *Repository) FindByOwnerID(ownerID string) ([]models.Order, error) {

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

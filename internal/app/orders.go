package app

import (
	"fmt"
	"strings"

	"github.com/tonydmorris/tranched/internal/models"
)

// createOrder creates an order
// It commits the assets for the order
// It creates the order
// It settles the order if the order is filled
// Returns the created order
// This shoud be a transaction but for the sake of simplicity we are not using transactions due to in-memory order database and sql user database
func (a *App) CreateOrder(order models.Order) (*models.Order, error) {
	// guard clauses
	if order.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than 0")
	}
	if order.Price <= 0 {
		return nil, fmt.Errorf("price must be greater than 0")
	}
	// check if the user has enough assets to reserve
	err := a.commitAssets(order)
	if err != nil {
		return nil, fmt.Errorf("error reserving assets: %w", err)
	}
	// create the order
	order, err = a.orderRepository.CreateOrder(order)
	if err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}
	// settle the order if it is filled
	if order.Status == models.OrderStatusfilled {
		err = a.settle(order)
		if err != nil {
			return nil, fmt.Errorf("error settling order: %w", err)
		}
	}
	return &order, nil
}

func (a *App) GetOrders(username string) ([]models.Order, error) {
	user, err := a.userRepository.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	return a.orderRepository.FindOrderByOwnerID(user.ID)
}

func (a *App) commitAssets(order models.Order) error {
	assets := strings.Split(order.AssetPair, "-")
	var sybol string
	var amount float64
	// - A `SELL` side for `EUR-USD` pair, for an amount of `1300.0` and price `1.3` would mean for a user to "sell 1000 EUR to receive 1300.0 USD"
	// - A `BUY` side for `EUR-USD` pair, for an amount of `1200.0` and price `1.2` would mean for a user to "buy 1000.0 EUR spending 1200.0 USD"
	switch order.Side {
	case models.SellSide:
		sybol = assets[0]
		amount = order.Quantity / order.Price
	case models.BuySide:
		sybol = assets[1]
		amount = order.Quantity
	}
	// check if the user has enough assets to reserve
	err := a.userRepository.UpdateAssetByUserID(order.OwnerID, sybol, -amount)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) settle(order models.Order) error {
	symbols := strings.Split(order.AssetPair, "-")

	var ordererPayout float64
	var ordererSymbol string

	var matchedOrderPayout float64
	var matchedOrderSymbol string

	switch order.Side {
	case models.SellSide:
		ordererPayout = order.Quantity
		ordererSymbol = symbols[1]
		matchedOrderPayout = order.Quantity / order.Price
		matchedOrderSymbol = symbols[0]
	case models.BuySide:
		ordererPayout = order.Quantity / order.Price
		ordererSymbol = symbols[0]
		matchedOrderPayout = order.Quantity
		matchedOrderSymbol = symbols[1]
	}

	err := a.userRepository.UpdateAssetByUserID(order.OwnerID, ordererSymbol, ordererPayout)
	if err != nil {
		return err
	}

	err = a.userRepository.UpdateAssetByUserID(order.FilledBy, matchedOrderSymbol, matchedOrderPayout)
	if err != nil {
		// rollback the orderer's payout
		err = a.userRepository.UpdateAssetByUserID(order.OwnerID, ordererSymbol, -ordererPayout)
		if err != nil {
			a.logger.Errorf("error rolling back orderer's payout %q, with order ID %s ", err, order.ID)
		}
		return err
	}

	return nil
}

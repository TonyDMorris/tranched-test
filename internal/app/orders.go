package app

import "github.com/tonydmorris/tranched/internal/models"

func (a *App) createOrder(req CreateOrderRequest, username string) (*models.Order, error) {
	user, err := a.userRepository.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	order := models.Order{
		OwnerID:   user.ID,
		Side:      req.Side,
		Price:     req.Price,
		Quantity:  req.Quantity,
		AssetPair: req.AssetPair,
	}
	order, err = a.orderRepository.CreateOrder(order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (a *App) getOrders(username string) ([]models.Order, error) {
	user, err := a.userRepository.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	return a.orderRepository.FindByOwnerID(user.ID)
}

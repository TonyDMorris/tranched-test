package app

import (
	"github.com/tonydmorris/tranched/internal/models"
	"github.com/tonydmorris/tranched/internal/repository/order"
	"github.com/tonydmorris/tranched/internal/repository/user"
)

var _ UserRepository = &user.Repository{}

var _ OrderRepository = &order.Repository{}

type UserRepository interface {
	FindByUsername(username string) (models.User, error)
	CreateUser(username, passwordHash string) (models.User, error)
	FindAssetsByUsername(username string) ([]models.Asset, error)
}
type OrderRepository interface {
	CreateOrder(order models.Order) (models.Order, error)
	FindByOwnerID(ownerID string) ([]models.Order, error)
}

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

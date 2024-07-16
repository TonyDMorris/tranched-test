package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	logger          Logger
	userRepository  UserRepository
	orderRepository OrderRepository
	router          *gin.Engine
}

func New(ur UserRepository, or OrderRepository, opts ...func(a *App)) *App {
	a := &App{
		userRepository:  ur,
		orderRepository: or,
	}

	for _, opt := range opts {
		opt(a)
	}

	if a.logger == nil {
		log, _ := zap.NewProduction()
		a.logger = log.Sugar()
	}

	if a.router == nil {
		a.router = gin.Default()
	}

	return a
}

func (a *App) Route() {

	a.router.GET("/orders", a.AuthenticateRequest, a.GetOrders)
	a.router.POST("/orders", a.AuthenticateRequest, a.CreateOrder)

	a.router.POST("/users", a.CreateUser)
	a.router.GET("/assets", a.AuthenticateRequest, a.GetUserAssets)

}

func (a *App) Run(port int) error {
	a.Route()
	return a.router.Run(fmt.Sprintf(":%d", port))
}

func WithLogger(logger Logger) func(a *App) {
	return func(a *App) {
		a.logger = logger
	}
}

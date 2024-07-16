package app

import (
	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tonydmorris/tranched/internal/models"
)

// AuthenticateRequest is a middleware that authenticates a request
func (a *App) AuthenticateRequest(c *gin.Context) {

	// get the auth header
	creds := c.Request.Header.Get("Authorization")

	if creds == "" {
		c.JSON(401, gin.H{"error": "no credentials provided"})
		c.Abort()
		return
	}
	// extract the base64 encoded credentials
	encodedCredentials := creds[6:]

	decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		c.Abort()
		return
	}
	// split the credentials into username and password
	split := strings.Split(string(decodedCredentials), ":")
	if len(split) != 2 {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		c.Abort()
		return
	}

	username := split[0]
	password := split[1]
	// authenticate the user
	ok, err := a.authenticate(username, password)
	if err != nil {
		a.logger.Errorf("error authenticating user %q, with username %s", err, username)
		c.JSON(500, gin.H{"error": "internal server error"})
		c.Abort()
		return
	}

	if !ok {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		c.Abort()
		return
	}

	// set the user in the context
	c.Set("username", username)

	c.Next()
}

// CreateUser creates a new user
func (a *App) CreateUser(c *gin.Context) {
	var req CreateUserRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		c.Abort()
		return
	}

	user, err := a.createUser(req.Username, req.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		c.Abort()
		return
	}

	c.JSON(201, CreateUserResponse{ID: user.ID})
}

// CreateOrder creates a new order
func (a *App) CreateOrder(c *gin.Context) {
	// get the user from the context
	username, ok := c.Get("username")
	if !ok {
		c.JSON(401, gin.H{"error": "no credentials provided"})
		c.Abort()
		return
	}

	var req CreateOrderRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		c.Abort()
		return
	}

	user, err := a.getUser(username.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		c.Abort()
		return
	}

	order := models.Order{
		OwnerID:   user.ID,
		Side:      req.Side,
		Price:     req.Price,
		Quantity:  req.Quantity,
		AssetPair: req.AssetPair,
	}

	createdOrder, err := a.createOrder(order)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		c.Abort()
		return
	}

	c.JSON(201, CreateOrderResponse{
		ID:     createdOrder.ID,
		Status: createdOrder.Status,
	})

}

func (a *App) GetOrders(c *gin.Context) {
	// get the user from the context
	username, ok := c.Get("username")
	if !ok {
		c.JSON(401, gin.H{"error": "no credentials provided"})
		c.Abort()
		return
	}

	orders, err := a.getOrders(username.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		c.Abort()
		return
	}

	c.JSON(200, orders)
}

func (a *App) GetUserAssets(c *gin.Context) {
	// get the user from the context
	username, ok := c.Get("username")
	if !ok {
		c.JSON(401, gin.H{"error": "no credentials provided"})
		c.Abort()
		return
	}

	assets, err := a.getAssetsByUsername(username.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		c.Abort()
		return
	}

	c.JSON(200, assets)
}

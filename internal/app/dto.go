package app

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID string `json:"id"`
}

type CreateOrderRequest struct {
	Side      string  `json:"side"`
	Price     float64 `json:"price"`
	Quantity  float64 `json:"quantity"`
	AssetPair string  `json:"asset_pair"`
}

type CreateOrderResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

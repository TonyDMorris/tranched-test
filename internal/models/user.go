package models

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Assets       []Asset
}

type Asset struct {
	ID     string  `json:"id"`
	UserID string  `json:"user_id"`
	Symbol string  `json:"symbol"`
	Amount float64 `json:"amount"`
}

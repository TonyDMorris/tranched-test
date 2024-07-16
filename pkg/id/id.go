package id

import "github.com/google/uuid"

type gen struct{}

func New() *gen {
	return &gen{}
}

// NewID generates a new ID.
func (g gen) New() string {
	return uuid.New().String()
}

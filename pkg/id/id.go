package id

import "github.com/google/uuid"

type gen struct{}

// NewID generates a new ID.
func (g gen) New() string {
	return uuid.New().String()
}

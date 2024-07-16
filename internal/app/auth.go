package app

import (
	"github.com/tonydmorris/tranched/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (a *App) Authenticate(username, password string) (bool, error) {
	DBUser, err := a.userRepository.FindByUsername(username)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(DBUser.PasswordHash), []byte(password))
	return err == nil, err

}

func (a *App) CreateUser(username, password string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user, err := a.userRepository.CreateUser(username, string(hashedPassword))
	return &user, err
}

func (a *App) GetUser(username string) (*models.User, error) {
	user, err := a.userRepository.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

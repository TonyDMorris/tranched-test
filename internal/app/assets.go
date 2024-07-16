package app

import "github.com/tonydmorris/tranched/internal/models"

func (a *App) getAssetsByUsername(username string) ([]models.Asset, error) {
	return a.userRepository.FindAssetsByUsername(username)
}

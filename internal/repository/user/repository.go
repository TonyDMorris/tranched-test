package user

import (
	"database/sql"
	"fmt"

	"github.com/tonydmorris/tranched/internal/models"
)

// Repository is a user repository
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new user repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// FindByUsername retrieves a user by their username
func (r *Repository) FindByUsername(username string) (models.User, error) {
	row := r.db.QueryRow("SELECT (id, username, password_hash) FROM users WHERE username = ?", username)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		return models.User{}, fmt.Errorf("error scanning user: %s with error: %w", username, err)
	}
	return user, nil
}

// CreateUser creates a new user
func (r *Repository) CreateUser(username, passwordHash string) (models.User, error) {
	var user models.User
	user.Username = username
	user.PasswordHash = passwordHash

	tx, err := r.db.Begin()
	if err != nil {
		return models.User{}, fmt.Errorf("error beginning transaction: %w", err)
	}
	stmnt, err := tx.Prepare("INSERT INTO users (username, password_hash) VALUES (?, ?, ?) RETURNING id")
	if err != nil {
		return models.User{}, fmt.Errorf("error preparing statement: %w", err)
	}

	defer stmnt.Close()

	err = stmnt.QueryRow(user.Username, user.PasswordHash).Scan(&user.ID)
	if err != nil {
		tx.Rollback()
		return user, fmt.Errorf("error inserting user: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return user, fmt.Errorf("error committing transaction: %w", err)
	}

	return user, nil

}

func (r *Repository) FindAssetsByUsername(username string) ([]models.Asset, error) {
	rows, err := r.db.Query("SELECT (asset.id, asset.symbol, asset.amount) FROM users as user LEFT JOIN assets as asset on user.id = asset.user_id  WHERE user.username = ?", username)
	if err != nil {
		return nil, fmt.Errorf("error querying assets: %w", err)
	}
	defer rows.Close()

	var assets []models.Asset
	for rows.Next() {
		var asset models.Asset
		err = rows.Scan(&asset.ID, &asset.Symbol, &asset.Amount)
		if err != nil {
			return nil, fmt.Errorf("error scanning asset: %w", err)
		}
		asset.UserID = username
		assets = append(assets, asset)
	}

	return assets, nil
}

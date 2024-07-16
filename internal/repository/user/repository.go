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
	row := r.db.QueryRow("SELECT id, username, password_hash FROM public.users WHERE username = $1", username)

	var idVal sql.NullString
	var passwordHashVal sql.NullString
	var usernameVal sql.NullString
	err := row.Scan(&idVal, &usernameVal, &passwordHashVal)
	if err != nil {
		return models.User{}, fmt.Errorf("error scanning user: %s with error: %w", username, err)
	}
	var user models.User
	user.ID = idVal.String
	user.Username = usernameVal.String
	user.PasswordHash = passwordHashVal.String
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
	stmnt, err := tx.Prepare("INSERT INTO public.users (username, password_hash) VALUES (?, ?) RETURNING id")
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
	rows, err := r.db.Query("SELECT user_id ,assets.id, assets.symbol, assets.amount FROM public.users LEFT JOIN public.assets on users.id = assets.user_id  WHERE users.username = $1", username)
	if err != nil {
		return nil, fmt.Errorf("error querying assets: %w", err)
	}
	defer rows.Close()

	var assets []models.Asset
	for rows.Next() {
		var id sql.NullString
		var userID sql.NullString
		var symbol sql.NullString
		var amount sql.NullFloat64
		err = rows.Scan(&userID, &id, &symbol, &amount)
		if err != nil {
			return nil, fmt.Errorf("error scanning asset: %w", err)
		}
		var asset models.Asset
		asset.ID = id.String
		asset.UserID = userID.String
		asset.Symbol = symbol.String
		asset.Amount = amount.Float64
		assets = append(assets, asset)
	}

	return assets, nil
}

// UpdateAssetByuserID updates an asset by user id
func (r *Repository) UpdateAssetByUserID(userID, symbol string, amount float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	stmnt, err := tx.Prepare("UPDATE public.assets SET amount = amount + $1 WHERE user_id = $2 AND symbol = $3")
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}

	defer stmnt.Close()

	_, err = stmnt.Exec(amount, userID, symbol)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error updating asset: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

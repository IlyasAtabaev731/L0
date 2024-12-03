package storage

import (
	"database/sql"
	"errors"
	"github.com/IlyasAtabaev731/L0/internal/models"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

// GetOrder fetches an order and its related data from the database
func GetOrder(db *sql.DB, orderUid string) (*models.Order, error) {
	var order models.Order
	// Fetch order data
	// Fetch delivery, payment, and items
	// Populate the order struct
	return &order, nil
}

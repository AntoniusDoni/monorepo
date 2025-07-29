package model

import (
	"time"

	"github.com/google/uuid"
)

type StockEntry struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	WarehouseID   uuid.UUID `json:"warehouse_id"`
	ProductID     uuid.UUID `json:"product_id"`
	BatchNumber   string    `json:"batch_number"`
	ExpiredAt     time.Time `json:"expired_at"`
	Date          time.Time `json:"date"`
	Margin        float64   `json:"margin"`
	Tax           float64   `json:"tax"`
	Price         float64   `json:"price"`
	Stock         int       `json:"stock"`
	PreviousStock int       `json:"previous_stock"`
	Status        string    `json:"status"`
	OrderID       uuid.UUID `json:"order_id"`
	Notes         string    `json:"notes"`
	ReferenceID   uuid.UUID `json:"reference_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

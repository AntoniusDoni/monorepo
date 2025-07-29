package model

import (
	"time"

	"github.com/google/uuid"
)

type UnitProduct struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"` // Unique identifier
	Code      string    `json:"code"`                                                      // Product code or SKU
	Name      string    `json:"name"`                                                      // Product name
	CreatedAt time.Time `json:"created_at"`                                                // Timestamp when created
	UpdatedAt time.Time `json:"updated_at"`                                                // Timestamp when updated
}

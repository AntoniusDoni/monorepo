package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID                  uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"` // Unique identifier
	Code                string          `json:"code"`                                                      // Product code or SKU
	Name                string          `json:"name"`                                                      // Product name
	LargeUnit           string          `json:"large_unit"`                                                // e.g., box, pack
	ContentPerLargeUnit int             `json:"content_per_large_unit"`                                    // e.g., 12 pieces per box
	SmallUnit           string          `json:"small_unit"`                                                // e.g., piece, tablet
	PurchasePrice       float64         `json:"purchase_price"`                                            // Cost price
	SellingPrice        float64         `json:"selling_price"`                                             // Sale price
	CategoryID          uuid.UUID       `gorm:"type:uuid" json:"category_id"`                              // Foreign key
	Category            CategoryProduct `gorm:"foreignKey:CategoryID" json:"category"`
	Indication          string          `json:"indication"` // Description or usage
	CreatedAt           time.Time       `json:"created_at"` // Timestamp when created
	UpdatedAt           time.Time       `json:"updated_at"` // Timestamp when updated
}

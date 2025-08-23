package model

import (
	"time"

	"github.com/google/uuid"
)

type CategoryProduct struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"` // Unique identifier
	Name      string     `json:"name"`                                                      // Category name
	ParentID  *uuid.UUID `json:"parent_id,omitempty"`                                       // Parent category ID
	CreatedAt *time.Time `json:"created_at"`                                                // Timestamp when created
	UpdatedAt *time.Time `json:"updated_at,omitempty"`                                      // Timestamp when updated
}

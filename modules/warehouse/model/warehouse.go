package model

import (
	"time"

	"github.com/google/uuid"
)

type Warehouse struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Code      string     `gorm:"unique;not null" json:"code"`
	Name      string     `gorm:"not null" json:"name"`
	Address   string     `json:"address"`
	Phone     string     `json:"phone"`
	Status    string     `json:"status"`
	BranchID  *uuid.UUID `gorm:"type:uuid" json:"branch_id"` // nullable for now
	OfficeID  *uuid.UUID `gorm:"type:uuid" json:"office_id"` // nullable for now
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

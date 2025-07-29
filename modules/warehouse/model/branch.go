package model

import (
	"time"

	"github.com/google/uuid"
)

type Branch struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Code    string    `gorm:"unique;not null" json:"code"`
	Name    string    `gorm:"not null" json:"name"`
	Address string    `json:"address"`
	City    string    `json:"city"`
	Phone   string    `json:"phone"`
	Status  string    `json:"status"`

	OfficeID uuid.UUID `gorm:"type:uuid;not null" json:"office_id"`
	Office   Office    `gorm:"foreignKey:OfficeID" json:"office"`

	Warehouses []Warehouse `gorm:"foreignKey:BranchID" json:"warehouses,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

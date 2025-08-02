package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"unique;not null" json:"username"`
	PasswordHash string    `gorm:"not null" json:"-"`
	OfficeID     uuid.UUID `gorm:"type:uuid;index" json:"officeId"`
	Email        string    `gorm:"unique;not null" json:"email"`
	ApiToken     string    `gorm:"uniqueIndex" json:"-"`
	CreatedAt    time.Time `json:"created_at"`

	Roles []Role `gorm:"many2many:user_roles;" json:"roles"`
}

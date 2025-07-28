package model

type Permission struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique;not null" json:"name"`

	// Many-to-many relation with Role through role_permissions join table
	Roles []Role `gorm:"many2many:role_permissions;" json:"-"`
}

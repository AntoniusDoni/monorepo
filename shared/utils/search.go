package utils

import (
	"gorm.io/gorm"
)

// BuildSearchQuery applies a global search filter with ILIKE for multiple fields.
func BuildSearchQuery(db *gorm.DB, searchTerm string, fields []string) *gorm.DB {
	if searchTerm == "" || len(fields) == 0 {
		return db
	}

	like := "%" + searchTerm + "%"
	query := db
	for i, field := range fields {
		if i == 0 {
			query = query.Where(field+" ILIKE ?", like)
		} else {
			query = query.Or(field+" ILIKE ?", like)
		}
	}
	return query
}

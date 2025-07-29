package repository

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewRepository(ctx context.Context, db *gorm.DB) *Repository {
	return &Repository{ctx: ctx, db: db}
}

func (r *Repository) DB() *gorm.DB {
	return r.db.WithContext(r.ctx)
}

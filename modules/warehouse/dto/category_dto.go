package dto

import "github.com/google/uuid"

// CategoryProductCreateRequest represents the request body for creating a category product

type CategoryProductCreateRequest struct {
	Name     string    `json:"name" validate:"required" example:"Precursor"`                                 // Category name
	ParentID uuid.UUID `json:"parent_id" validate:"required" example:"123e4567-e89b-12d3-a456-426614174000"` //example:"123e4567-e89b-12d3-a456-426614174000"
}

package dto

import (
	"github.com/antoniusDoni/monorepo/modules/warehouse/model"
	"github.com/google/uuid"
)

// ProductCreateRequest represents the request body for creating a product
type ProductCreateRequest struct {
	Code                string    `json:"code" validate:"required" example:"PROD001"`                                     // Product code or SKU
	Name                string    `json:"name" validate:"required" example:"Laptop Dell XPS 13"`                          // Product name
	LargeUnit           string    `json:"large_unit" validate:"required" example:"box"`                                   // e.g., box, pack
	ContentPerLargeUnit int       `json:"content_per_large_unit" validate:"required,min=1" example:"12"`                  // e.g., 12 pieces per box
	SmallUnit           string    `json:"small_unit" validate:"required" example:"piece"`                                 // e.g., piece, tablet
	PurchasePrice       float64   `json:"purchase_price" validate:"required,min=0" example:"500.00"`                      // Cost price
	SellingPrice        float64   `json:"selling_price" validate:"required,min=0" example:"750.00"`                       // Sale price
	CategoryID          uuid.UUID `json:"category_id" validate:"required" example:"123e4567-e89b-12d3-a456-426614174000"` // Foreign key to category
	Indication          string    `json:"indication" example:"High-performance laptop for professionals"`                 // Description or usage
}

// ProductUpdateRequest represents the request body for updating a product
type ProductUpdateRequest struct {
	Code                string    `json:"code" validate:"required" example:"PROD001"`                                     // Product code or SKU
	Name                string    `json:"name" validate:"required" example:"Laptop Dell XPS 13"`                          // Product name
	LargeUnit           string    `json:"large_unit" validate:"required" example:"box"`                                   // e.g., box, pack
	ContentPerLargeUnit int       `json:"content_per_large_unit" validate:"required,min=1" example:"12"`                  // e.g., 12 pieces per box
	SmallUnit           string    `json:"small_unit" validate:"required" example:"piece"`                                 // e.g., piece, tablet
	PurchasePrice       float64   `json:"purchase_price" validate:"required,min=0" example:"500.00"`                      // Cost price
	SellingPrice        float64   `json:"selling_price" validate:"required,min=0" example:"750.00"`                       // Sale price
	CategoryID          uuid.UUID `json:"category_id" validate:"required" example:"123e4567-e89b-12d3-a456-426614174000"` // Foreign key to category
	Indication          string    `json:"indication" example:"High-performance laptop for professionals"`                 // Description or usage
}

// ToProduct converts ProductCreateRequest to Product model
func (req *ProductCreateRequest) ToProduct() *model.Product {
	return &model.Product{
		Code:                req.Code,
		Name:                req.Name,
		LargeUnit:           req.LargeUnit,
		ContentPerLargeUnit: req.ContentPerLargeUnit,
		SmallUnit:           req.SmallUnit,
		PurchasePrice:       req.PurchasePrice,
		SellingPrice:        req.SellingPrice,
		CategoryID:          req.CategoryID,
		Indication:          req.Indication,
	}
}

// ToProduct converts ProductUpdateRequest to Product model
func (req *ProductUpdateRequest) ToProduct() *model.Product {
	return &model.Product{
		Code:                req.Code,
		Name:                req.Name,
		LargeUnit:           req.LargeUnit,
		ContentPerLargeUnit: req.ContentPerLargeUnit,
		SmallUnit:           req.SmallUnit,
		PurchasePrice:       req.PurchasePrice,
		SellingPrice:        req.SellingPrice,
		CategoryID:          req.CategoryID,
		Indication:          req.Indication,
	}
}

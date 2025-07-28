package contract

// Generic PaginationRequest used across modules for paged endpoints
type PaginationRequest struct {
	Page     int `json:"page" form:"page" validate:"gte=1"`
	PageSize int `json:"page_size" form:"page_size" validate:"gte=1,lte=100"`
}

type ListRequest struct {
	PaginationRequest
	SearchTerm string `json:"search_term,omitempty"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

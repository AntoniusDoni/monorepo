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
	Email    string `json:"email" validate:"required,email"`
	OfficeID string `json:"office_id" validate:"required,uuid"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required" example:"superadmin"`
	Password string `json:"password" validate:"required" example:"securepassword123"`
}

type RegisterWithOfficeRequest struct {
	// User fields
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
	Email    string `json:"email" validate:"required,email"`

	// Office fields
	OfficeCode    string `json:"office_code" validate:"required,min=2,max=10"`
	OfficeName    string `json:"office_name" validate:"required,min=3,max=100"`
	OfficeAddress string `json:"office_address"`
	OfficeCity    string `json:"office_city"`
	OfficePhone   string `json:"office_phone"`
}

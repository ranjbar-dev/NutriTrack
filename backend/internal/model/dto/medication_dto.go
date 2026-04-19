package dto

// CreateMedicationRequest represents the request body for creating a medication.
type CreateMedicationRequest struct {
	Name        string  `json:"name" binding:"required,max=200"`
	GenericName *string `json:"generic_name" binding:"omitempty,max=200"`
	Form        string  `json:"form" binding:"required,oneof=tablet capsule syrup injection drop powder other"`
	DosageUnit  *string `json:"dosage_unit" binding:"omitempty,max=50"`
	Description *string `json:"description" binding:"omitempty,max=1000"`
}

// UpdateMedicationRequest represents the request body for updating a medication.
type UpdateMedicationRequest struct {
	Name        string  `json:"name" binding:"required,max=200"`
	GenericName *string `json:"generic_name" binding:"omitempty,max=200"`
	Form        string  `json:"form" binding:"required,oneof=tablet capsule syrup injection drop powder other"`
	DosageUnit  *string `json:"dosage_unit" binding:"omitempty,max=50"`
	Description *string `json:"description" binding:"omitempty,max=1000"`
}

// MedicationListQueryParams represents query parameters for listing medications.
type MedicationListQueryParams struct {
	Search   *string `form:"search"`
	IsActive *bool   `form:"is_active"`
	Page     int     `form:"page,default=1"`
	Limit    int     `form:"limit,default=20"`
}

// MedicationResponse represents a medication returned in API responses.
type MedicationResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	GenericName *string `json:"generic_name,omitempty"`
	Form        string  `json:"form"`
	DosageUnit  *string `json:"dosage_unit,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    bool    `json:"is_active"`
	CreatedBy   string  `json:"created_by"`
	CreatorName string  `json:"creator_name"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// MedicationListResponse represents a paginated list of medications.
type MedicationListResponse struct {
	Data    []MedicationResponse `json:"data"`
	Total   int64                `json:"total"`
	Page    int                  `json:"page"`
	Limit   int                  `json:"limit"`
	HasMore bool                 `json:"has_more"`
}

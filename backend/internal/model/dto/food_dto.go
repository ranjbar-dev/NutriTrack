package dto

type CreateFoodRequest struct {
	Name              string   `json:"name" binding:"required,max=200"`
	Description       *string  `json:"description" binding:"omitempty,max=1000"`
	Categories        []string `json:"categories" binding:"required,min=1,dive,oneof=breakfast lunch dinner snack fruit beverage supplement other"`
	Calories          float64  `json:"calories" binding:"gte=0,lte=9999.99"`
	ProteinG          float64  `json:"protein_g" binding:"gte=0,lte=999.99"`
	CarbsG            float64  `json:"carbs_g" binding:"gte=0,lte=999.99"`
	FatG              float64  `json:"fat_g" binding:"gte=0,lte=999.99"`
	FiberG            float64  `json:"fiber_g" binding:"gte=0,lte=999.99"`
	SugarG            float64  `json:"sugar_g" binding:"gte=0,lte=999.99"`
	SodiumMg          float64  `json:"sodium_mg" binding:"gte=0,lte=9999.99"`
	MeasurementUnit   string   `json:"measurement_unit" binding:"required,oneof=gram kg tablespoon teaspoon cup piece slice palm matchbox bowl ml liter"`
	MeasurementAmount float64  `json:"measurement_amount" binding:"required,gt=0"`
}

type UpdateFoodRequest struct {
	Name              string   `json:"name" binding:"required,max=200"`
	Description       *string  `json:"description" binding:"omitempty,max=1000"`
	Categories        []string `json:"categories" binding:"required,min=1,dive,oneof=breakfast lunch dinner snack fruit beverage supplement other"`
	Calories          float64  `json:"calories" binding:"gte=0,lte=9999.99"`
	ProteinG          float64  `json:"protein_g" binding:"gte=0,lte=999.99"`
	CarbsG            float64  `json:"carbs_g" binding:"gte=0,lte=999.99"`
	FatG              float64  `json:"fat_g" binding:"gte=0,lte=999.99"`
	FiberG            float64  `json:"fiber_g" binding:"gte=0,lte=999.99"`
	SugarG            float64  `json:"sugar_g" binding:"gte=0,lte=999.99"`
	SodiumMg          float64  `json:"sodium_mg" binding:"gte=0,lte=9999.99"`
	MeasurementUnit   string   `json:"measurement_unit" binding:"required,oneof=gram kg tablespoon teaspoon cup piece slice palm matchbox bowl ml liter"`
	MeasurementAmount float64  `json:"measurement_amount" binding:"required,gt=0"`
}

type FoodListQueryParams struct {
	Search   *string `form:"search"`
	Category *string `form:"category"`
	IsActive *bool   `form:"is_active"`
	Page     int     `form:"page,default=1"`
	Limit    int     `form:"limit,default=20"`
}

type FoodResponse struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Description       *string  `json:"description,omitempty"`
	Categories        []string `json:"categories"`
	Calories          float64  `json:"calories"`
	ProteinG          float64  `json:"protein_g"`
	CarbsG            float64  `json:"carbs_g"`
	FatG              float64  `json:"fat_g"`
	FiberG            float64  `json:"fiber_g"`
	SugarG            float64  `json:"sugar_g"`
	SodiumMg          float64  `json:"sodium_mg"`
	MeasurementUnit   string   `json:"measurement_unit"`
	MeasurementAmount float64  `json:"measurement_amount"`
	IsActive          bool     `json:"is_active"`
	CreatedBy         string   `json:"created_by"`
	CreatorName       string   `json:"creator_name"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
}

type FoodListResponse struct {
	Data    []FoodResponse `json:"data"`
	Total   int64          `json:"total"`
	Page    int            `json:"page"`
	Limit   int            `json:"limit"`
	HasMore bool           `json:"has_more"`
}

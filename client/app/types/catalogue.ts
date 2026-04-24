export interface FoodCategory {
  id: string
  name: string
  created_at?: string
}

export interface FoodItem {
  id: string
  name: string
  unit: string
  amount: number
  calories: number
  protein: number
  carbohydrate: number
  fat: number
  fiber: number
  sugar: number
  sodium: number
  is_active: boolean
  created_by: string | null
  categories: FoodCategory[]
  created_at: string
  updated_at: string
}

export interface MedicationItem {
  id: string
  name: string
  description: string | null
  unit: string
  is_active: boolean
  created_by: string | null
  created_at: string
  updated_at: string
}

export interface PaginationMeta {
  page: number
  page_size: number
  total: number
}

export interface PaginatedCatalogueResponse<T> {
  data: T[]
  meta: PaginationMeta
}

export interface FoodSearchQuery {
  q?: string
  category_id?: string
  page?: number
  page_size?: number
}

export interface MedicationSearchQuery {
  q?: string
  page?: number
  page_size?: number
}

export interface AdminFoodSearchQuery {
  q?: string
  category_id?: string
  page?: number
  page_size?: number
}

export interface AdminMedicationSearchQuery {
  q?: string
  page?: number
  page_size?: number
}

export interface AdminCreateFoodCategoryRequest {
  name: string
}

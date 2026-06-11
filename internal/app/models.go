package app

// Item represents an item in the database.
type Item struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name" gorm:"column:name"`
	Description string  `json:"description" gorm:"column:description"`
	Price       float64 `json:"price" gorm:"column:price"`
	Category    string  `json:"category" gorm:"column:category"`
}

// ItemCreate represents the request body for creating an item.
type ItemCreate struct {
	Name        string  `json:"name" binding:"required" validate:"required,min=1,max=100"`
	Description string  `json:"description" validate:"omitempty,max=500"`
	Price       float64 `json:"price" binding:"required" validate:"required,gt=0"`
	Category    string  `json:"category" validate:"omitempty,max=50"`
}

// ItemUpdate represents the request body for updating an item (all fields optional).
type ItemUpdate struct {
	Name        *string  `json:"name" validate:"omitempty,min=1,max=100"`
	Description *string  `json:"description" validate:"omitempty,max=500"`
	Price       *float64 `json:"price" validate:"omitempty,gt=0"`
	Category    *string  `json:"category" validate:"omitempty,max=50"`
}

package main

// Item represents an item in the database
type Item struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name" gorm:"column:name"`
	Description string  `json:"description" gorm:"column:description"`
	Price       float64 `json:"price" gorm:"column:price"`
	Category    string  `json:"category" gorm:"column:category"`
}

// ItemCreate represents the request body for creating an item
type ItemCreate struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Category    string  `json:"category"`
}

// ItemUpdate represents the request body for updating an item (all fields optional)
type ItemUpdate struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Category    *string  `json:"category"`
}

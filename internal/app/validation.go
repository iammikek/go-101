package app

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// InitValidator initializes the request validator.
func InitValidator() {
	validate = validator.New()
}

// ValidateItemCreate validates a create-item payload.
func ValidateItemCreate(item ItemCreate) error {
	return validate.Struct(item)
}

// ValidateItemUpdate validates an update-item payload.
func ValidateItemUpdate(item ItemUpdate) error {
	return validate.Struct(item)
}

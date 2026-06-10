package main

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func initValidator() {
	validate = validator.New()
}

func validateItemCreate(item ItemCreate) error {
	return validate.Struct(item)
}

func validateItemUpdate(item ItemUpdate) error {
	return validate.Struct(item)
}

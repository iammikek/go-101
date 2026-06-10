package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitValidator(t *testing.T) {
	initValidator()
	assert.NotNil(t, validate)
}

func TestValidateItemCreateValid(t *testing.T) {
	initValidator()

	item := ItemCreate{
		Name:        "Test Item",
		Description: "Test Description",
		Price:       9.99,
		Category:    "test",
	}

	err := validateItemCreate(item)
	assert.NoError(t, err)
}

func TestValidateItemCreateMissingName(t *testing.T) {
	initValidator()

	item := ItemCreate{
		Description: "Test Description",
		Price:       9.99,
		Category:    "test",
	}

	err := validateItemCreate(item)
	require.Error(t, err)
}

func TestValidateItemCreateMissingPrice(t *testing.T) {
	initValidator()

	item := ItemCreate{
		Name:        "Test Item",
		Description: "Test Description",
		Category:    "test",
	}

	err := validateItemCreate(item)
	require.Error(t, err)
}

func TestValidateItemUpdateValid(t *testing.T) {
	initValidator()

	item := ItemUpdate{
		Name: stringPtr("Updated Name"),
	}

	err := validateItemUpdate(item)
	assert.NoError(t, err)
}

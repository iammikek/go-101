package unit_test

import (
	"testing"

	"github.com/iammikek/go-101/internal/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func stringPtr(s string) *string {
	return &s
}

func TestValidateItemCreateValid(t *testing.T) {
	item := app.ItemCreate{
		Name:        "Test Item",
		Description: "Test Description",
		Price:       9.99,
		Category:    "test",
	}

	err := app.ValidateItemCreate(item)
	assert.NoError(t, err)
}

func TestValidateItemCreateMissingName(t *testing.T) {
	item := app.ItemCreate{
		Description: "Test Description",
		Price:       9.99,
		Category:    "test",
	}

	err := app.ValidateItemCreate(item)
	require.Error(t, err)
}

func TestValidateItemCreateMissingPrice(t *testing.T) {
	item := app.ItemCreate{
		Name:        "Test Item",
		Description: "Test Description",
		Category:    "test",
	}

	err := app.ValidateItemCreate(item)
	require.Error(t, err)
}

func TestValidateItemCreateInvalidPrice(t *testing.T) {
	item := app.ItemCreate{
		Name:  "Test Item",
		Price: -1.0,
	}

	err := app.ValidateItemCreate(item)
	require.Error(t, err)
}

func TestValidateItemUpdateValid(t *testing.T) {
	item := app.ItemUpdate{
		Name: stringPtr("Updated Name"),
	}

	err := app.ValidateItemUpdate(item)
	assert.NoError(t, err)
}

func TestValidateItemUpdateInvalidPrice(t *testing.T) {
	price := -5.0
	item := app.ItemUpdate{
		Price: &price,
	}

	err := app.ValidateItemUpdate(item)
	require.Error(t, err)
}

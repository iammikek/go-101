package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	// Test that logger can be initialized without panic
	assert.NotNil(t, logger)
}

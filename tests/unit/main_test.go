package unit_test

import (
	"os"
	"testing"

	"github.com/iammikek/go-101/internal/app"
)

func TestMain(m *testing.M) {
	app.InitValidator()
	os.Exit(m.Run())
}

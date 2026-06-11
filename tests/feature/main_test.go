package feature_test

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/iammikek/go-101/internal/app"
	"github.com/iammikek/go-101/tests/testcase"
)

var featureApp *app.Application

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	tmpFile, err := os.CreateTemp("", "go-101-feature-*.db")
	if err != nil {
		panic(err)
	}
	dbPath := tmpFile.Name()
	tmpFile.Close()

	featureApp, err = app.New(app.Config{
		DatabaseURL: dbPath,
		APIKey:      testcase.DefaultAPIKey,
	})
	if err != nil {
		panic(err)
	}

	code := m.Run()
	_ = os.Remove(dbPath)
	os.Exit(code)
}

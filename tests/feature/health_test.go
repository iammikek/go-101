package feature_test

import (
	"net/http"
	"testing"

	"github.com/iammikek/go-101/tests/testcase"
)

func TestRoot(t *testing.T) {
	tc := testcase.NewFeature(t, featureApp)

	rec := tc.Get("/")
	rec.AssertStatus(http.StatusOK)
	rec.AssertJSON(`{"message":"Hello from Go!"}`)
}

func TestHealth(t *testing.T) {
	tc := testcase.NewFeature(t, featureApp)

	rec := tc.Get("/health")
	rec.AssertStatus(http.StatusOK)
	rec.AssertJSON(`{"status":"ok"}`)
}

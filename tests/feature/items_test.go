package feature_test

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iammikek/go-101/tests/testcase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListItemsEmpty(t *testing.T) {
	tc := testcase.NewFeature(t, featureApp)

	rec := tc.Get("/items")
	rec.AssertStatus(http.StatusOK)
	rec.AssertJSON(`[]`)
}

func TestListItemsWithPagination(t *testing.T) {
	tc := testcase.NewFeature(t, featureApp)

	for _, item := range []struct {
		name  string
		price float64
	}{
		{"A", 1.0}, {"B", 2.0}, {"C", 3.0},
	} {
		rec := tc.Post("/items", `{"name":"`+item.name+`","price":`+formatPrice(item.price)+`}`)
		rec.AssertStatus(http.StatusCreated)
	}

	rec := tc.Get("/items?skip=1&limit=2")
	rec.AssertStatus(http.StatusOK)

	data := rec.JSONArray()
	require.Len(t, data, 2)
	assert.Equal(t, "B", data[0]["name"])
	assert.Equal(t, "C", data[1]["name"])
}

func TestCreateItem(t *testing.T) {
	tc := testcase.NewFeature(t, featureApp)

	rec := tc.Post("/items", `{"name":"Widget","description":"A nice widget","price":9.99}`)
	rec.AssertStatus(http.StatusCreated)

	data := rec.JSON()
	assert.GreaterOrEqual(t, int(data["id"].(float64)), 1)
	assert.Equal(t, "Widget", data["name"])
	assert.Equal(t, "A nice widget", data["description"])
	assert.Equal(t, 9.99, data["price"])
}

func TestCreateItemOptionalDescription(t *testing.T) {
	tc := testcase.NewFeature(t, featureApp)

	rec := tc.Post("/items", `{"name":"Thing","price":5.0}`)
	rec.AssertStatus(http.StatusCreated)
	assert.Equal(t, "", rec.JSON()["description"])
}

func TestCreateItemWithCategory(t *testing.T) {
	tc := testcase.NewFeature(t, featureApp)

	rec := tc.Post("/items", `{"name":"Gadget","price":15.0,"category":"Electronics"}`)
	rec.AssertStatus(http.StatusCreated)
	assert.Equal(t, "Electronics", rec.JSON()["category"])
}

func TestUpdateItemCategory(t *testing.T) {
	tc := testcase.NewFeature(t, featureApp)

	created := tc.Post("/items", `{"name":"Item","price":10.0}`)
	created.AssertStatus(http.StatusCreated)
	itemID := int(created.JSON()["id"].(float64))

	rec := tc.Patch("/items/"+strconv.Itoa(itemID), `{"category":"Tools"}`)
	rec.AssertStatus(http.StatusOK)
	assert.Equal(t, "Tools", rec.JSON()["category"])
}

func TestGetItem(t *testing.T) {
	tc := testcase.NewFeature(t, featureApp)

	created := tc.Post("/items", `{"name":"Widget","price":9.99}`)
	created.AssertStatus(http.StatusCreated)
	itemID := int(created.JSON()["id"].(float64))

	rec := tc.Get("/items/" + strconv.Itoa(itemID))
	rec.AssertStatus(http.StatusOK)
	assert.Equal(t, "Widget", rec.JSON()["name"])
}

func TestGetItemNotFound(t *testing.T) {
	tc := testcase.NewFeature(t, featureApp)

	rec := tc.Get("/items/99")
	rec.AssertStatus(http.StatusNotFound)
	rec.AssertJSON(`{"detail":"Item not found"}`)
}

func TestCreateItemInvalidBody(t *testing.T) {
	tc := testcase.NewFeature(t, featureApp)

	rec := tc.Post("/items", `{"name":"No price"}`)
	rec.AssertStatus(http.StatusUnprocessableEntity)
}

func TestUpdateItemPartial(t *testing.T) {
	tc := testcase.NewFeature(t, featureApp)

	created := tc.Post("/items", `{"name":"Widget","description":"Original","price":10.0}`)
	created.AssertStatus(http.StatusCreated)
	itemID := int(created.JSON()["id"].(float64))

	rec := tc.Patch("/items/"+strconv.Itoa(itemID), `{"price":5.99}`)
	rec.AssertStatus(http.StatusOK)

	data := rec.JSON()
	assert.Equal(t, "Widget", data["name"])
	assert.Equal(t, "Original", data["description"])
	assert.Equal(t, 5.99, data["price"])
}

func TestUpdateItemFull(t *testing.T) {
	tc := testcase.NewFeature(t, featureApp)

	created := tc.Post("/items", `{"name":"Old","price":1.0}`)
	created.AssertStatus(http.StatusCreated)
	itemID := int(created.JSON()["id"].(float64))

	rec := tc.Patch("/items/"+strconv.Itoa(itemID), `{"name":"New","description":"Updated","price":2.5}`)
	rec.AssertStatus(http.StatusOK)

	data := rec.JSON()
	assert.Equal(t, float64(itemID), data["id"])
	assert.Equal(t, "New", data["name"])
	assert.Equal(t, "Updated", data["description"])
	assert.Equal(t, 2.5, data["price"])
	_, hasCategory := data["category"]
	assert.True(t, hasCategory)
}

func TestUpdateItemNotFound(t *testing.T) {
	tc := testcase.NewFeature(t, featureApp)

	rec := tc.Patch("/items/99", `{"name":"Nope"}`)
	rec.AssertStatus(http.StatusNotFound)
	rec.AssertJSON(`{"detail":"Item not found"}`)
}

func TestDeleteItemWithAPIKey(t *testing.T) {
	tc := testcase.NewFeature(t, featureApp)

	created := tc.Post("/items", `{"name":"To Delete","price":1.0}`)
	created.AssertStatus(http.StatusCreated)
	itemID := int(created.JSON()["id"].(float64))

	del := tc.Delete("/items/"+strconv.Itoa(itemID), tc.WithAPIKey())
	del.AssertStatus(http.StatusNoContent)

	rec := tc.Get("/items/" + strconv.Itoa(itemID))
	rec.AssertStatus(http.StatusNotFound)
}

func TestDeleteItemWithoutAPIKey(t *testing.T) {
	tc := testcase.NewFeature(t, featureApp)

	created := tc.Post("/items", `{"name":"Test","price":1.0}`)
	created.AssertStatus(http.StatusCreated)
	itemID := int(created.JSON()["id"].(float64))

	rec := tc.Delete("/items/"+strconv.Itoa(itemID), nil)
	rec.AssertStatus(http.StatusUnauthorized)
	rec.AssertJSON(`{"detail":"Invalid or missing API key"}`)
}

func TestDeleteItemInvalidAPIKey(t *testing.T) {
	tc := testcase.NewFeature(t, featureApp)

	created := tc.Post("/items", `{"name":"Test","price":1.0}`)
	created.AssertStatus(http.StatusCreated)
	itemID := int(created.JSON()["id"].(float64))

	rec := tc.Delete("/items/"+strconv.Itoa(itemID), map[string]string{"X-API-Key": "wrong-key"})
	rec.AssertStatus(http.StatusUnauthorized)
	rec.AssertJSON(`{"detail":"Invalid or missing API key"}`)
}

func TestDeleteItemNotFound(t *testing.T) {
	tc := testcase.NewFeature(t, featureApp)

	rec := tc.Delete("/items/99", tc.WithAPIKey())
	rec.AssertStatus(http.StatusNotFound)
	rec.AssertJSON(`{"detail":"Item not found"}`)
}

func TestGetItemsStatsEmpty(t *testing.T) {
	tc := testcase.NewFeature(t, featureApp)

	rec := tc.Get("/items/stats/summary")
	rec.AssertStatus(http.StatusOK)
	rec.AssertJSON(`{"average_price":0,"max_price":null,"min_price":null,"total_items":0}`)
}

func TestGetItemsStats(t *testing.T) {
	tc := testcase.NewFeature(t, featureApp)

	for _, body := range []string{
		`{"name":"A","price":10.0}`,
		`{"name":"B","price":20.0}`,
		`{"name":"C","price":30.0}`,
	} {
		rec := tc.Post("/items", body)
		rec.AssertStatus(http.StatusCreated)
	}

	rec := tc.Get("/items/stats/summary")
	rec.AssertStatus(http.StatusOK)

	data := rec.JSON()
	assert.Equal(t, float64(3), data["total_items"])
	assert.Equal(t, 20.0, data["average_price"])
	assert.Equal(t, 10.0, data["min_price"])
	assert.Equal(t, 30.0, data["max_price"])
}

func formatPrice(price float64) string {
	return strconv.FormatFloat(price, 'f', -1, 64)
}

package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "app.db"
	}

	var err error
	db, err = gorm.Open(sqlite.Open(databaseURL), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	// Auto migrate the Item model
	db.AutoMigrate(&Item{})
}

func main() {
	router := gin.Default()

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from Go!",
		})
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// List items with pagination
	router.GET("/items", listItems)

	// Get single item by ID
	router.GET("/items/:item_id", getItem)

	// Create item
	router.POST("/items", createItem)

	// Update item (partial)
	router.PATCH("/items/:item_id", updateItem)

	// Delete item (requires API key)
	router.DELETE("/items/:item_id", authMiddleware(), deleteItem)

	// Stats endpoint
	router.GET("/items/stats/summary", getItemsStats)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	if err := router.Run(":" + port); err != nil {
		panic(fmt.Sprintf("failed to start server: %v", err))
	}
}

// Handler functions

func listItems(c *gin.Context) {
	var items []Item
	skip := c.DefaultQuery("skip", "0")
	limit := c.DefaultQuery("limit", "10")

	db.Offset(skip).Limit(limit).Find(&items)
	c.JSON(200, items)
}

func getItem(c *gin.Context) {
	itemID := c.Param("item_id")
	var item Item

	result := db.First(&item, itemID)
	if result.Error != nil {
		c.JSON(404, gin.H{"detail": "Item not found"})
		return
	}

	c.JSON(200, item)
}

func createItem(c *gin.Context) {
	var req ItemCreate

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"detail": "Invalid request body"})
		return
	}

	item := Item{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
	}

	if err := db.Create(&item).Error; err != nil {
		c.JSON(500, gin.H{"detail": "Failed to create item"})
		return
	}

	c.JSON(201, item)
}

func updateItem(c *gin.Context) {
	itemID := c.Param("item_id")
	var item Item

	if err := db.First(&item, itemID).Error; err != nil {
		c.JSON(404, gin.H{"detail": "Item not found"})
		return
	}

	var req ItemUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"detail": "Invalid request body"})
		return
	}

	// Update only provided fields
	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}

	if err := db.Model(&item).Updates(updates).Error; err != nil {
		c.JSON(500, gin.H{"detail": "Failed to update item"})
		return
	}

	c.JSON(200, item)
}

func deleteItem(c *gin.Context) {
	itemID := c.Param("item_id")
	var item Item

	if err := db.First(&item, itemID).Error; err != nil {
		c.JSON(404, gin.H{"detail": "Item not found"})
		return
	}

	if err := db.Delete(&item).Error; err != nil {
		c.JSON(500, gin.H{"detail": "Failed to delete item"})
		return
	}

	c.Status(204)
}

func getItemsStats(c *gin.Context) {
	var count int64
	var avgPrice float64
	var minPrice float64
	var maxPrice float64

	db.Model(&Item{}).Count(&count)

	if count == 0 {
		c.JSON(200, gin.H{
			"total_items":   0,
			"average_price": 0.0,
			"min_price":     nil,
			"max_price":     nil,
		})
		return
	}

	db.Model(&Item{}).Select("AVG(price)").Row().Scan(&avgPrice)
	db.Model(&Item{}).Select("MIN(price)").Row().Scan(&minPrice)
	db.Model(&Item{}).Select("MAX(price)").Row().Scan(&maxPrice)

	c.JSON(200, gin.H{
		"total_items":   count,
		"average_price": avgPrice,
		"min_price":     minPrice,
		"max_price":     maxPrice,
	})
}

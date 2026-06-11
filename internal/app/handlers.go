package app

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func listItems(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var items []Item
		skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

		db.Offset(skip).Limit(limit).Find(&items)
		c.JSON(200, items)
	}
}

func getItem(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		itemID := c.Param("item_id")
		var item Item

		result := db.First(&item, itemID)
		if result.Error != nil {
			c.JSON(404, gin.H{"detail": "Item not found"})
			return
		}

		c.JSON(200, item)
	}
}

func createItem(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ItemCreate

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(422, gin.H{"detail": "Invalid request body"})
			return
		}

		if err := ValidateItemCreate(req); err != nil {
			c.JSON(422, gin.H{"detail": err.Error()})
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
}

func updateItem(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		itemID := c.Param("item_id")
		var item Item

		if err := db.First(&item, itemID).Error; err != nil {
			c.JSON(404, gin.H{"detail": "Item not found"})
			return
		}

		var req ItemUpdate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(422, gin.H{"detail": "Invalid request body"})
			return
		}

		if err := ValidateItemUpdate(req); err != nil {
			c.JSON(422, gin.H{"detail": err.Error()})
			return
		}

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

		db.First(&item, itemID)
		c.JSON(200, item)
	}
}

func deleteItem(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
}

func getItemsStats(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
}

package app

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Config holds application configuration.
type Config struct {
	DatabaseURL string
	APIKey      string
	LogLevel    string
}

// Application is the bootstrapped HTTP application.
type Application struct {
	DB     *gorm.DB
	Router *gin.Engine
	Logger *slog.Logger
	Config Config
}

// New boots the application with database, migrations, and routes.
func New(cfg Config) (*Application, error) {
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = "app.db"
	}
	if cfg.APIKey == "" {
		cfg.APIKey = "dev-key-123"
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = os.Getenv("LOG_LEVEL")
	}

	log := newLogger(cfg.LogLevel)

	db, err := gorm.Open(sqlite.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if err := runMigrations(cfg.DatabaseURL); err != nil {
		return nil, err
	}

	InitValidator()

	router := setupRouter(db, cfg.APIKey, log)

	return &Application{
		DB:     db,
		Router: router,
		Logger: log,
		Config: cfg,
	}, nil
}

// ResetDatabase clears all items between tests.
func (a *Application) ResetDatabase() error {
	if err := a.DB.Exec("DELETE FROM items").Error; err != nil {
		return err
	}
	return a.DB.Exec("DELETE FROM sqlite_sequence WHERE name='items'").Error
}

func setupRouter(db *gorm.DB, apiKey string, log *slog.Logger) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestLoggingMiddleware(log))

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from Go!"})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.GET("/items", listItems(db))
	router.GET("/items/:item_id", getItem(db))
	router.POST("/items", createItem(db))
	router.PATCH("/items/:item_id", updateItem(db))
	router.DELETE("/items/:item_id", AuthMiddleware(apiKey), deleteItem(db))
	router.GET("/items/stats/summary", getItemsStats(db))

	return router
}

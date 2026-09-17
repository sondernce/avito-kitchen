package postgres

import (
    "fmt"
    "log"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "avito-kitchen/main-service/internal/models"
)

type Config struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

func NewConnection(cfg Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    return db, nil
}

func AutoMigrate(db *gorm.DB) error {
    log.Println("Running database migrations...")
    
    err := db.AutoMigrate(
        &models.Restaurant{},
        &models.Dish{},
        &models.Order{},
        &models.OrderItem{},
        &models.User{},
    )
    if err != nil {
        return fmt.Errorf("failed to run migrations: %w", err)
    }

    log.Println("Migrations completed successfully")
    return nil
}
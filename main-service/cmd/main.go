package main

import (
    "log"
    "os"

    "avito-kitchen/main-service/internal/handlers"
    "avito-kitchen/main-service/internal/repository/postgres"
    "avito-kitchen/main-service/internal/service"
)

func main() {
    // Конфигурация БД
    dbConfig := postgres.Config{
        Host:     getEnv("DB_HOST", "localhost"),
        Port:     getEnv("DB_PORT", "5432"),
        User:     getEnv("DB_USER", "postgres"),
        Password: getEnv("DB_PASSWORD", "postgres"),
        DBName:   getEnv("DB_NAME", "avito_kitchen"),
        SSLMode:  getEnv("DB_SSLMODE", "disable"),
    }

    // Подключение к БД
    db, err := postgres.NewConnection(dbConfig)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Миграции
    if err := postgres.AutoMigrate(db); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }

    // Инициализация репозиториев (PostgreSQL)
    restaurantRepo := postgres.NewRestaurantRepository(db)
    dishRepo := postgres.NewDishRepository(db)
    orderRepo := postgres.NewOrderRepository(db)

    // Инициализация сервисов
    restaurantService := service.NewRestaurantService(restaurantRepo)
    dishService := service.NewDishService(dishRepo)
    webhookService := service.NewWebhookService()
    orderService := service.NewOrderService(orderRepo, webhookService, restaurantRepo, dishRepo)

    // Инициализация обработчиков
    restaurantHandlers := handlers.NewRestaurantHandlers(restaurantService)
    dishHandlers := handlers.NewDishHandlers(dishService)
    orderHandlers := handlers.NewOrderHandlers(orderService)

    // Создание сервера
    server := handlers.NewServer(restaurantHandlers, dishHandlers, orderHandlers)

    // Запуск сервера
    addr := ":8080"
    log.Printf("Starting Avito.Kitchen Main Service on %s", addr)
    if err := server.Start(addr); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
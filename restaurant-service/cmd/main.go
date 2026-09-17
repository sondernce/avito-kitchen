package main

import (
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"

    "avito-kitchen/restaurant-service/internal/client"
    "avito-kitchen/restaurant-service/internal/handlers"
)

func main() {
    // Конфигурация
    mainServiceURL := getEnv("MAIN_SERVICE_URL", "http://localhost:8080")
    restaurantID := 1 // ID нашего ресторана в main-service

    log.Printf("Starting Restaurant Service for restaurant #%d", restaurantID)
    log.Printf("Main Service URL: %s", mainServiceURL)

    // Инициализация клиента для main-service
    mainServiceClient := client.NewMainServiceClient(mainServiceURL)

    // Инициализация обработчиков
    webhookHandler := handlers.NewWebhookHandler(mainServiceClient, restaurantID)

    // Настройка маршрутов
    router := mux.NewRouter()
    router.HandleFunc("/webhook/order", webhookHandler.HandleOrderWebhook).Methods("POST")
    
    // Health check
    router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    }).Methods("GET")

    // Запуск сервера
    addr := ":8081"
    log.Printf("Restaurant Service listening on %s", addr)
    if err := http.ListenAndServe(addr, router); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
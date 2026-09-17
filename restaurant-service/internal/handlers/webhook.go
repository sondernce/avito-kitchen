package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "time"

    "avito-kitchen/restaurant-service/internal/client"
    "avito-kitchen/restaurant-service/internal/models"
)

type WebhookHandler struct {
    mainServiceClient *client.MainServiceClient
    restaurantID      int
}

func NewWebhookHandler(mainServiceClient *client.MainServiceClient, restaurantID int) *WebhookHandler {
    return &WebhookHandler{
        mainServiceClient: mainServiceClient,
        restaurantID:      restaurantID,
    }
}

// HandleOrderWebhook принимает уведомление о новом заказе от main-service
// POST /webhook/order
func (h *WebhookHandler) HandleOrderWebhook(w http.ResponseWriter, r *http.Request) {
    var webhook models.OrderWebhook
    if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
        log.Printf("Failed to decode webhook: %v", err)
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    log.Printf("Received new order #%d for restaurant #%d", webhook.OrderID, webhook.RestaurantID)
    log.Printf("Delivery address: %s, Phone: %s", webhook.DeliveryAddress, webhook.Phone)
    log.Printf("Total price: %.2f", webhook.TotalPrice)

    // Проверяем, что заказ для нашего ресторана
    if webhook.RestaurantID != h.restaurantID {
        log.Printf("Order #%d is for restaurant #%d, but we are #%d", 
            webhook.OrderID, webhook.RestaurantID, h.restaurantID)
        http.Error(w, "wrong restaurant", http.StatusBadRequest)
        return
    }

    // Автоматически подтверждаем заказ (заглшушка)
    go h.processOrder(webhook.OrderID)

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "received",
        "message": "Order received and being processed",
    })
}

// processOrder обрабатывает заказ и обновляет статусы
func (h *WebhookHandler) processOrder(orderID int) {
    // 1. Подтверждаем заказ
    if err := h.mainServiceClient.UpdateOrderStatus(orderID, "confirmed"); err != nil {
        log.Printf("Failed to confirm order #%d: %v", orderID, err)
        return
    }
    log.Printf("Order #%d confirmed", orderID)

    // 2. Начинаем готовить
    time.Sleep(5 * time.Second)
    if err := h.mainServiceClient.UpdateOrderStatus(orderID, "preparing"); err != nil {
        log.Printf("Failed to set preparing for order #%d: %v", orderID, err)
        return
    }
    log.Printf("Order #%d is being prepared", orderID)

    // 3. Отправляем доставку
    time.Sleep(10 * time.Second)
    if err := h.mainServiceClient.UpdateOrderStatus(orderID, "delivering"); err != nil {
        log.Printf("Failed to set delivering for order #%d: %v", orderID, err)
        return
    }
    log.Printf("Order #%d is out for delivery", orderID)

    // 4. Завершаем заказ
    time.Sleep(15 * time.Second)
    if err := h.mainServiceClient.UpdateOrderStatus(orderID, "completed"); err != nil {
        log.Printf("Failed to complete order #%d: %v", orderID, err)
        return
    }
    log.Printf("Order #%d completed", orderID)
}
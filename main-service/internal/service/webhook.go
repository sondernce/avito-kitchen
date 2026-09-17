package service

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    "avito-kitchen/main-service/internal/models"
)

type WebhookService struct {
    httpClient *http.Client
}

func NewWebhookService() *WebhookService {
    return &WebhookService{
        httpClient: &http.Client{
            Timeout: 5 * time.Second,
        },
    }
}

// NotifyRestaurant отправляет webhook уведомление ресторану о новом заказе
func (s *WebhookService) NotifyRestaurant(webhookURL string, order *models.Order, items []models.OrderItem) error {
    if webhookURL == "" {
        log.Printf("No webhook URL for restaurant #%d, skipping notification", order.RestaurantID)
        return nil
    }

    payload := map[string]interface{}{
        "order_id":         order.ID,
        "restaurant_id":    order.RestaurantID,
        "status":           order.Status,
        "total_price":      order.TotalPrice,
        "delivery_address": order.DeliveryAddress,
        "phone":            order.Phone,
        "created_at":       order.CreatedAt,
        "items":            items,
    }

    jsonData, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("failed to marshal webhook payload: %w", err)
    }

    req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("failed to create webhook request: %w", err)
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := s.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send webhook: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("webhook returned status %d", resp.StatusCode)
    }

    log.Printf("Webhook sent successfully to %s for order #%d", webhookURL, order.ID)
    return nil
}
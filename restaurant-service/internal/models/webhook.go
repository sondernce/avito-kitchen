package models

import "time"

type OrderWebhook struct {
    OrderID       int       `json:"order_id"`
    RestaurantID  int       `json:"restaurant_id"`
    Status        string    `json:"status"`
    TotalPrice    float64   `json:"total_price"`
    DeliveryAddress string  `json:"delivery_address"`
    Phone         string    `json:"phone"`
    CreatedAt     time.Time `json:"created_at"`
    Items         []OrderItem `json:"items"`
}

type OrderItem struct {
    DishID   int     `json:"dish_id"`
    Quantity int     `json:"quantity"`
    Price    float64 `json:"price"`
}

type StatusUpdateRequest struct {
    Status string `json:"status"`
}
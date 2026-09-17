package models

import "time"

type OrderStatus string

const (
    StatusCreated    OrderStatus = "created"
    StatusConfirmed  OrderStatus = "confirmed"
    StatusPreparing  OrderStatus = "preparing"
    StatusDelivering OrderStatus = "delivering"
    StatusCompleted  OrderStatus = "completed"
    StatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
    ID              int         `gorm:"primaryKey" json:"id"`
    UserID          int         `gorm:"not null;index" json:"user_id"`
    RestaurantID    int         `gorm:"not null;index" json:"restaurant_id"`
    Status          OrderStatus `gorm:"type:varchar(50);not null;default:'created'" json:"status"`
    TotalPrice      float64     `gorm:"type:decimal(10,2);not null" json:"total_price"`
    DeliveryAddress string      `gorm:"not null" json:"delivery_address"`
    Phone           string      `json:"phone"`
    CreatedAt       time.Time   `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt       time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
    Items           []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

func NewOrder(userID, restaurantID int, totalPrice float64, deliveryAddress, phone string) *Order {
    return &Order{
        UserID:          userID,
        RestaurantID:    restaurantID,
        Status:          StatusCreated,
        TotalPrice:      totalPrice,
        DeliveryAddress: deliveryAddress,
        Phone:           phone,
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
    }
}

func (o *Order) SetStatus(status OrderStatus) {
    o.Status = status
    o.UpdatedAt = time.Now()
}

func (o *Order) CanTransitionTo(newStatus OrderStatus) bool {
    switch o.Status {
    case StatusCreated:
        return newStatus == StatusConfirmed || newStatus == StatusCancelled
    case StatusConfirmed:
        return newStatus == StatusPreparing || newStatus == StatusCancelled
    case StatusPreparing:
        return newStatus == StatusDelivering
    case StatusDelivering:
        return newStatus == StatusCompleted
    default:
        return false
    }
}

type OrderItem struct {
    ID       int     `gorm:"primaryKey" json:"id"`
    OrderID  int     `gorm:"not null;index" json:"order_id"`
    DishID   int     `gorm:"not null" json:"dish_id"`
    Quantity int     `gorm:"not null" json:"quantity"`
    Price    float64 `gorm:"type:decimal(10,2);not null" json:"price"`
}
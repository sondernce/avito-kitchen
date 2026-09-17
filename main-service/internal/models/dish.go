package models

import "time"

type Dish struct {
    ID           int       `gorm:"primaryKey" json:"id"`
    RestaurantID int       `gorm:"not null;index" json:"restaurant_id"`
    Name         string    `gorm:"not null" json:"name"`
    Description  string    `json:"description"`
    Price        float64   `gorm:"type:decimal(10,2);not null" json:"price"`
    Category     string    `json:"category"`
    IsAvailable  bool      `gorm:"default:true" json:"is_available"`
    CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func NewDish(restaurantID int, name, description string, price float64, category string) *Dish {
    return &Dish{
        RestaurantID: restaurantID,
        Name:         name,
        Description:  description,
        Price:        price,
        Category:     category,
        IsAvailable:  true,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
}

func (d *Dish) SetAvailable(available bool) {
    d.IsAvailable = available
    d.UpdatedAt = time.Now()
}

func (d *Dish) UpdatePrice(price float64) {
    d.Price = price
    d.UpdatedAt = time.Now()
}
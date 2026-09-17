package models

import "time"

type Restaurant struct {
    ID          int       `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"not null" json:"name"`
    Description string    `json:"description"`
    Address     string    `gorm:"not null" json:"address"`
    Phone       string    `json:"phone"`
    IsActive    bool      `gorm:"default:true" json:"is_active"`
    WebhookURL  string    `json:"webhook_url"`
    CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func NewRestaurant(name, description, address, phone string) *Restaurant {
    return &Restaurant{
        Name:        name,
        Description: description,
        Address:     address,
        Phone:       phone,
        IsActive:    true,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
}

func (r *Restaurant) Activate() {
    r.IsActive = true
    r.UpdatedAt = time.Now()
}

func (r *Restaurant) Deactivate() {
    r.IsActive = false
    r.UpdatedAt = time.Now()
}
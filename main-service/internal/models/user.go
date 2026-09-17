package models

import "time"

type User struct {
    ID        int       `gorm:"primaryKey" json:"id"`
    Name      string    `gorm:"not null" json:"name"`
    Phone     string    `gorm:"uniqueIndex" json:"phone"`
    Address   string    `json:"address"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func NewUser(name, phone, address string) *User {
    return &User{
        Name:      name,
        Phone:     phone,
        Address:   address,
        CreatedAt: time.Now(),
    }
}
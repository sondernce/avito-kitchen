package handlers

import (
    "encoding/json"
    "errors"
    "time"
)

type ErrorDTO struct {
    Message string    `json:"message"`
    Time    time.Time `json:"time"`
}

func (e ErrorDTO) ToString() string {
    b, err := json.MarshalIndent(e, "", "   ")
    if err != nil {
        panic(err)
    }
    return string(b)
}

type CreateRestaurantDTO struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Address     string `json:"address"`
    Phone       string `json:"phone"`
}

func (d CreateRestaurantDTO) Validate() error {
    if d.Name == "" {
        return errors.New("name is required")
    }
    if d.Address == "" {
        return errors.New("address is required")
    }
    return nil
}

type CreateDishDTO struct {
    RestaurantID int     `json:"restaurant_id"`
    Name         string  `json:"name"`
    Description  string  `json:"description"`
    Price        float64 `json:"price"`
    Category     string  `json:"category"`
}

func (d CreateDishDTO) Validate() error {
    if d.Name == "" {
        return errors.New("name is required")
    }
    if d.Price <= 0 {
        return errors.New("price must be positive")
    }
    return nil
}

type CreateOrderDTO struct {
    UserID          int     `json:"user_id"`
    RestaurantID    int     `json:"restaurant_id"`
    DeliveryAddress string  `json:"delivery_address"`
    Phone           string  `json:"phone"`
    Items           []ItemDTO `json:"items"`
}

type ItemDTO struct {
    DishID   int `json:"dish_id"`
    Quantity int `json:"quantity"`
}

func (d CreateOrderDTO) Validate() error {
    if d.UserID <= 0 {
        return errors.New("user_id is required")
    }
    if d.RestaurantID <= 0 {
        return errors.New("restaurant_id is required")
    }
    if d.DeliveryAddress == "" {
        return errors.New("delivery_address is required")
    }
    if len(d.Items) == 0 {
        return errors.New("items cannot be empty")
    }
    return nil
}

type UpdateStatusDTO struct {
    Status string `json:"status"`
}

func (d UpdateStatusDTO) Validate() error {
    if d.Status == "" {
        return errors.New("status is required")
    }
    return nil
}
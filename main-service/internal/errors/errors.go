package errors

import "errors"

var (
    // Restaurant errors
    ErrRestaurantNotFound = errors.New("restaurant not found")
    ErrRestaurantExists   = errors.New("restaurant already exists")

    // Dish errors
    ErrDishNotFound = errors.New("dish not found")
    ErrDishExists   = errors.New("dish already exists")

    // Order errors
    ErrOrderNotFound          = errors.New("order not found")
    ErrOrderExists            = errors.New("order already exists")
    ErrInvalidStatusTransition = errors.New("invalid status transition")

    // User errors
    ErrUserNotFound = errors.New("user not found")
    ErrUserExists   = errors.New("user already exists")
)
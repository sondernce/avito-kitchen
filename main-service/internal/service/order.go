package service

import (
    "log"

    "avito-kitchen/main-service/internal/errors"
    "avito-kitchen/main-service/internal/models"
    "avito-kitchen/main-service/internal/repository/postgres"
)

type OrderService struct {
    repo           *postgres.OrderRepository
    webhookSvc     *WebhookService
    restaurantRepo *postgres.RestaurantRepository
    dishRepo       *postgres.DishRepository
}

func NewOrderService(
    repo *postgres.OrderRepository,
    webhookSvc *WebhookService,
    restaurantRepo *postgres.RestaurantRepository,
    dishRepo *postgres.DishRepository,
) *OrderService {
    return &OrderService{
        repo:           repo,
        webhookSvc:     webhookSvc,
        restaurantRepo: restaurantRepo,
        dishRepo:       dishRepo,
    }
}

func (s *OrderService) CreateOrder(userID, restaurantID int, items []models.OrderItem, deliveryAddress, phone string) (*models.Order, error) {
    totalPrice := 0.0
    for i, item := range items {
        dish, err := s.dishRepo.GetByID(item.DishID)
        if err != nil {
            return nil, err
        }
        item.Price = dish.Price * float64(item.Quantity)
        items[i] = item
        totalPrice += item.Price
    }

    order := models.NewOrder(userID, restaurantID, totalPrice, deliveryAddress, phone)
    order.Items = items
    if err := s.repo.Create(order); err != nil {
        return nil, err
    }
    go s.sendWebhookToRestaurant(order)
    return order, nil
}

func (s *OrderService) sendWebhookToRestaurant(order *models.Order) {
    restaurant, err := s.restaurantRepo.GetByID(order.RestaurantID)
    if err != nil {
        log.Printf("Failed to get restaurant #%d for webhook: %v", order.RestaurantID, err)
        return
    }
    if restaurant.WebhookURL == "" {
        log.Printf("Restaurant #%d has no webhook URL, skipping", order.RestaurantID)
        return
    }
    if err := s.webhookSvc.NotifyRestaurant(restaurant.WebhookURL, order, order.Items); err != nil {
        log.Printf("Failed to send webhook for order #%d: %v", order.ID, err)
    }
}

func (s *OrderService) GetOrder(id int) (*models.Order, error) {
    return s.repo.GetByID(id)
}

func (s *OrderService) GetUserOrders(userID int) ([]*models.Order, error) {
    return s.repo.GetByUserID(userID)
}

func (s *OrderService) GetRestaurantOrders(restaurantID int) ([]*models.Order, error) {
    return s.repo.GetByRestaurantID(restaurantID)
}

func (s *OrderService) UpdateOrderStatus(id int, newStatus models.OrderStatus) error {
    order, err := s.repo.GetByID(id)
    if err != nil {
        return err
    }
    if !order.CanTransitionTo(newStatus) {
        return errors.ErrInvalidStatusTransition
    }
    order.SetStatus(newStatus)
    return s.repo.Update(order)
}

func (s *OrderService) CancelOrder(id int) error {
    return s.UpdateOrderStatus(id, models.StatusCancelled)
}
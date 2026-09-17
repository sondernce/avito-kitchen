package service

import (
    "avito-kitchen/main-service/internal/models"
    "avito-kitchen/main-service/internal/repository/postgres"
)

type DishService struct {
    repo *postgres.DishRepository
}

func NewDishService(repo *postgres.DishRepository) *DishService {
    return &DishService{repo: repo}
}

func (s *DishService) CreateDish(restaurantID int, name, description string, price float64, category string) (*models.Dish, error) {
    dish := models.NewDish(restaurantID, name, description, price, category)
    if err := s.repo.Create(dish); err != nil {
        return nil, err
    }
    return dish, nil
}

func (s *DishService) GetDish(id int) (*models.Dish, error) {
    return s.repo.GetByID(id)
}

func (s *DishService) GetMenu(restaurantID int) ([]*models.Dish, error) {
    return s.repo.GetByRestaurantID(restaurantID)
}

func (s *DishService) UpdateDishAvailability(id int, available bool) error {
    dish, err := s.repo.GetByID(id)
    if err != nil {
        return err
    }
    dish.SetAvailable(available)
    return s.repo.Update(dish)
}

func (s *DishService) UpdateDishPrice(id int, price float64) error {
    dish, err := s.repo.GetByID(id)
    if err != nil {
        return err
    }
    dish.UpdatePrice(price)
    return s.repo.Update(dish)
}

func (s *DishService) DeleteDish(id int) error {
    return s.repo.Delete(id)
}
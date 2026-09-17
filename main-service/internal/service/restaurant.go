package service

import (
    "avito-kitchen/main-service/internal/models"
    "avito-kitchen/main-service/internal/repository/postgres"
)

type RestaurantService struct {
    repo *postgres.RestaurantRepository
}

func NewRestaurantService(repo *postgres.RestaurantRepository) *RestaurantService {
    return &RestaurantService{repo: repo}
}

func (s *RestaurantService) CreateRestaurant(name, description, address, phone string) (*models.Restaurant, error) {
    restaurant := models.NewRestaurant(name, description, address, phone)
    if err := s.repo.Create(restaurant); err != nil {
        return nil, err
    }
    return restaurant, nil
}

func (s *RestaurantService) GetRestaurant(id int) (*models.Restaurant, error) {
    return s.repo.GetByID(id)
}

func (s *RestaurantService) GetAllRestaurants() ([]*models.Restaurant, error) {
    return s.repo.GetAll()
}

func (s *RestaurantService) ActivateRestaurant(id int) error {
    restaurant, err := s.repo.GetByID(id)
    if err != nil {
        return err
    }
    restaurant.Activate()
    return s.repo.Update(restaurant)
}

func (s *RestaurantService) DeactivateRestaurant(id int) error {
    restaurant, err := s.repo.GetByID(id)
    if err != nil {
        return err
    }
    restaurant.Deactivate()
    return s.repo.Update(restaurant)
}
package postgres

import (
    "gorm.io/gorm"

    "avito-kitchen/main-service/internal/errors"
    "avito-kitchen/main-service/internal/models"
)

type RestaurantRepository struct {
    db *gorm.DB
}

func NewRestaurantRepository(db *gorm.DB) *RestaurantRepository {
    return &RestaurantRepository{db: db}
}

func (r *RestaurantRepository) Create(restaurant *models.Restaurant) error {
    return r.db.Create(restaurant).Error
}

func (r *RestaurantRepository) GetByID(id int) (*models.Restaurant, error) {
    var restaurant models.Restaurant
    if err := r.db.First(&restaurant, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errors.ErrRestaurantNotFound
        }
        return nil, err
    }
    return &restaurant, nil
}

func (r *RestaurantRepository) GetAll() ([]*models.Restaurant, error) {
    var restaurants []*models.Restaurant
    if err := r.db.Find(&restaurants).Error; err != nil {
        return nil, err
    }
    return restaurants, nil
}

func (r *RestaurantRepository) Update(restaurant *models.Restaurant) error {
    return r.db.Save(restaurant).Error
}

func (r *RestaurantRepository) Delete(id int) error {
    return r.db.Delete(&models.Restaurant{}, id).Error
}
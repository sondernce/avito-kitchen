package postgres

import (
    "gorm.io/gorm"

    "avito-kitchen/main-service/internal/errors"
    "avito-kitchen/main-service/internal/models"
)

type DishRepository struct {
    db *gorm.DB
}

func NewDishRepository(db *gorm.DB) *DishRepository {
    return &DishRepository{db: db}
}

func (r *DishRepository) Create(dish *models.Dish) error {
    return r.db.Create(dish).Error
}

func (r *DishRepository) GetByID(id int) (*models.Dish, error) {
    var dish models.Dish
    if err := r.db.First(&dish, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errors.ErrDishNotFound
        }
        return nil, err
    }
    return &dish, nil
}

func (r *DishRepository) GetByRestaurantID(restaurantID int) ([]*models.Dish, error) {
    var dishes []*models.Dish
    if err := r.db.Where("restaurant_id = ?", restaurantID).Find(&dishes).Error; err != nil {
        return nil, err
    }
    return dishes, nil
}

func (r *DishRepository) Update(dish *models.Dish) error {
    return r.db.Save(dish).Error
}

func (r *DishRepository) Delete(id int) error {
    return r.db.Delete(&models.Dish{}, id).Error
}
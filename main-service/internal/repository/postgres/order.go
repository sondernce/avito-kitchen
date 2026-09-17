package postgres

import (
    "gorm.io/gorm"

    "avito-kitchen/main-service/internal/errors"
    "avito-kitchen/main-service/internal/models"
)

type OrderRepository struct {
    db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
    return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *models.Order) error {
    return r.db.Create(order).Error
}

func (r *OrderRepository) GetByID(id int) (*models.Order, error) {
    var order models.Order
    if err := r.db.Preload("Items").First(&order, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errors.ErrOrderNotFound
        }
        return nil, err
    }
    return &order, nil
}

func (r *OrderRepository) GetByUserID(userID int) ([]*models.Order, error) {
    var orders []*models.Order
    if err := r.db.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
        return nil, err
    }
    return orders, nil
}

func (r *OrderRepository) GetByRestaurantID(restaurantID int) ([]*models.Order, error) {
    var orders []*models.Order
    if err := r.db.Where("restaurant_id = ?", restaurantID).Find(&orders).Error; err != nil {
        return nil, err
    }
    return orders, nil
}

func (r *OrderRepository) Update(order *models.Order) error {
    return r.db.Save(order).Error
}

func (r *OrderRepository) Delete(id int) error {
    return r.db.Delete(&models.Order{}, id).Error
}
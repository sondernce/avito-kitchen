package handlers

import (
    "encoding/json"
    stderrors "errors"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"

    kitchenerrors "avito-kitchen/main-service/internal/errors"
    "avito-kitchen/main-service/internal/models"
    "avito-kitchen/main-service/internal/service"
)

type OrderHandlers struct {
    orderService *service.OrderService
}

func NewOrderHandlers(orderService *service.OrderService) *OrderHandlers {
    return &OrderHandlers{
        orderService: orderService,
    }
}

// POST /api/v1/orders
func (h *OrderHandlers) CreateOrder(w http.ResponseWriter, r *http.Request) {
    var dto CreateOrderDTO
    if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
        writeError(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := dto.Validate(); err != nil {
        writeError(w, err.Error(), http.StatusBadRequest)
        return
    }

    items := make([]models.OrderItem, len(dto.Items))
    for i, item := range dto.Items {
        items[i] = models.OrderItem{
            DishID:   item.DishID,
            Quantity: item.Quantity,
        }
    }

    order, err := h.orderService.CreateOrder(
        dto.UserID,
        dto.RestaurantID,
        items,
        dto.DeliveryAddress,
        dto.Phone,
    )
    if err != nil {
        writeError(w, err.Error(), http.StatusInternalServerError)
        return
    }

    writeJSON(w, order, http.StatusCreated)
}

// GET /api/v1/orders/{id}
func (h *OrderHandlers) GetOrder(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        writeError(w, "invalid order id", http.StatusBadRequest)
        return
    }

    order, err := h.orderService.GetOrder(id)
    if err != nil {
        if stderrors.Is(err, kitchenerrors.ErrOrderNotFound) {
            writeError(w, err.Error(), http.StatusNotFound)
        } else {
            writeError(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    writeJSON(w, order, http.StatusOK)
}

// GET /api/v1/users/{id}/orders
func (h *OrderHandlers) GetUserOrders(w http.ResponseWriter, r *http.Request) {
    userID, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        writeError(w, "invalid user id", http.StatusBadRequest)
        return
    }

    orders, err := h.orderService.GetUserOrders(userID)
    if err != nil {
        writeError(w, err.Error(), http.StatusInternalServerError)
        return
    }
    writeJSON(w, orders, http.StatusOK)
}

// GET /api/v1/restaurants/{id}/orders
func (h *OrderHandlers) GetRestaurantOrders(w http.ResponseWriter, r *http.Request) {
    restaurantID, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        writeError(w, "invalid restaurant id", http.StatusBadRequest)
        return
    }

    orders, err := h.orderService.GetRestaurantOrders(restaurantID)
    if err != nil {
        writeError(w, err.Error(), http.StatusInternalServerError)
        return
    }
    writeJSON(w, orders, http.StatusOK)
}

// PATCH /api/v1/orders/{id}/status
func (h *OrderHandlers) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        writeError(w, "invalid order id", http.StatusBadRequest)
        return
    }

    var dto UpdateStatusDTO
    if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
        writeError(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := dto.Validate(); err != nil {
        writeError(w, err.Error(), http.StatusBadRequest)
        return
    }

    status := models.OrderStatus(dto.Status)
    if err := h.orderService.UpdateOrderStatus(id, status); err != nil {
        if stderrors.Is(err, kitchenerrors.ErrOrderNotFound) {
            writeError(w, err.Error(), http.StatusNotFound)
        } else if stderrors.Is(err, kitchenerrors.ErrInvalidStatusTransition) {
            writeError(w, err.Error(), http.StatusConflict)
        } else {
            writeError(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    w.WriteHeader(http.StatusOK)
}

// PATCH /api/v1/orders/{id}/cancel
func (h *OrderHandlers) CancelOrder(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        writeError(w, "invalid order id", http.StatusBadRequest)
        return
    }

    if err := h.orderService.CancelOrder(id); err != nil {
        if stderrors.Is(err, kitchenerrors.ErrOrderNotFound) {
            writeError(w, err.Error(), http.StatusNotFound)
        } else {
            writeError(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    w.WriteHeader(http.StatusOK)
}
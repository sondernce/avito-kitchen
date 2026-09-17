package handlers

import (
    "encoding/json"
    stderrors "errors"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"

    kitchenerrors "avito-kitchen/main-service/internal/errors"
    "avito-kitchen/main-service/internal/service"
)

type DishHandlers struct {
    dishService *service.DishService
}

func NewDishHandlers(dishService *service.DishService) *DishHandlers {
    return &DishHandlers{
        dishService: dishService,
    }
}

// POST /api/v1/dishes
func (h *DishHandlers) CreateDish(w http.ResponseWriter, r *http.Request) {
    var dto CreateDishDTO
    if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
        writeError(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := dto.Validate(); err != nil {
        writeError(w, err.Error(), http.StatusBadRequest)
        return
    }

    dish, err := h.dishService.CreateDish(
        dto.RestaurantID,
        dto.Name,
        dto.Description,
        dto.Price,
        dto.Category,
    )
    if err != nil {
        writeError(w, err.Error(), http.StatusInternalServerError)
        return
    }

    writeJSON(w, dish, http.StatusCreated)
}

// GET /api/v1/restaurants/{id}/menu
func (h *DishHandlers) GetMenu(w http.ResponseWriter, r *http.Request) {
    restaurantID, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        writeError(w, "invalid restaurant id", http.StatusBadRequest)
        return
    }

    dishes, err := h.dishService.GetMenu(restaurantID)
    if err != nil {
        writeError(w, err.Error(), http.StatusInternalServerError)
        return
    }
    writeJSON(w, dishes, http.StatusOK)
}

// PATCH /api/v1/dishes/{id}/availability
func (h *DishHandlers) UpdateDishAvailability(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        writeError(w, "invalid dish id", http.StatusBadRequest)
        return
    }

    var dto struct {
        Available bool `json:"available"`
    }
    if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
        writeError(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.dishService.UpdateDishAvailability(id, dto.Available); err != nil {
        if stderrors.Is(err, kitchenerrors.ErrDishNotFound) {
            writeError(w, err.Error(), http.StatusNotFound)
        } else {
            writeError(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    w.WriteHeader(http.StatusOK)
}

// DELETE /api/v1/dishes/{id}
func (h *DishHandlers) DeleteDish(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        writeError(w, "invalid dish id", http.StatusBadRequest)
        return
    }

    if err := h.dishService.DeleteDish(id); err != nil {
        if stderrors.Is(err, kitchenerrors.ErrDishNotFound) {
            writeError(w, err.Error(), http.StatusNotFound)
        } else {
            writeError(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
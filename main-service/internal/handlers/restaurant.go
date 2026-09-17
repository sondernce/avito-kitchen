package handlers

import (
    "encoding/json"
    stderrors "errors"
    "net/http"
    "strconv"
    "time"

    "github.com/gorilla/mux"

    kitchenerrors "avito-kitchen/main-service/internal/errors"
    "avito-kitchen/main-service/internal/service"
)

type RestaurantHandlers struct {
    restaurantService *service.RestaurantService
}

func NewRestaurantHandlers(restaurantService *service.RestaurantService) *RestaurantHandlers {
    return &RestaurantHandlers{
        restaurantService: restaurantService,
    }
}

// POST /api/v1/restaurants
func (h *RestaurantHandlers) CreateRestaurant(w http.ResponseWriter, r *http.Request) {
    var dto CreateRestaurantDTO
    if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
        writeError(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := dto.Validate(); err != nil {
        writeError(w, err.Error(), http.StatusBadRequest)
        return
    }

    restaurant, err := h.restaurantService.CreateRestaurant(
        dto.Name,
        dto.Description,
        dto.Address,
        dto.Phone,
    )
    if err != nil {
        writeError(w, err.Error(), http.StatusInternalServerError)
        return
    }

    writeJSON(w, restaurant, http.StatusCreated)
}

// GET /api/v1/restaurants
func (h *RestaurantHandlers) GetAllRestaurants(w http.ResponseWriter, r *http.Request) {
    restaurants, err := h.restaurantService.GetAllRestaurants()
    if err != nil {
        writeError(w, err.Error(), http.StatusInternalServerError)
        return
    }
    writeJSON(w, restaurants, http.StatusOK)
}

// GET /api/v1/restaurants/{id}
func (h *RestaurantHandlers) GetRestaurant(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        writeError(w, "invalid id", http.StatusBadRequest)
        return
    }

    restaurant, err := h.restaurantService.GetRestaurant(id)
    if err != nil {
        if stderrors.Is(err, kitchenerrors.ErrRestaurantNotFound) {
            writeError(w, err.Error(), http.StatusNotFound)
        } else {
            writeError(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    writeJSON(w, restaurant, http.StatusOK)
}

// PATCH /api/v1/restaurants/{id}/activate
func (h *RestaurantHandlers) ActivateRestaurant(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        writeError(w, "invalid id", http.StatusBadRequest)
        return
    }

    if err := h.restaurantService.ActivateRestaurant(id); err != nil {
        if stderrors.Is(err, kitchenerrors.ErrRestaurantNotFound) {
            writeError(w, err.Error(), http.StatusNotFound)
        } else {
            writeError(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    w.WriteHeader(http.StatusOK)
}

// PATCH /api/v1/restaurants/{id}/deactivate
func (h *RestaurantHandlers) DeactivateRestaurant(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        writeError(w, "invalid id", http.StatusBadRequest)
        return
    }

    if err := h.restaurantService.DeactivateRestaurant(id); err != nil {
        if stderrors.Is(err, kitchenerrors.ErrRestaurantNotFound) {
            writeError(w, err.Error(), http.StatusNotFound)
        } else {
            writeError(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    w.WriteHeader(http.StatusOK)
}

func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
    b, err := json.MarshalIndent(data, "", "   ")
    if err != nil {
        writeError(w, "failed to marshal response", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    w.Write(b)
}

func writeError(w http.ResponseWriter, message string, statusCode int) {
    errDTO := ErrorDTO{
        Message: message,
        Time:    time.Now(),
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    w.Write([]byte(errDTO.ToString()))
}
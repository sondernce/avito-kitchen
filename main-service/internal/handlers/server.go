package handlers

import (
    "net/http"

    "github.com/gorilla/mux"
)

type Server struct {
    restaurantHandlers *RestaurantHandlers
    dishHandlers       *DishHandlers
    orderHandlers      *OrderHandlers
}

func NewServer(
    restaurantHandlers *RestaurantHandlers,
    dishHandlers *DishHandlers,
    orderHandlers *OrderHandlers,
) *Server {
    return &Server{
        restaurantHandlers: restaurantHandlers,
        dishHandlers:       dishHandlers,
        orderHandlers:      orderHandlers,
    }
}

func (s *Server) Start(addr string) error {
    router := mux.NewRouter()

    // Health check
    router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("{\"status\":\"ok\"}"))
    }).Methods("GET")

    api := router.PathPrefix("/api/v1").Subrouter()

    // Restaurant routes
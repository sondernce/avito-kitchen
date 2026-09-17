package client

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type MainServiceClient struct {
    baseURL    string
    httpClient *http.Client
}

func NewMainServiceClient(baseURL string) *MainServiceClient {
    return &MainServiceClient{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

// UpdateOrderStatus отправляет обновление статуса заказа в main-service
func (c *MainServiceClient) UpdateOrderStatus(orderID int, status string) error {
    url := fmt.Sprintf("%s/api/v1/orders/%d/status", c.baseURL, orderID)
    
    payload := map[string]string{
        "status": status,
    }
    
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("failed to marshal payload: %w", err)
    }

    req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    return nil
}

// GetMenu получает меню ресторана из main-service
func (c *MainServiceClient) GetMenu(restaurantID int) ([]map[string]interface{}, error) {
    url := fmt.Sprintf("%s/api/v1/restaurants/%d/menu", c.baseURL, restaurantID)
    
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to get menu: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    var menu []map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&menu); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    return menu, nil
}
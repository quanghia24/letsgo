package alihunter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/quanghia24/letsgo/internal/model"
)

const (
	ServiceURL  = "https://product-source-api.staging.alihunter.io/aliexpress/api/products/ds-image-search-v2"
	MaxProducts = 3
)

type AliHunterSearchByImageRequest struct {
	ImageURL   string `json:"image_url"`
	SearchType string `json:"search_type"`
	Currency   string `json:"currency"`
	Lang       string `json:"lang"`
	ShipTo     string `json:"ship_to"`
}

// AliHunterSearchByImage fetches product data from alihunter API
func AliHunterSearchByImage(url string) ([]model.AliHunterProduct, error) {
	// Validate input
	if url == "" {
		return nil, fmt.Errorf("image URL cannot be empty")
	}

	arg := AliHunterSearchByImageRequest{
		ImageURL:   url,
		SearchType: "same",
		Currency:   "USD",
		Lang:       "en",
		ShipTo:     "US",
	}

	body, err := json.Marshal(arg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, ServiceURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data model.AliHunterSearchByImageResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Filter and collect products
	var products []model.AliHunterProduct

	for _, item := range data.Result.Data.Data {
		// Skip products without ratings or image URL
		if item.EvaluateRate == "" || item.ProductMainImageURL == "" {
			continue
		}

		products = append(products, item)
		if len(products) >= MaxProducts {
			break
		}
	}

	return products, nil
}

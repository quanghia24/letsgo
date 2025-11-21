package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type FetchProductsRequest struct {
	ImageURL   string `json:"image_url"`
	SearchType string `json:"search_type"`
	Currency   string `json:"currency"`
	Lang       string `json:"lang"`
	ShipTo     string `json:"ship_to"`
}

// fetchProducts fetches product data from alihunter API
func fetchProducts(url string) (*AliHunterResponse, error) {
	// call external service to fetch products
	arg := FetchProductsRequest{
		ImageURL:   url,
		SearchType: "same",
		Currency:   "USD",
		Lang:       "en",
		ShipTo:     "US",
	}

	serviceURL := "https://product-source-api.staging.alihunter.io/aliexpress/api/products/ds-image-search-v2"

	body, err := json.Marshal(arg)
	if err != nil {
		log.Fatal("failed to marshal request body")
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, serviceURL, bytes.NewReader(body))
	if err != nil {
		log.Fatal("failed to create new request")
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("failed to perform request")
		return nil, err
	}
	defer resp.Body.Close()

	var data AliHunterResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatal("failed to decode response body")
		return nil, err
	}

	// empty data case
	// data = AliHunterResponse{}

	return &data, nil
}

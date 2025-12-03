package rapidapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/quanghia24/letsgo/configs"
	"github.com/quanghia24/letsgo/internal/model"
)

// AliExpressSearchByImage fetches products from AliExpress API with endpoint get from .env
// Return top 3 products
func AliExpressSearchByImage(image string) ([]model.AliExpressProduct, []model.AliExpressProduct, error) {
	if image == "" {
		return nil, nil, fmt.Errorf("image URL is empty")
	}

	// URL encode the image parameter to handle special characters
	encodedImage := url.QueryEscape(image)
	serviceURL := fmt.Sprintf("https://%s/item_search_image?sort=default&catId=0&imgUrl=%s", configs.GetRapidAPIConfig().Host, encodedImage)

	req, err := http.NewRequest(http.MethodGet, serviceURL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-RapidAPI-Key", configs.GetRapidAPIConfig().APIKey)
	req.Header.Set("X-RapidAPI-Host", configs.GetRapidAPIConfig().Host)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, resp.Status)
	}

	var data model.AliExpressSearchByImageResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	// Debug: log the response for troubleshooting
	log.Printf("AliExpress API Response - URL: %s, Results count: %d", serviceURL, len(data.Result.ResultList))

	var products []model.AliExpressProduct
	var originProducts []model.AliExpressProduct

	for _, result := range data.Result.ResultList {
		item := result.Item

		// Skip items with missing critical data early
		if item.Image == "" {
			continue
		}

		// Safely extract price
		price, ok := item.Sku.Def.Price.(float64)
		if !ok {
			price = item.Sku.Def.PromotionPrice
		}

		// Safely extract rating
		avgRating := 0.0
		if rating, ok := item.AverageStarRate.(float64); ok {
			avgRating = rating
		}

		product := model.AliExpressProduct{
			ProductID:     item.ItemID,
			URL:           item.ItemURL,
			Title:         item.Title,
			ImageURL:      item.Image,
			AvgRatingStar: avgRating,
			Volume:        item.Sales,
			SalePrice:     item.Sku.Def.PromotionPrice,
			OriginalPrice: price,
		}

		originProducts = append(originProducts, product)

		// Only add to filtered products if it has ratings
		if item.AverageStarRate != nil {
			products = append(products, product)
		}
	}

	// get max 3 products
	if len(products) > 3 {
		products = products[:3]
	}
	if len(originProducts) > 3 {
		originProducts = originProducts[:3]
	}

	return products, originProducts, nil
}

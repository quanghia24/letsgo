package rapidapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/quanghia24/letsgo/internal/config"
	"github.com/quanghia24/letsgo/internal/model"
)

// AliExpressSearchByImage fetches products from AliExpress API with endpoint get from .env
// Return top 3 products
func AliExpressSearchByImage(image string) ([]model.AliExpressProduct, []model.AliExpressProduct, error) {
	if image == "" {
		return nil, nil, fmt.Errorf("image URL is empty")
	}

	serviceURL := fmt.Sprintf("https://%s/item_search_image?sort=default&catId=0&imgUrl=%s", config.GetRapidAPIConfig().Host, image)
	req, err := http.NewRequest(http.MethodGet, serviceURL, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-RapidAPI-Key", config.GetRapidAPIConfig().APIKey)
	req.Header.Set("X-RapidAPI-Host", config.GetRapidAPIConfig().Host)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("failed to perform request")
		return nil, nil, err
	}
	defer resp.Body.Close()

	var data model.AliExpressSearchByImageResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatal("failed to decode response body")
		return nil, nil, err
	}

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
			ID:            item.ItemID,
			URL:           item.ItemURL,
			Title:         item.Title,
			ImageURL:      item.Image,
			AvgRatingStar: avgRating,
			Volume:        item.Sales,
			SalePrice:     item.Sku.Def.PromotionPrice,
			OriginalPrice: price,
		}

		originProducts = append(originProducts, product)

		if item.AverageStarRate == nil {
			continue
		}

		products = append(products, product)
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

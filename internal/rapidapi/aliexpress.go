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
func AliExpressSearchByImage(image string) ([]model.AliExpressProduct, error) {
	if image == "" {
		return nil, fmt.Errorf("image URL is empty")
	}

	serviceURL := fmt.Sprintf("https://%s/item_search_image?sort=default&catId=0&imgUrl=%s", config.GetRapidAPIConfig().Host, image)
	req, err := http.NewRequest(http.MethodGet, serviceURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-RapidAPI-Key", config.GetRapidAPIConfig().APIKey)
	req.Header.Set("X-RapidAPI-Host", config.GetRapidAPIConfig().Host)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("failed to perform request")
		return nil, err
	}
	defer resp.Body.Close()

	var data model.AliExpressSearchByImageResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatal("failed to decode response body")
		return nil, err
	}

	var products []model.AliExpressProduct
	for _, result := range data.Result.ResultList {
		item := result.Item

		// filter out items with no rating or no sales or no image
		if item.AverageStarRate == nil || item.Sales == 0 || item.Image == "" {
			continue
		}

		price, ok := item.Sku.Def.Price.(float64)
		if !ok {
			price = item.Sku.Def.PromotionPrice
		}

		products = append(products, model.AliExpressProduct{
			ID:            item.ItemID,
			URL:           item.ItemURL,
			Title:         item.Title,
			ImageURL:      item.Image,
			AvgRatingStar: item.AverageStarRate.(float64),
			Volume:        item.Sales,
			SalePrice:     item.Sku.Def.PromotionPrice,
			OriginalPrice: price,
		})
	}

	// get max 3 products
	if len(products) > 3 {
		products = products[:3]
	}

	return products, nil
}

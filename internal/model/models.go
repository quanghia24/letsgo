package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AliExpressSearchByImageResponse struct {
	Result struct {
		ResultList []*ResultListSearchByImage `json:"resultList"`
	} `json:"result"`
}

type ResultListSearchByImage struct {
	Item struct {
		ItemID  string `json:"itemId"`
		Title   string `json:"title"`
		Sales   int64  `json:"sales"`
		ItemURL string `json:"itemUrl"`
		Image   string `json:"image"`
		Sku     struct {
			Def struct {
				Price          interface{} `json:"price"`
				PromotionPrice float64     `json:"promotionPrice"`
			} `json:"def"`
		} `json:"sku"`
		AverageStarRate interface{} `json:"averageStarRate"`
	}
}

type AliExpressProduct struct {
	ID            string  // Product ID
	URL           string  // Product URL
	Title         string  // Product title
	ImageURL      string  // Main product image URL
	AvgRatingStar float64 // Average rating (0-5)
	Volume        int64   // Sales volume
	SalePrice     float64 // Current sale price
	OriginalPrice float64 // Original price
}

type AliHunterSearchByImageResponse struct {
	Result struct {
		Ret  bool `json:"ret"`
		Data struct {
			Data []AliHunterProduct `json:"data"`
		} `json:"data"`
	} `json:"result"`
}

type AliHunterProduct struct {
	ProductID               string `json:"product_id"`
	EvaluateRate            string `json:"evaluate_rate"`
	ProductTitle            string `json:"product_title"`
	ProductMainImageURL     string `json:"product_main_image_url"`
	ProductDetailURL        string `json:"product_detail_url"`
	TargetSalePrice         string `json:"target_sale_price"`
	TargetOriginalPrice     string `json:"target_original_price"`
	LatestVolume            string `json:"latest_volume"`
	SimilarityScore         string `json:"similarity_score"`
	ShipFrom                string `json:"ship_from"`
	TargetSalePriceCurrency string `json:"target_sale_price_currency"`
	Matching                bool   `json:"matching"`
}

type RapidapiResponse struct {
	ID        string            `json:"_id"`
	CreatedAt string            `json:"created_at"`
	ImageURL  string            `json:"image_url"`
	JobID     string            `json:"job_id"`
	ProductID string            `json:"product_id"`
	Products  []RapidapiProduct `json:"products"`
	ShopID    string            `json:"shop_id"`
	Status    string            `json:"status"`
	UpdatedAt string            `json:"updated_at"`
}

type RapidapiProduct struct {
	Avgstar             float32 `json:"avgstar"`
	Platform            string  `json:"platform"`
	ProductID           string  `json:"productid"`
	ProductMainImageURL string  `json:"productmainimageurl"`
	ProductTitle        string  `json:"producttitle"`
	ProductURL          string  `json:"producturl"`
	Sale                int     `json:"sale"`
	TargetOriginalPrice string  `json:"targetoriginalprice"`
	TargetSalePrice     string  `json:"targetsaleprice"`
	TotalReview         string  `json:"totalreview"`
	Type                string  `json:"type"`
	Matching            bool    `json:"matching"`
}

type ShopGroup struct {
	ShopID             int64               `json:"shop_id"`
	Shop               Shop                `json:"shop"`
	ProductCount       int                 `json:"product_count"`
	SuggestionProducts []SuggestionProduct `json:"suggestion_products"`
}

type Shop struct {
	ShopID          int64  `bson:"shop_id" json:"shop_id"`
	MyshopifyDomain string `bson:"myshopify_domain" json:"myshopify_domain"`
	PlanDisplayName string `bson:"plan_display_name" json:"plan_display_name"`
	AppPlan         string `bson:"app_plan" json:"app_plan"`
	Domain          string `bson:"domain" json:"domain"`
}

type SuggestionProduct struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	ShopID     int64              `bson:"shop_id" json:"shop_id"`
	ProductID  int64              `bson:"product_id" json:"product_id"`
	JobID      string             `bson:"job_id" json:"job_id"`
	ImageURL   string             `bson:"image_url" json:"image_url"`
	Status     string             `bson:"status" json:"status"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
	Products   []ProductItem      `bson:"products" json:"products"`
	Product    Product            `bson:"product" json:"product"`
	ProductURL string             `bson:"producturl" json:"producturl"`
}

type ProductItem struct {
	Type                string      `bson:"type" json:"type"`
	Platform            string      `bson:"platform" json:"platform"`
	ProductID           string      `bson:"productid" json:"productid"`
	ProductURL          string      `bson:"producturl" json:"producturl"`
	ProductMainImageURL string      `bson:"productmainimageurl" json:"productmainimageurl"`
	ProductTitle        string      `bson:"producttitle" json:"producttitle"`
	TargetSalePrice     string      `bson:"targetsaleprice" json:"targetsaleprice"`
	TargetOriginalPrice string      `bson:"targetoriginalprice" json:"targetoriginalprice"`
	AvgStar             float64     `bson:"avgstar" json:"avgstar"`
	Sale                int         `bson:"sale" json:"sale"`
	TotalReview         interface{} `bson:"totalreview" json:"totalreview"` // Can be int or object like {"$numberLong": "..."}
}

type Product struct {
	ProductID    int64  `bson:"product_id" json:"product_id"`
	ShopID       int64  `bson:"shop_id" json:"shop_id"`
	Title        string `bson:"title" json:"title"`
	Image        string `bson:"image" json:"image"`
	Type         string `bson:"type" json:"type"`
	Handle       string `bson:"handle" json:"handle"`
	ProductURL   string `bson:"product_url" json:"product_url"`
	TotalReviews int64  `bson:"total_reviews" json:"total_reviews"`
}

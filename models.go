package main

type AliHunterResponse struct {
	Result struct {
		Ret  bool `json:"ret"`
		Data struct {
			Data []AliHunterProduct `json:"data"`
		} `json:"data"`
	} `json:"result"`
}

type AliHunterProduct struct {
	ProductID               string `json:"product_id"`
	ProductTitle            string `json:"product_title"`
	ProductMainImageURL     string `json:"product_main_image_url"`
	ProductDetailURL        string `json:"product_detail_url"`
	TargetSalePrice         string `json:"target_sale_price"`
	TargetOriginalPrice     string `json:"target_original_price"`
	LatestVolume            string `json:"latest_volume"`
	SimilarityScore         string `json:"similarity_score"`
	ShipFrom                string `json:"ship_from"`
	TargetSalePriceCurrency string `json:"target_sale_price_currency"`
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
}

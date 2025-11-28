package report

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/quanghia24/letsgo/internal/model"
)

type ListReports struct {
	GeneratedAt     string
	Comparisons     []Report
	ComparisonsJSON string
}

// Report is the view-model passed to the HTML template
type Report struct {
	ProductTitle        string
	ProductID           int64
	ImageURL            string
	ShopID              int64
	LocalRapidAPITop    []model.ProductItem
	LocalRapidAPIOrigin []model.ProductItem
	AliHunterTop        []model.AliHunterProduct
	AliHunterOrigin     []model.AliHunterProduct
	AliExpressTop       []model.AliExpressProduct
	AliExpressOrigin    []model.AliExpressProduct
}

func TakeTopProducts(input []model.ProductItem) ([]model.ProductItem, []model.ProductItem) {
	var filtered []model.ProductItem
	count := 0

	for _, item := range input {
		if count < 3 && item.ProductMainImageURL != "" {
			filtered = append(filtered, item)
			count++
		}
	}

	// Return original slice up to its length or 3, whichever is smaller
	maxLen := len(input)
	if maxLen > 3 {
		maxLen = 3
	}
	return filtered, input[:maxLen]
}

// GenerateHTMLReport writes a single HTML file containing all provided report
func GenerateHTMLReport(reports []Report, outDir string) error {
	// Marshal comparisons to JSON for JavaScript
	comparisonsJSON, err := json.Marshal(reports)
	if err != nil {
		return fmt.Errorf("failed to marshal comparisons to JSON: %w", err)
	}

	listReports := ListReports{
		GeneratedAt:     time.Now().Format(time.RFC3339),
		Comparisons:     reports,
		ComparisonsJSON: string(comparisonsJSON),
	}

	tmplPath := "./templates/report.tmpl"
	// register template functions
	funcMap := template.FuncMap{
		"formatPrice": formatPrice,
	}
	t, err := template.New("report.tmpl").Funcs(funcMap).ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", tmplPath, err)
	}

	f, err := os.Create(outDir)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outDir, err)
	}
	defer f.Close()

	if err := t.Execute(f, listReports); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// format price from cents to dollars
func formatPrice(price string) string {
	for len(price) < 3 {
		price = "0" + price
	}
	price = price[:len(price)-2] + "." + price[len(price)-2:]
	return "$" + price
}

func GenerateJSONComparisonReport(reports []Report) error {
	data, err := json.MarshalIndent(reports, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal reports to json: %w", err)
	}
	if err := os.WriteFile("report.json", data, 0644); err != nil {
		return fmt.Errorf("failed to write json file: %w", err)
	}
	return nil
}

type Review struct {
	ProductID    string `json:"product_id"`
	ReviewsCount string `json:"reviews_count"`
}

type ListReviews struct {
	AliHunterTop     []Review
	AliHunterOrigin  []Review
	AliExpressTop    []Review
	AliExpressOrigin []Review
}

func UpdateJSONComparisonReview() error {
	// get list of reports from report.json
	var reports []Report
	data, err := os.ReadFile("report.json")
	if err != nil {
		return fmt.Errorf("failed to read report.json file: %w", err)
	}
	if err := json.Unmarshal(data, &reports); err != nil {
		return fmt.Errorf("failed to unmarshal report.json file: %w", err)
	}

	// update report.json, on each report with reviews count
	for i := range reports { // Use index to modify the actual slice element
		for j := range reports[i].AliHunterTop {
			productID := reports[i].AliHunterTop[j].ProductID
			count, err := getReviewsCount(productID)
			if err != nil {
				fmt.Printf("failed to get reviews count for AliHunter product %s: %v\n", productID, err)
			}
			reports[i].AliHunterTop[j].TotalReview = count
		}
		for j := range reports[i].AliHunterOrigin {
			productID := reports[i].AliHunterOrigin[j].ProductID
			count, err := getReviewsCount(productID)
			if err != nil {
				fmt.Printf("failed to get reviews count for AliHunter product %s: %v\n", productID, err)
			}
			reports[i].AliHunterOrigin[j].TotalReview = count
		}
		for j := range reports[i].AliExpressTop {
			productID := reports[i].AliExpressTop[j].ID
			count, err := getReviewsCount(productID)
			if err != nil {
				fmt.Printf("failed to get reviews count for AliExpress product %s: %v\n", productID, err)
			}
			reports[i].AliExpressTop[j].TotalReview = count
		}
		for j := range reports[i].AliExpressOrigin {
			productID := reports[i].AliExpressOrigin[j].ID
			count, err := getReviewsCount(productID)
			if err != nil {
				fmt.Printf("failed to get reviews count for AliExpress product %s: %v\n", productID, err)
			}
			reports[i].AliExpressOrigin[j].TotalReview = count
		}
	}

	// write back to report.json
	updatedData, err := json.MarshalIndent(reports, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated reports to json: %w", err)
	}
	if err := os.WriteFile("report.json", updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write updated json file: %w", err)
	}

	return nil
}

type getReviewsCountResponse struct {
	Data struct {
		TotalNum int `json:"totalNum"`
	} `json:"data"`
}

func getReviewsCount(productID string) (string, error) {
	serviceURL := fmt.Sprintf("https://feedback.aliexpress.com/pc/searchEvaluation.do?productId=%s&page=1", productID)

	req, err := http.NewRequest(http.MethodGet, serviceURL, nil)
	if err != nil {
		return "0 ratings", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "0 ratings", fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return "0 ratings", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data getReviewsCountResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "0 ratings", fmt.Errorf("failed to decode response: %w", err)
	}

	return fmt.Sprintf("%d ratings", data.Data.TotalNum), nil
}

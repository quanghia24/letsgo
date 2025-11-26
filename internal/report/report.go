package report

import (
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/quanghia24/letsgo/internal/model"
)

type ListReports struct {
	GeneratedAt string
	Comparisons []Report
}

// Report is the view-model passed to the HTML template
type Report struct {
	ProductID        string
	ImageURL         string
	LocalRapidAPITop []model.RapidapiProduct
	AliHunterTop     []model.AliHunterProduct
	AliExpressTop    []model.AliExpressProduct
}

func TakeTopProducts(input []model.RapidapiProduct) []model.RapidapiProduct {
	var res []model.RapidapiProduct
	len := 0

	for _, item := range input {
		if len < 3 && item.ProductMainImageURL != "" {
			res = append(res, item)
			len++
		} else {
			break
		}
	}
	return res
}

// GenerateHTMLReport writes a single HTML file containing all provided report
func GenerateHTMLReport(reports []Report, outDir string) error {
	listReports := ListReports{
		GeneratedAt: time.Now().Format(time.RFC3339),
		Comparisons: reports,
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

package main

import (
	"fmt"
	"html/template"
	"os"
	"time"
)

type ListReports struct {
	GeneratedAt string
	Comparisons []Report
}

// Report is the view-model passed to the HTML template
type Report struct {
	ProductID string
	ImageURL  string
	RapidTop  []RapidapiProduct
	AliTop    []AliHunterProduct
	AliError  string
}

// take top N products with non-empty image URLs
func takeTopRapid(input []RapidapiProduct, n int) []RapidapiProduct {
	var res []RapidapiProduct
	len := 0

	for _, item := range input {
		if len < n {
			if item.ProductMainImageURL != "" {
				res = append(res, item)
			}
			len++
		} else {
			break
		}
	}
	return res
}

// take top N products with non-empty image URLs
func takeTopAli(input []AliHunterProduct, n int) []AliHunterProduct {
	var res []AliHunterProduct
	len := 0

	for _, item := range input {
		if len < n {
			if item.ProductMainImageURL != "" {
				res = append(res, item)
			}
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

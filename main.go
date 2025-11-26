package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/quanghia24/letsgo/internal/alihunter"
	"github.com/quanghia24/letsgo/internal/model"
	"github.com/quanghia24/letsgo/internal/rapidapi"
	"github.com/quanghia24/letsgo/internal/report"
)

func main() {
	// 1. read JSON file path from flag
	filePath := flag.String("local", "./docs/suggest_products.json", "path to local JSON file with RapidAPI product suggestions")
	outDir := flag.String("output", "./report.html", "output directory for the generated HTML report")
	flag.Parse()

	fmt.Println("1️⃣ Reading: ", *filePath)
	fileBytes, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatal("cannot read local json file:", err)
	}
	var rapidapiResponses []model.RapidapiResponse
	if err := json.Unmarshal(fileBytes, &rapidapiResponses); err != nil {
		log.Fatal("cannot unmarshal json data:", err)
	}
	fmt.Println("⭐ Done reading JSON file")

	// 2. request product data from alihunter API and aliexpress then collect comparisons
	fmt.Println("2️⃣ Fetching data:", *filePath)

	type result struct {
		index      int
		comparison report.Report
	}

	total := len(rapidapiResponses)
	results := make([]result, total)
	resultsChan := make(chan result, total)
	sem := make(chan struct{}, 5) // semaphore limit to 5 concurrent requests, not to overwhelm the APIs
	var wg sync.WaitGroup

	// Process each product concurrently
	for i, item := range rapidapiResponses {
		wg.Add(1)
		go func(idx int, product model.RapidapiResponse) {
			defer wg.Done()
			sem <- struct{}{}        // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			var aliHunterProducts []model.AliHunterProduct
			var aliexpressProducts []model.AliExpressProduct
			var innerWg sync.WaitGroup

			innerWg.Add(2)

			// Fetch AliHunter
			go func() {
				defer innerWg.Done()
				products, err := alihunter.AliHunterSearchByImage(product.ImageURL)
				if err != nil {
					log.Printf("fliHunter failed for %s: %v\n", product.ProductID, err)
					aliHunterProducts = []model.AliHunterProduct{}
				} else {
					aliHunterProducts = products
				}
			}()

			// Fetch AliExpressparallel
			go func() {
				defer innerWg.Done()
				products, err := rapidapi.AliExpressSearchByImage(product.ImageURL)
				if err != nil {
					log.Printf("aliExpress failed for %s: %v\n", product.ProductID, err)
					aliexpressProducts = []model.AliExpressProduct{}
				} else {
					aliexpressProducts = products
				}
			}()

			innerWg.Wait() // Wait for both API calls to complete

			// Take top 3 local rapidapi products
			localRapidAPIProducts := report.TakeTopProducts(product.Products)

			// Send result to channel
			resultsChan <- result{
				index: idx,
				comparison: report.Report{
					ProductID:        product.ProductID,
					ImageURL:         product.ImageURL,
					LocalRapidAPITop: localRapidAPIProducts,
					AliHunterTop:     aliHunterProducts,
					AliExpressTop:    aliexpressProducts,
				},
			}
		}(i, item)
	}

	// Close channel when all goroutines complete
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	for res := range resultsChan {
		results[res.index] = res
	}

	// Build ordered comparisons list
	comparisons := make([]report.Report, total)
	for i, res := range results {
		comparisons[i] = res.comparison
	}

	fmt.Println("⭐ Finished fetching from alihunter API and preparing comparisons")

	// 3. Generates an interactive HTML comparison report
	fmt.Println("3️⃣ Generating HTML report")
	if err := report.GenerateHTMLReport(comparisons, *outDir); err != nil {
		log.Fatalf("failed to generate report: %v", err)
	}
	fmt.Println("⭐ Report successfully generated and saved to:", *outDir)
}

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
	// Parse command-line flags
	filePath := flag.String("local", "./docs/suggest_products.json", "path to local JSON file with RapidAPI product suggestions")
	htmlFlag := flag.Bool("html", false, "generate HTML report")
	flag.Parse()

	// Generates an interactive HTML comparison report: only run on htmlFlag set to true
	if *htmlFlag {
		fmt.Println("🌶️ Generating HTML report")
		var comparisons []report.Report
		fileBytes, err := os.ReadFile("report.json")
		if err != nil {
			log.Fatalf("failed to read report.json: %v", err)
		}
		if err := json.Unmarshal(fileBytes, &comparisons); err != nil {
			log.Fatalf("failed to unmarshal report.json: %v", err)
		}

		if err := report.GenerateHTMLReport(comparisons, "report.html"); err != nil {
			log.Fatalf("failed to generate report: %v", err)
		}

		fmt.Println("⭐ Report successfully generated and saved to report.html")
		return
	}

	// Gererate comparison report from local JSON file
	fmt.Println("1️⃣ Reading: ", *filePath)
	fileBytes, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatal("cannot read local json file:", err)
	}
	var ShopGroupResponses []model.ShopGroup
	if err := json.Unmarshal(fileBytes, &ShopGroupResponses); err != nil {
		log.Fatal("cannot unmarshal json data:", err)
	}
	fmt.Println("⭐ Done reading JSON file")

	// 2. request product data from alihunter API and aliexpress then collect comparisons
	fmt.Println("2️⃣ Fetching data:", *filePath)

	type result struct {
		index      int
		comparison report.Report
	}

	// Calculate total products across all shops
	total := 0
	for _, shop := range ShopGroupResponses {
		total += len(shop.SuggestionProducts)
	}

	results := make([]result, total)
	resultsChan := make(chan result, total)
	sem := make(chan struct{}, 7) // semaphore limit to 7 concurrent requests, not to overwhelm the APIs, why 7? cuz I like it =))
	var wg sync.WaitGroup

	// Use a global index counter
	globalIdx := 0
	for _, shop := range ShopGroupResponses {
		for _, product := range shop.SuggestionProducts {
			wg.Add(1)
			currIndex := globalIdx // index for storing result in order
			globalIdx++

			go func(idx int, prod model.SuggestionProduct) {
				defer wg.Done()
				sem <- struct{}{}        // Acquire semaphore
				defer func() { <-sem }() // Release semaphore

				var aliHunterProducts []model.AliHunterProduct
				var aliHunterOrigin []model.AliHunterProduct
				var aliexpressProducts []model.AliExpressProduct
				var aliexpressOrigin []model.AliExpressProduct
				var innerWg sync.WaitGroup

				innerWg.Add(2)

				// Fetch AliHunter
				go func() {
					defer innerWg.Done()
					products, originals, err := alihunter.AliHunterSearchByImage(prod.ImageURL)
					if err != nil {
						log.Printf("aliHunter failed for %d: %v\n", prod.ProductID, err)
						aliHunterProducts = []model.AliHunterProduct{}
						aliHunterOrigin = []model.AliHunterProduct{}
					} else { // query total reviews for each product
						for i := range products {
							count, err := report.GetReviewsCount(products[i].ProductID)
							if err != nil {
								log.Printf("failed to get reviews count for aliHunter product %s: %v\n", products[i].ProductID, err)
							}
							products[i].TotalReview = count
						}
						for i := range originals {
							count, err := report.GetReviewsCount(originals[i].ProductID)
							if err != nil {
								log.Printf("failed to get reviews count for aliHunter product %s: %v\n", originals[i].ProductID, err)
							}
							originals[i].TotalReview = count
						}

						aliHunterProducts = products
						aliHunterOrigin = originals
					}

				}()

				// Fetch AliExpress
				go func() {
					defer innerWg.Done()
					products, originals, err := rapidapi.AliExpressSearchByImage(prod.ImageURL)
					if err != nil {
						log.Printf("aliExpress failed for %d: %v\n", prod.ProductID, err)
						aliexpressProducts = []model.AliExpressProduct{}
						aliexpressOrigin = []model.AliExpressProduct{}
					} else {
						for i := range products {
							count, err := report.GetReviewsCount(products[i].ProductID)
							if err != nil {
								log.Printf("failed to get reviews count for aliExpress product %s: %v\n", products[i].ProductID, err)
							}
							products[i].TotalReview = count
						}
						for i := range originals {
							count, err := report.GetReviewsCount(originals[i].ProductID)
							if err != nil {
								log.Printf("failed to get reviews count for AliExpress product %s: %v\n", originals[i].ProductID, err)
							}
							originals[i].TotalReview = count
						}

						aliexpressProducts = products
						aliexpressOrigin = originals
					}
				}()

				innerWg.Wait() // Wait for both API calls to complete

				// Take top 3 local products
				localProducts, localOrigin := report.TakeTopProducts(prod.Products)

				// Send result to channel
				resultsChan <- result{
					index: idx,
					comparison: report.Report{
						ProductTitle:        prod.Product.Title,
						ProductID:           prod.ProductID,
						ImageURL:            prod.ImageURL,
						ShopID:              prod.ShopID,
						LocalRapidAPITop:    localProducts,
						LocalRapidAPIOrigin: localOrigin,
						AliHunterTop:        aliHunterProducts,
						AliHunterOrigin:     aliHunterOrigin,
						AliExpressTop:       aliexpressProducts,
						AliExpressOrigin:    aliexpressOrigin,
					},
				}

			}(currIndex, product)
		}
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

	// Build ordered comparisons list -> default behaviour: export to json
	comparisons := make([]report.Report, total)
	for i, res := range results {
		comparisons[i] = res.comparison
	}

	if err := report.GenerateJSONComparisonReport(comparisons); err != nil {
		log.Fatalf("failed to generate JSON report: %v", err)
	}

	fmt.Println("⭐ Finished fetching from alihunter API and preparing comparisons")
}

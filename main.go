package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// 1. read JSON file path from flag
	filePath := flag.String("local", "./docs/suggest_products.json", "path to local JSON file with RapidAPI product suggestions")
	outDir := flag.String("output", "./report.html", "output directory for the generated HTML report")
	flag.Parse()

	fmt.Println("⏰ Reading: ", *filePath)
	fileBytes, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatal("💥 cannot read local json file:", err)
	}
	var rapidapiResponses []RapidapiResponse
	if err := json.Unmarshal(fileBytes, &rapidapiResponses); err != nil {
		log.Fatal("💥 cannot unmarshal json data:", err)
	}
	fmt.Println("⭐ Done reading JSON file")

	// 2. request product data from alihunter API and collect comparisons
	fmt.Println("⏰ Fetching AliHunter API:", *filePath)
	var comparisons []Report
	for _, item := range rapidapiResponses {
		// call alihunter API to fetch products
		alihunterResponse, err := fetchProducts(item.ImageURL)
		var aliTop []AliHunterProduct
		var aliErr string
		if err != nil {
			log.Printf("💥 fetching err for product %s: %v", item.ProductID, err)
			aliErr = err.Error()
			aliTop = []AliHunterProduct{}
		} else {
			aliTop = takeTopAli(alihunterResponse.Result.Data.Data, 3)
		}

		// make comparison report
		report := Report{
			ProductID: item.ProductID,
			ImageURL:  item.ImageURL,
			RapidTop:  takeTopRapid(item.Products, 3),
			AliTop:    aliTop,
			AliError:  aliErr,
		}
		comparisons = append(comparisons, report)
	}
	fmt.Println("⭐ Finished fetching from alihunter API and preparing comparisons")

	// 3. Generates an interactive HTML comparison report
	fmt.Println("⏰ Generating HTML report")
	if err := GenerateHTMLReport(comparisons, *outDir); err != nil {
		log.Fatalf("💥 failed to generate report: %v", err)
	}
	fmt.Println("⭐ Report successfully generated and safed to:", *outDir)
}

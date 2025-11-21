# Comparison Tool - AliHunter vs RapidAPI

A Go command-line tool that compares product search results from local RapidAPI data with real-time AliHunter API responses.

**Usage:**

```Bash
go run . -local <products.json> -output <report.html>
```

---

## 📁 Code Structure

### Project Structure

```Bash
gofun/
├── main.go              # Entry point and orchestration
├── models.go            # Data structures
├── fetcher.go           # API integration
├── report.go            # Report generation logic
├── go.mod
├── products.json        # Input data (RapidAPI results)
└── templates/
    └── report.tmpl      # HTML template with styling and JS
```

## 🔄 Process Flow

```Text
1. INITIALIZATION
   ├─ Parse command-line flags (-local, -output)
   └─ Read and unmarshal input JSON file
          ↓
2. DATA PROCESSING (for each product)
   ├─ Extract image URL from RapidAPI data
   ├─ Call AliHunter API with image URL
   ├─ Handle API errors (log and continue)
   ├─ Filter top 3 products with valid images
   └─ Build Report struct with comparison data
          ↓
3. REPORT GENERATION
   ├─ Prepare ListReports view model
   ├─ Parse HTML template with custom functions
   ├─ Execute template with all comparison data
   └─ Write HTML file to specified output path
          ↓
4. COMPLETION
   └─ Log success message with output file path
```

**Console Output Example:**

```Bash
⏰ Reading: products.json
⭐ Done reading JSON file
⏰ Fetching AliHunter API: products.json
⭐ Finished fetching from alihunter API and preparing comparisons
⏰ Generating HTML report
⭐ Report successfully generated and saved to: report.html
```

---

## ✅ Requirements Compliance Review

### Error Handling Requirements

| Requirement | Status | Implementation Details |
|------------|--------|------------------------|
| **Processing multiple products** | ✅ **IMPLEMENTED** | `main.go` loops through all products in input JSON array |
| **Skipping products without image URLs** | ✅ **IMPLEMENTED** | `takeTopRapid()` and `takeTopAli()` filter out empty image URLs |
| **API errors (log and continue)** | ✅ **IMPLEMENTED** | API errors logged with `log.Printf()`, stored in `Report.AliError`, processing continues |
| **Empty API results** | ✅ **IMPLEMENTED** | Empty arrays handled gracefully, displays "No results" in HTML |
| **Invalid JSON input** | ✅ **IMPLEMENTED** | `json.Unmarshal()` returns error, triggers `log.Fatal()` |
| **Missing files** | ✅ **IMPLEMENTED** | `os.ReadFile()` returns error, triggers `log.Fatal()` |

### Acceptance Criteria Status

| Criterion | Status | Evidence |
|-----------|--------|----------|
| 1. Tool processes all products from input file | ✅ **COMPLETE** | Loop in `main.go` processes entire `rapidapiResponses` array |
| 2. API integration works correctly | ✅ **COMPLETE** | `fetchProducts()` makes POST request with proper headers and JSON body |
| 3. HTML report displays in 3-column layout | ✅ **COMPLETE** | Template uses flexbox: `w-1/5` (product), `w-2/5` (RapidAPI), `w-2/5` (AliHunter) |
| 4. All product information displays correctly | ✅ **COMPLETE** | Images, titles, prices, ratings, links all rendered from template data |
| 5. Checkboxes toggle card highlighting | ✅ **COMPLETE** | JavaScript adds/removes `.matched` class with green border |
| 6. Summary section shows/hides dynamically | ✅ **COMPLETE** | Fixed panel with `.active` class toggles `transform: translateY()` |
| 7. Statistics calculate correctly | ✅ **COMPLETE** | JS calculates percentages: `(checkedCount/totalCheckboxes)*100` |
| 8. Copy button exports data to clipboard | ✅ **COMPLETE** | `navigator.clipboard.writeText()` exports tab-separated table data |
| 9. Errors are handled gracefully | ✅ **COMPLETE** | All error paths have proper handling (see error table above) |
| 10. Progress is logged to console | ✅ **COMPLETE** | Console messages at each major step with emoji indicators |

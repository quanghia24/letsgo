# Product Comparison Tool

A Go command-line tool that compares product search results across three sources: local RapidAPI data, AliHunter API, and AliExpress API.

Features concurrent processing with goroutines for fast performance.

**Usage:**

```Bash
go run . -local <products.json> -output <report.html>
```

---

## 📁 Code Structure

### Project Structure

```Text
gofun/
├── main.go                      # Entry point with concurrent processing
├── products.json                # Input data
├── internal/
│   ├── alihunter/
│   │   └── alihunter.go        # AliHunter API client
│   ├── rapidapi/
│   │   └── aliexpress.go       # AliExpress/RapidAPI client
│   ├── model/
│   │   └── models.go           # Data structures
│   ├── report/
│   │   └── report.go           # Report generation logic
│   └── config/
│       ├── rapidapi.go         # API configuration
│       └── .env                # Environment variables
└── templates/
    └── report.tmpl             # HTML template with JS
```

## 🔄 Process Flow

```
1. INITIALIZATION
   ├─ Parse command-line flags (-local, -output)
   └─ Read and unmarshal input JSON file
          ↓
2. CONCURRENT DATA PROCESSING
   ├─ Launch goroutines for each product
   │  ├─ For each product:
   │  │  ├─ Fetch AliHunter API (goroutine)
   │  │  ├─ Fetch AliExpress API (goroutine)
   │  │  └─ Wait for both to complete
   │  ├─ Filter top 3 products with valid images
   │  └─ Build Report struct with comparison data
   ├─ Collect results via channel
   └─ Preserve original order
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

---

## 🔄 Data Flow Diagram

### 1️⃣ Input Phase
```
┌─────────────────┐
│  products.json  │  ← User provides local RapidAPI product data
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   main.go       │  ← Reads file, unmarshals JSON into []RapidapiResponse
└────────┬────────┘
         │
         ▼
  [Product Array]
```

### 2️⃣ Concurrent Processing Phase
```
                    ┌──────────────────────────────┐
                    │  Goroutine Pool (Max 5)      │
                    │  + Semaphore Rate Limiting   │
                    └──────────┬───────────────────┘
                               │
           ┌───────────────────┼───────────────────┐
           ▼                   ▼                   ▼
    ┌──────────┐        ┌──────────┐        ┌──────────┐
    │Product #1│        │Product #2│        │Product #N│
    └─────┬────┘        └─────┬────┘        └─────┬────┘
          │                   │                   │
     [Parallel API Calls]  [Parallel API Calls]  [...]
          │                   │
    ┌─────┴─────┐       ┌─────┴─────┐
    ▼           ▼       ▼           ▼
┌────────┐  ┌────────┐
│AliHunt │  │AliExpr │  ← Both APIs called simultaneously
│ API    │  │ API    │     using goroutines
└───┬────┘  └────┬───┘
    │            │
    └─────┬──────┘
          ▼
    [WaitGroup.Wait()]  ← Wait for both responses
          │
          ▼
    ┌─────────────┐
    │ Build Report│  ← Combine: Local + AliHunter + AliExpress
    │   Object    │
    └──────┬──────┘
           │
           ▼
    [Send to Channel] ← Results sent via buffered channel
```

### 3️⃣ Collection & Ordering Phase
```
┌─────────────────┐
│  Results Chan   │  ← All goroutines send results here
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Order by Index  │  ← Preserve original product order
└────────┬────────┘
         │
         ▼
[Ordered Comparisons Array]
```

### 4️⃣ Report Generation Phase
```
[Comparisons Array]
         │
         ▼
┌─────────────────────┐
│ report.go           │
│ ├─ Create view model│  ← ListReports{GeneratedAt, Comparisons}
│ ├─ Load template    │  ← Parse report.tmpl
│ ├─ Execute template │  ← Inject data into HTML
│ └─ Write to file    │  ← Save report.html
└──────────┬──────────┘
           │
           ▼
    ┌────────────┐
    │report.html │  ← Final interactive report
    └────────────┘
```

### 5️⃣ User Interaction Phase (Browser)
```
┌────────────────────┐
│   report.html      │
│  opened in browser │
└─────────┬──────────┘
          │
          ▼
┌─────────────────────┐
│ User Actions:       │
│ ├─ Check/uncheck    │  ← Mark matching products
│ │   product boxes   │
│ ├─ View statistics  │  ← Real-time counts per source
│ └─ Export JSON      │  ← Download with "matching" flags
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│ comparison.json     │  ← Exported data with user selections
└─────────────────────┘
```

---

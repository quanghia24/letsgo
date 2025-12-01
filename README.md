# Product Comparison Tool

Tool that compares product search results: local RapidAPI data, AliHunter API, and AliExpress API.

## 🚀 Quick Start

**Generate comparison data:**

```bash
go run . -local products.json
```

**Create HTML report:**

```bash
go run . -html true
```

**Usage workflow:**

1. Generate data → Open `report.html` → Mark matches → Export results

## 📁 Project Structure

```
gofun/
├── main.go                     # CLI entry point with concurrent processing
├── internal/
│   ├── alihunter/alihunter.go  # AliHunter API client
│   ├── rapidapi/aliexpress.go  # AliExpress/RapidAPI client
│   ├── model/models.go         # Data structures
│   ├── report/report.go        # Report generation & review fetching
│   └── config/rapidapi.go      # API configuration
└── templates/report.tmpl       # HTML template with JavaScript
```

## 🔄 Command Options

| Flag | Description | Example |
|------|-------------|---------|
| `-local <file>` | Input JSON file path | `go run . -local products.json` |
| `-html true` | Generate HTML from existing report.json | `go run . -html true` |

## 🏗️ Architecture Overview

### Data Generation Phase

```Text
Input JSON → Concurrent API Calls → Review Fetching → JSON Report
     ↓              ↓                     ↓              ↓
Products → [AliHunter + AliExpress] → Review Counts → report.json
```

### Key Features

**💾 Data Flow:**

1. **Generate**: `go run . -local products.json` → Creates fresh `report.json` with API results
2. **Visualize**: `go run . -html true` → Creates `report.html` for analysis
3. **Analyze**: Open browser → Mark products → Export selections
4. **Share**: Send exported JSON + HTML

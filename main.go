package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strings"
)

type Allowance struct {
    AllowanceType string `json:"allowanceType"`
    Amount        int64  `json:"amount"`
}

type TaxRequest struct {
    TotalIncome int64       `json:"totalIncome"`
    WHT         int64       `json:"wht"`
    Allowances  []Allowance `json:"allowances"`
}

type TaxLevel struct {
    Level string `json:"level"`
    Tax   int64  `json:"tax"`
}

type TaxResponse struct {
    Tax      int64      `json:"tax"`
    TaxLevel []TaxLevel `json:"taxLevel,omitempty"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}

const personalAllowance int64 = 60000
const donationMax int64 = 100000

// ==============================
//     TAX CALCULATION
// ==============================
func calculateTax(taxableIncome int64) (int64, []TaxLevel) {

    type bracket struct {
        start int64
        end   int64
        rate  float64
        label string
    }

    brackets := []bracket{
        {0, 150000, 0.00, "0-150,000"},
        {150001, 500000, 0.10, "150,001-500,000"},
        {500001, 1000000, 0.15, "500,001-1,000,000"},
        {1000001, 2000000, 0.20, "1,000,001-2,000,000"},
        {2000001, 1<<60 - 1, 0.35, "2,000,001 ขึ้นไป"},
    }

    var totalTax int64
    var levels []TaxLevel

    for _, b := range brackets {
        if taxableIncome < b.start {
            levels = append(levels, TaxLevel{b.label, 0})
            continue
        }

        taxable := min(taxableIncome, b.end) - b.start + 1
        if taxable < 0 {
            taxable = 0
        }

        tax := int64(float64(taxable) * b.rate)
        totalTax += tax

        levels = append(levels, TaxLevel{b.label, tax})
    }

    return totalTax, levels
}

func min(a, b int64) int64 {
    if a < b {
        return a
    }
    return b
}

// ==============================
//          API HANDLER
// ==============================
func taxHandler(w http.ResponseWriter, r *http.Request) {

    if r.Method != "POST" {
        w.WriteHeader(http.StatusMethodNotAllowed)
        json.NewEncoder(w).Encode(ErrorResponse{"method not allowed"})
        return
    }

    var req TaxRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(ErrorResponse{"invalid JSON body"})
        return
    }

    // Validation
    if req.TotalIncome < 0 {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(ErrorResponse{"totalIncome must be a positive number"})
        return
    }

    if req.WHT < 0 {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(ErrorResponse{"wht must be a positive number"})
        return
    }

    if req.WHT > req.TotalIncome {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(ErrorResponse{"wht cannot be greater than totalIncome"})
        return
    }

    // DONATION (max 100k)
    donation := int64(0)
    for _, a := range req.Allowances {
        if strings.ToLower(strings.TrimSpace(a.AllowanceType)) == "donation" {
            if a.Amount > donationMax {
                donation += donationMax
            } else if a.Amount > 0 {
                donation += a.Amount
            }
        }
    }
    if donation > donationMax {
        donation = donationMax
    }

    // TAXABLE INCOME
    taxableIncome := req.TotalIncome - personalAllowance - donation
    if taxableIncome < 0 {
        taxableIncome = 0
    }

    // Calculate TAX
    totalTax, levels := calculateTax(taxableIncome)

    // Final tax after WHT
    finalTax := totalTax - req.WHT

    resp := TaxResponse{
        Tax:      finalTax,
        TaxLevel: levels,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func main() {
    http.HandleFunc("/tax/calculations", taxHandler)
    log.Println("Server running at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}

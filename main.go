package main

//=========================
// IMPORTS
//=========================

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//=========================
// MAIN
//=========================

func main() {
	dir, err := GetArg(1)
	if err != nil {
		panic(err)
	}

	if !DirExists(dir) {
		panic(fmt.Errorf("%s directory does not exist", dir))
	}

	mp := make(map[string]DollarAmount)
	receipts := []ReceiptPdf{}

	err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".pdf" {
			return nil
		}
		fileName := strings.TrimSuffix(filepath.Base(path), ".pdf")
		r, err := NewReceiptPdfFromFileName(fileName)
		if err != nil {
			return err
		}
		cost, err := r.GetCost()
		if err != nil {
			return err
		}
		category := r.GetCategory()
		if !IsKeyInMap(mp, category) {
			mp[category] = cost
		} else {
			mp[category] = mp[category].Add(cost)
		}
		receipts = append(receipts, r)
		return nil
	})
	if err != nil {
		panic(err)
	}

	outPath, err := GetArg(2)
	if err != nil {
		panic(err)
	}

	invoiceName, err := GetArg(3)
	if err != nil {
		panic(err)
	}

	receiptMap := make(map[string][]ReceiptPdf)
	for _, r := range receipts {
		category := r.GetCategory()
		receiptMap[category] = append(receiptMap[category], r)
	}

	var output strings.Builder

	output.WriteString(strings.ToUpper(invoiceName) + "\n")

	grandTotal := DollarAmount{}
	for _, amt := range mp {
		grandTotal = grandTotal.Add(amt)
	}
	output.WriteString("TOTAL: " + grandTotal.String() + "\n\n")

	for category, total := range mp {
		output.WriteString(fmt.Sprintf("%s %s:\n", strings.ToUpper(category), total.String()))
		for _, receipt := range receiptMap[category] {
			output.WriteString("\t" + receipt.String() + "\n")
		}
		output.WriteString("\n\n")
	}

	err = os.WriteFile(outPath, []byte(output.String()), 0644)
	if err != nil {
		panic(fmt.Errorf("failed to write output file: %w", err))
	}

	fmt.Println("Wrote output to", outPath)
}

//=========================
// RECEIPT PDF
//=========================

type ReceiptPdf struct {
	FileName string
	Parts    []string
}

func NewReceiptPdfFromFileName(fileName string) (ReceiptPdf, error) {
	r := &ReceiptPdf{FileName: fileName}
	r.Parts = strings.Split(fileName, "-")
	if len(r.Parts) != 6 {
		return *r, fmt.Errorf("each PdfReceipt requires 6 filename parts split by '-', like:\n041525-target-29.87-desc-category-store.pdf")
	}
	return *r, nil
}

func (r *ReceiptPdf) GetCost() (DollarAmount, error) {
	return NewDollarAmountFromString(r.Parts[2])
}

func (r *ReceiptPdf) GetCategory() string {
	return r.Parts[4]
}

func (r *ReceiptPdf) String() string {
	return fmt.Sprintf("[%s] [%s] [%s] => %s", r.Parts[0], r.Parts[1], r.Parts[3], r.Parts[2])
}

//=========================
// UTILS
//=========================

func GetArg(i int) (string, error) {
	if len(os.Args) > i {
		return os.Args[i], nil
	}
	return "", fmt.Errorf("arg of index '%d' does not exist", i)
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func ParseToFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	return strconv.ParseFloat(s, 64)
}

func IsKeyInMap[K comparable, V any](m map[K]V, key K) bool {
	_, ok := m[key]
	return ok
}

func RoundUpToTwoDecimals(f float64) float64 {
	return math.Ceil(f*100) / 100
}

//=========================
// DOLLAR AMOUNT
//=========================

type DollarAmount struct {
	Dollars int
	Cents   int
}

func NewDollarAmountFromString(s string) (DollarAmount, error) {
	s = strings.TrimSpace(s)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return DollarAmount{}, err
	}
	totalCents := int(math.Round(f * 100))
	return DollarAmount{
		Dollars: totalCents / 100,
		Cents:   totalCents % 100,
	}, nil
}

func (d DollarAmount) Add(other DollarAmount) DollarAmount {
	totalCents := d.TotalCents() + other.TotalCents()
	return DollarAmount{
		Dollars: totalCents / 100,
		Cents:   totalCents % 100,
	}
}

func (d DollarAmount) TotalCents() int {
	return d.Dollars*100 + d.Cents
}

func (d DollarAmount) ToFloat() float64 {
	return float64(d.TotalCents()) / 100
}

func (d DollarAmount) String() string {
	return fmt.Sprintf("$%d.%02d", d.Dollars, d.Cents)
}

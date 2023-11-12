package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/gofpdi"
)

// for when we upgrade to a db rather than csv
//type Pet struct{
// 	SpeciesGenderCode string //`csv:"Pet1Code"`
// 	Age string //`csv:"Pet1Age"`
// 	WeightEst string //`csv:"Pet1Weight"`
// 	Name string //`csv:"Pet1Name"`
// }
// type Customer struct{
// 	Own Owner
// 	Pets []Pet
// }
type Owner struct{
	Name string `csv:"Name"`
	Phone string `csv:"Phone"`
	AltPhone string `csv:"AltPhone"`
	StreetAddress string `csv:"Address"`
	City string `csv:"City State Zip"` 
	Email string `csv:"Email"`
	Pet1 string `csv:"Pet1"`
	Pet2 string `csv:"Pet2"`
	Pet3 string `csv:"Pet3"`
	Pet4 string `csv:"Pet4"`
	Pet5 string `csv:"Pet5"`
}
func main() {
	// allow user to input starting ID number. Expected format defined by owner
	var startID []byte
	fmt.Println("Please enter the starting Invoice ID")
	fmt.Scan(&startID)
	matched, err := regexp.Match(`[A-Z][0-9]*`, startID)
	if !matched || err != nil{
		fmt.Println("Error: Invalid ID format. Expected exactly 1 capital letter followed by a number.")
		os.Exit(1)
	}

	// split prefix letter form ID number
	prefix := string(startID[0])
	ID, _ := strconv.Atoi(string(startID[1:])) 

	// load input csv into memory
	bs, err := os.ReadFile(".//Input/Invoice_Input.csv")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(2)
	}
	// check for voucher ID to avoid overwriting
	// TODO

	// Parse and unmarshal CSV to Owner struct
	r := csv.NewReader(bytes.NewReader(bs)) 
	r.Comment = '/'
	r.LazyQuotes = true
	owners := []*Owner{}
	err = gocsv.UnmarshalCSV(r, &owners)
	if err != nil && err != io.EOF{
		log.Fatal(err)
	}

	// Generate PDF Invoices
	for _, owner := range owners{
		// Each pet gets its own invoice
		for i := 0; i < owner.getNumPets(); i++{
		generateInvoice(*owner, prefix, ID)
		ID++
		}
	}
}

func generateInvoice(o Owner, prefix string, DocID int) error{
	// generate a new document
	pdf := gofpdf.New("P", "mm", "A4", "")
	// Import Invoice pdf with gofpdi free pdf document importer
	tpl1 := gofpdi.ImportPage(pdf, ".//Input/BLANK_Voucher_Current.pdf", 1, "/MediaBox")
	pdf.AddPage()
	
	// Draw imported template onto page
	gofpdi.UseImportedTemplate(pdf, tpl1, 0, 5, 210, 0)
	
	// Draw Voucher Header, included ID number and Date printed
	pdf.SetFont("Helvetica", "B", 16)
	pdf.SetTextColor(255,0,0) // Red text
	draw(pdf, 57, 82, prefix+strconv.Itoa(DocID))
	year, month, day := time.Now().Date() 
	draw(pdf, 40, 90, strconv.Itoa(int(month))+"/"+strconv.Itoa(day)+"/"+strconv.Itoa(year))
	// Draw Customer data
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(0, 0, 0) // Black text
	o.fillCustomerData(pdf, 25, 100)

	// Saves voucher as PDF
	err := os.MkdirAll("Output/", 0750) 
	if err != nil{
		fmt.Println("Error: Unable to create Output folder.", err)
		os.Exit(3)
	}
	err = pdf.OutputFileAndClose("Output/"+prefix+strconv.Itoa(DocID)+".pdf")
	if err != nil {
		panic(err)
	}
	return nil
}

// draws text at specified coordiantes on the PDF
func draw(p *gofpdf.Fpdf, x, y float64, data string){
	p.SetXY(x, y)
	p.Cell(100, 0, data)
}

// each pet gets its own voucher, so we need to know how many pets an owner has
// using reflection allows us to check for the existance of a pet iteratively
func (o Owner) getNumPets() int {
	numpets := 0
	vals := reflect.ValueOf(o)
	// Loop through elements 8-end AKA all Pet fields
	for i := 6; i < vals.NumField(); i++ {
		if !vals.Field(i).IsZero(){
			numpets++
		}
	}
	return numpets
}

//draws each element of the Owner struct if and only if data exists in any given field
func (o Owner) fillCustomerData(p *gofpdf.Fpdf, xStart, yStart float64){
	yOffset := 0
	vals := reflect.ValueOf(o)
	for i := 0; i < vals.NumField(); i++ {
		if !vals.Field(i).IsZero(){
			draw(p, 25, 100+float64(yOffset), vals.Field(i).String())
			yOffset += 5
		} 
	}
}
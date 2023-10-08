package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

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
	City string `csv:"City"` 
	St string `csv:"State"`
	Zip string `csv:"Zip"`
	Email string `csv:"Email"`
	Pet1 string `csv:"Pet1"`
	Pet2 string `csv:"Pet2"`
	Pet3 string `csv:"Pet3"`
	Pet4 string `csv:"Pet4"`
	Pet5 string `csv:"Pet5"`
}
func main() {
	 bs, err := os.ReadFile("Sample_DataSimple.csv")
	//bs, err := os.ReadFile(os.Args[1])
	//startID := strings.TrimSuffix(os.Args[1], ".csv")
	startID := 1
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	r := csv.NewReader(bytes.NewReader(bs))
	r.Comment = '/'
	r.LazyQuotes = true
	owners := []*Owner{}
	err = gocsv.UnmarshalCSV(r, &owners)
	if err != nil && err != io.EOF{
		log.Fatal(err)
	}
	for _, owner := range owners{
		//fmt.Println(owner)
		GenerateInvoice(*owner, startID)
		startID++
	}
}

func GenerateInvoice(o Owner, DocID int) error{
	// generate a new document
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Import Invoice pdf with gofpdi free pdf document importer
	tpl1 := gofpdi.ImportPage(pdf, "Voucher_BLANK_TAC_2023_Sample.pdf", 1, "/MediaBox")
	pdf.AddPage()
	
	// Draw imported template onto page
	gofpdi.UseImportedTemplate(pdf, tpl1, 0, 5, 210, 0)

	// Draw Customer data
 	pdf.SetFont("Helvetica", "", 20)
	pdf.Cell(0, 0,o.Name)

	err := pdf.OutputFileAndClose("Output/Voucher_"+strconv.Itoa(DocID)+".pdf")
	if err != nil {
		panic(err)
	}
	return nil
}
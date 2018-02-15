package main

import (
	"fmt"
	"log"
	"os"

	"github.com/makpoc/hades-sheet/gsclient"
)

func main() {
	spreadsheetId, ok := os.LookupEnv("SHEET_ID")
	if !ok {
		log.Fatal("Please set SHEET_ID in environment")
	}
	sheet, err := gsclient.New(spreadsheetId)
	if err != nil {
		log.Fatal("Failed to create spreadsheet client", err)
	}
	result, err := sheet.GetTimeZones()
	if err != nil {
		log.Fatal("GetTimeZones failed", err)
	}

	fmt.Printf("Success\n%v\n", result)
}

package zipCodesService

import (
	"fmt"
	"log"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type ZipCodeNode struct {
	ZipCode string
	Lat float64
	Lng float64
	City string
	State string
}

type ZipCodes struct {
	DatasetList map[string]ZipCodeNode
}

func New(datasetPath string) (*ZipCodes, error) {
	zipcodes, err := LoadDataset(datasetPath)
	if err != nil {
		return nil, err
	}
	return &zipcodes, nil
}

func LoadDataset(dataPath string)(ZipCodes, error)  {
	file, err := excelize.OpenFile(dataPath)
	rows, err := file.Rows("Sheet1")
	if err != nil {
		log.Fatal(err)
		return ZipCodes{}, fmt.Errorf("zipcodes: error while opening file %v", err)
	}
	for rows.Next() {
		row := rows.Columns()
		fmt.Printf("%s\t%s\n", row[1], row[3]) // Print values in columns B and D
	}
	return ZipCodes{}, nil
}


package zipCodesService

import (
	"fmt"
	"log"
	"strconv"

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

func NewZipCodesDataService(datasetPath string) (*ZipCodes, error) {
	zipcodes, err := loadDataset(datasetPath)
	if err != nil {
		return nil, err
	}
	return &zipcodes, nil
}

func loadDataset(dataPath string)(ZipCodes, error)  {
	file, err := excelize.OpenFile(dataPath)
	if err != nil {
		log.Fatal(err)
		return ZipCodes{}, fmt.Errorf("zipcodes: error while opening file %v", err)
	}
	rows, error := file.Rows("Sheet1")
	if error != nil {
		log.Fatal(error)
		return ZipCodes{}, fmt.Errorf("zipcodes: error while opening spreadsheet %v", error)
	}
	zipCodeMap := ZipCodes{DatasetList: make(map[string]ZipCodeNode)}

	for rows.Next() {
		row, _ := rows.Columns()
		lat,_ := strconv.ParseFloat(row[1], 64)
		lng,_ := strconv.ParseFloat(row[2], 64)
		zipCodeMap.DatasetList[row[0]] = ZipCodeNode{
			ZipCode: row[0],
			Lat:     lat,
			Lng:     lng,
			City:    row[3],
			State:   row[4],
		}
	}
	return zipCodeMap, err
}


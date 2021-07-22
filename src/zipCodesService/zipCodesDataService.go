package zipCodesService

import (
	"fmt"
	"log"
	"os"

	"github.com/gocarina/gocsv"
)

type ZipCodeNode struct {
	Zip string `csv:"zip"`
	Lat float64 `csv:"lat"`
	Lng float64 `csv:"lng"`
	City string `csv:"city"`
	State string `csv:"state_id"`
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
	in, err := os.Open(dataPath)
	if err != nil {
		panic(err)
	}
	defer in.Close()
	list := []*ZipCodeNode{}
	zipCodeMap := ZipCodes{DatasetList: make(map[string]ZipCodeNode)}
	if err := gocsv.UnmarshalFile(in, &list); err != nil {
		log.Fatal(err)
		return ZipCodes{}, fmt.Errorf("zipcodes: error while opening spreadsheet %v", err)
	}
	for _, v := range list {
		zipCodeMap.DatasetList[v.Zip] = *v
	}
	return zipCodeMap, err
}


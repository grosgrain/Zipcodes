package zipCodesService

import (
	"fmt"
	"log"
)

type RequestService struct {
	ZipCodesDataService *ZipCodes
	Error error
}

func NewRequestService() *RequestService {
	s := new(RequestService)
	s.ZipCodesDataService,s.Error = NewZipCodesDataService("../public/assets/uszips.csv")
	return s
}

func (s *RequestService)GetAllZipCodesData()(map[string]ZipCodeNode, error) {
	zipCodesDataset,err := s.ZipCodesDataService, s.Error
	if err != nil {
		log.Fatal(err)
	}
	return zipCodesDataset.DatasetList, err
}

func (s *RequestService)Lookup(zip string)(*ZipCodeNode, error) {
	foundedZipCode := s.ZipCodesDataService.DatasetList[zip]
	if (foundedZipCode == ZipCodeNode{}) {
		return &ZipCodeNode{}, fmt.Errorf("zipcodes: zipcode %s not found !", zip)
	}
	return &foundedZipCode, nil
}

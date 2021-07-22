package zipCodesService

import (
	"log"
)

type RequestService struct {
	ZipCodesDataService *ZipCodes
	Error error
}

func NewRequestService() *RequestService {
	s := new(RequestService)
	s.ZipCodesDataService,s.Error = NewZipCodesDataService("../public/assets/uszips.xlsx")
	return s
}

func (s *RequestService)GetAllZipCodesData()(map[string]ZipCodeNode, error) {
	zipCodesDataset,err := s.ZipCodesDataService, s.Error
	if err != nil {
		log.Fatal(err)
	}
	return zipCodesDataset.DatasetList, err
}

func Lookup() {

}

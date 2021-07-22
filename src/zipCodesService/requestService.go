package zipCodesService

import (
	"fmt"
	"log"
	"math"
)

const (
	earthRadiusKm = 6371
	earthRadiusMi = 3958
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

func (s *RequestService) GetZipCodesWithinRadius(zipCode string, radius float64, inMiles bool) ([]string, error) {
	var earthRadius float64
	if (inMiles) {
		earthRadius = earthRadiusMi
	} else {
		earthRadius = earthRadiusKm
	}
	list := []string{}
	location, err := s.Lookup(zipCode)
	if err != nil {
		return list, err
	}
	for _, elm := range s.ZipCodesDataService.DatasetList {
		if elm.Zip != location.Zip {
			distance := distanceBetweenPoints(location.Lat, location.Lng, elm.Lat, elm.Lng, earthRadius)
			if distance < radius {
				list = append(list, elm.Zip)
			}
		}
	}
	return list, nil
}

// DistanceBetweenPoints returns the distance between two lat/lng points using the Haversin distance formula.
func distanceBetweenPoints(latitude1, longitude1, latitude2, longitude2 float64, radius float64) float64 {
	lat1 := degreesToRadians(latitude1)
	lon1 := degreesToRadians(longitude1)
	lat2 := degreesToRadians(latitude2)
	lon2 := degreesToRadians(longitude2)
	diffLat := lat2 - lat1
	diffLon := lon2 - lon1

	a := hsin(diffLat) + math.Cos(lat1) * math.Cos(lat2) * hsin(diffLon)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1 - a))
	distance := c * radius

	return math.Round(distance*100) / 100
}

func hsin(t float64) float64 {
	return math.Pow(math.Sin(t / 2), 2)
}

func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}


package zipCodesService

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"sync"

	"googlemaps.github.io/maps"
)

const (
	earthRadiusKm = 6371
	earthRadiusMi = 3958
	meterInMiles = 1609.34
)

type ZipToDistanceMapping struct {
	Zip string
	OriginZip string
	Distance float32
}

type RequestService struct {
	ZipCodesDataService *ZipCodes
	Error error
}

func NewRequestService() *RequestService {
	s := new(RequestService)
	pwd, _ := os.Getwd()
	path := pwd + "/public/assets/uszips.csv"
	s.ZipCodesDataService,s.Error = NewZipCodesDataService(path)
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
	list = append(list, zipCode)
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

func (s *RequestService) GetDistanceFromOnePointToAnother(originZip string, destZip string) (ZipToDistanceMapping, error)  {
    googleMatrixAPIKey,_ := os.LookupEnv("GOOGLE_DISTANCE_MATRIX_API_KEY")
	c, err := maps.NewClient(maps.WithAPIKey(googleMatrixAPIKey))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	r := &maps.DistanceMatrixRequest{
		Origins:      []string{originZip},
		Destinations: []string{destZip},
		Units:        maps.UnitsImperial,
		Language:     "en",
	}
	route, err := c.DistanceMatrix(context.Background(), r)

	if err != nil {
		log.Fatalf("Fatal error: %s", err)
	}
    res := ZipToDistanceMapping{
		Zip:       destZip,
		OriginZip: originZip,
		Distance:  float32(route.Rows[0].Elements[0].Distance.Meters) / meterInMiles,
	}
	return res, nil
}

func (s *RequestService) GetDistanceFromOnePointToMultiplePoints(originZip string, destZips []string, wg *sync.WaitGroup) ([]ZipToDistanceMapping, error)  {
	if originZip == "" || destZips == nil || len(destZips) <= 0 {
		log.Fatal("Inputs are not valid")
		return nil, nil
	}
	googleMatrixAPIKey,_ := os.LookupEnv("GOOGLE_DISTANCE_MATRIX_API_KEY")
	c, err := maps.NewClient(maps.WithAPIKey(googleMatrixAPIKey))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	r := &maps.DistanceMatrixRequest{
		Origins:      []string{originZip},
		Destinations: destZips,
		Units:        maps.UnitsImperial,
		Language:     "en",
	}
	route, err := c.DistanceMatrix(context.Background(), r)

	if err != nil {
		log.Fatalf("Fatal error from Google Distance Matrix API: %s", err)
	}
	res := make([]ZipToDistanceMapping, len(destZips))
	for i, destZip := range destZips {
		res[i] = ZipToDistanceMapping{
			Zip:       destZip,
			OriginZip: originZip,
			Distance:  float32(route.Rows[0].Elements[i].Distance.Meters) / meterInMiles,
		}
	}
	defer wg.Done()
	return res, nil
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


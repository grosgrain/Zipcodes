package zipCodesService

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

const googleMaximumZipsPerRequest  = 20

func GetAllZipCodesData(w http.ResponseWriter, r *http.Request) {
	requestService := NewRequestService()
	if requestService.Error != nil {
		log.Fatal(requestService.Error)
		return
	}
	data,err := requestService.GetAllZipCodesData()
	if err != nil {
		log.Fatal(err)
		return
	}
	json.NewEncoder(w).Encode(data)
	return
}

func Lookup(w http.ResponseWriter, r *http.Request) {
	zip := validateZip(w, r)
	c := make(chan *ZipCodeNode, 20)
	go func() {
		data, _ := NewRequestService().Lookup(*zip)
		c <- data
	}()
	json.NewEncoder(w).Encode(<-c)
	return
}

func GetZipCodesWithinRadius(w http.ResponseWriter, r *http.Request)  {
	zip, radius, inMiles := validateZipRadiusMatrix(w, r)
	c := make(chan []string, 20)
	go func() {
		data, _ := NewRequestService().GetZipCodesWithinRadius(*zip, *radius, *inMiles)
		c <- data
	}()
	json.NewEncoder(w).Encode(<-c)
	return
}

func GetDistanceFromOneZipToMultipleZips(w http.ResponseWriter, r *http.Request)  {
	zip, radius, inMiles := validateZipRadiusMatrix(w, r)
	c := make(chan []string, 20)
	go func() {
		data, _ := NewRequestService().GetZipCodesWithinRadius(*zip, *radius, *inMiles)
		c <- data
	}()
	zipLists := <-c
	var lists []ZipToDistanceMapping
	var wg sync.WaitGroup
	for i := 0; i < len(zipLists) / googleMaximumZipsPerRequest + 1; i++ {
		var end int
		if i * googleMaximumZipsPerRequest + googleMaximumZipsPerRequest < len(zipLists) {
			end = i * googleMaximumZipsPerRequest + googleMaximumZipsPerRequest
		} else {
			end = len(zipLists)
		}
		temp := zipLists[i * googleMaximumZipsPerRequest : end]
		wg.Add(1)
		go func() {
			data,_ := NewRequestService().GetDistanceFromOnePointToMultiplePoints(*zip, temp, &wg)
			lists = append(lists, data...)
		}()
	}
	wg.Wait()
	for _, v:= range lists {
		if float64(v.Distance) <= *radius {
			res := make(map[string]float32)
			res[v.Zip] = v.Distance
			json.NewEncoder(w).Encode(res)
		}
	}
	return
}

func validateZipRadiusMatrix(w http.ResponseWriter, r *http.Request)(*string, *float64, *bool)  {
	body := json.NewDecoder(r.Body)
	body.DisallowUnknownFields()
	temp := struct {
		Zip *string `json:"zip"`
		Radius *float64 `json:"radius"`
		InMiles *bool `json:"inMiles"`
	}{}
	err := body.Decode(&temp)
	if err != nil {
		// bad JSON or unrecognized json field
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, nil, nil
	}
	if temp.Zip == nil {
		http.Error(w, "missing field 'zip' from JSON object", http.StatusBadRequest)
		return nil, nil, nil
	} else if temp.Radius == nil {
		http.Error(w, "missing field 'radius' from JSON object", http.StatusBadRequest)
		return nil, nil, nil
	} else if temp.InMiles == nil {
		http.Error(w, "missing field 'inMiles' from JSON object", http.StatusBadRequest)
		return nil, nil, nil
	}
	return temp.Zip, temp.Radius, temp.InMiles
}

func validateZip(w http.ResponseWriter, r *http.Request) *string {
	body := json.NewDecoder(r.Body)
	body.DisallowUnknownFields()
	temp := struct {
		Zip *string `json:"zip"` // pointer so we can test for field absence
	}{}
	err := body.Decode(&temp)
	if err != nil {
		// bad JSON or unrecognized json field
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	if temp.Zip == nil {
		http.Error(w, "missing field 'zip' from JSON object", http.StatusBadRequest)
		return nil
	}
	return temp.Zip
}

package freightMatchingService

import (
	"encoding/json"
	"log"
	"net/http"
)


func GetNearbyCarriers(w http.ResponseWriter, r *http.Request) {
	requestService := NewRequestService()
	if requestService.Error != nil {
		log.Fatal("Failed to set up freightMatchingService", requestService.Error)
		return
	}
	zip, radius, page := validateZipRadiusPage(w, r)
	data,err := requestService.GetNearbyCarriers(*zip, *radius, *page)
	if err != nil {
		log.Fatal("Failed to fetch nearby carriers", err)
		return
	}
	json.NewEncoder(w).Encode(data)
	return
}

func validateZipRadiusPage(w http.ResponseWriter, r *http.Request)(*string, *float64, *int)  {
	body := json.NewDecoder(r.Body)
	body.DisallowUnknownFields()
	temp := struct {
		Zip *string `json:"zip"`
		Radius *float64 `json:"radius"`
		Page *int `json:"page"`
	}{}
	err := body.Decode(&temp)
	if err != nil {
		// bad JSON or unrecognized json field
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, nil, nil
	}
	if temp.Zip == nil {
		http.Error(w, "missing field 'zip' from JSON object", http.StatusBadRequest)
	} else if temp.Radius == nil {
		http.Error(w, "missing field 'radius' from JSON object", http.StatusBadRequest)
	} else if temp.Page == nil {
		http.Error(w, "missing field 'page' from JSON object", http.StatusBadRequest)
	}
	return temp.Zip, temp.Radius, temp.Page
}

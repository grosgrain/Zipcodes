package zipCodesService

import (
	"encoding/json"
	"log"
	"net/http"
)

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
	body := json.NewDecoder(r.Body)
	body.DisallowUnknownFields()
	temp := struct {
		Zip *string `json:"zip"` // pointer so we can test for field absence
	}{}
	err := body.Decode(&temp)
	if err != nil {
		// bad JSON or unrecognized json field
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if temp.Zip == nil {
		http.Error(w, "missing field 'test' from JSON object", http.StatusBadRequest)
		return
	}
	c := make(chan *ZipCodeNode, 20)
	go func() {
		data, _ := NewRequestService().Lookup(*temp.Zip)
		c <- data
	}()
	/*requestService := NewRequestService()
	data,err := requestService.Lookup(*temp.Zip)
	if err != nil {
		log.Fatal(err)
		return
	}*/
	json.NewEncoder(w).Encode(<-c)
	return
}

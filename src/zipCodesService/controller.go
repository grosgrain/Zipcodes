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
	json.NewEncoder(w).Encode(requestService.ZipCodesDataService.DatasetList)
	return
}

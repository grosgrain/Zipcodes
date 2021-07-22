package zipCodesService

import (
	"log"
	"net/http"
	"encoding/json"
)

func GetAllZipCodesData(w http.ResponseWriter, r *http.Request) {
	zipCodesDataset,err := New("../public/assets/uszips.xlsx")
	if err != nil {
		log.Fatal(err)
		return
	}
	json.NewEncoder(w).Encode(zipCodesDataset.DatasetList)
	return
}

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/grosgrain/zipcodes/src/zipCodesService"
	"github.com/grosgrain/zipcodes/src/freightMatchingService"
	"github.com/joho/godotenv"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	port, exists := os.LookupEnv("ZIP_CODE_SERVICE_PORT")
	if !exists {
		port = ":4001"
	}

	//Freight Matching route
	freightMatchingRouter := router.PathPrefix("/zipCode").Subrouter()
	freightMatchingRouter.HandleFunc("/allZipCodes", zipCodesService.GetAllZipCodesData).Methods("GET")
	freightMatchingRouter.HandleFunc("/zipCodeLookup", zipCodesService.Lookup).Methods("POST")
	freightMatchingRouter.HandleFunc("/zipCodesWithinRadius", zipCodesService.GetZipCodesWithinRadius).Methods("POST")
	freightMatchingRouter.HandleFunc("/distanceFromOneZipToMultipleZips", zipCodesService.GetDistanceFromOneZipToMultipleZips).Methods("POST")
	freightMatchingRouter.HandleFunc("/getNearbyCarriers", freightMatchingService.GetNearbyCarriers).Methods("POST")

	log.Fatal(http.ListenAndServe(port, router))

}

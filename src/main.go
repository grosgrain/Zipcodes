package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grosgrain/zipcodes/src/zipCodesService"
)


func main() {
	router := mux.NewRouter().StrictSlash(true)

	//Freight Matching route
	freightMatchingRouter := router.PathPrefix("/zipCode").Subrouter()
	freightMatchingRouter.HandleFunc("/allZipCodes", zipCodesService.GetAllZipCodesData).Methods("GET")
	freightMatchingRouter.HandleFunc("/zipCodeLookup", zipCodesService.Lookup).Methods("POST")
	freightMatchingRouter.HandleFunc("/zipCodesWithinRadius", zipCodesService.GetZipCodesWithinRadius).Methods("POST")
	log.Fatal(http.ListenAndServe(":4001", router))

}

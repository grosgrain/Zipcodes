package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	//"github.com/zipcodes/zipCodesService"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	//router.HandleFunc("/", homeLink)

	//Freight Matching route
	freightMatchingRouter := router.PathPrefix("/freightMatching").Subrouter()
	freightMatchingRouter.HandleFunc("/zipCodes", homeLink).Methods("POST")
	log.Fatal(http.ListenAndServe(":4001", router))

	zipCodesDataset,err := zipCodesService.New('public/assets/uszips.xlsx')
}

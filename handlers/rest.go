package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func SyncData(w http.ResponseWriter, req *http.Request) {
	dg, cancel := getDgraphClient()
	defer cancel()
	SetDataToGraph(dg)
	json.NewEncoder(w).Encode("Data download")
}
func predicateData(w http.ResponseWriter, req *http.Request) {
	dg, cancel := getDgraphClient()
	defer cancel()
	CreateRelations(dg)
	json.NewEncoder(w).Encode("Data update")
}
func Controller() {
	router := mux.NewRouter()
	router.HandleFunc("/people", SyncData).Methods("GET")
	router.HandleFunc("/v1", predicateData).Methods("GET")
	log.Fatal(http.ListenAndServe(":3030", router))

}

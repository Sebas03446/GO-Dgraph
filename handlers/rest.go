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
func ShowPurchaser(w http.ResponseWriter, req *http.Request) {
	dg, cancel := getDgraphClient()
	defer cancel()
	GetAllPurchaser(dg)
	json.NewEncoder(w).Encode("List succesful")
}
func Controller() {
	router := mux.NewRouter()
	router.HandleFunc("/v1", SyncData).Methods("GET")
	router.HandleFunc("/v1/purchaser", ShowPurchaser).Methods("GET")
	log.Fatal(http.ListenAndServe(":3030", router))

}

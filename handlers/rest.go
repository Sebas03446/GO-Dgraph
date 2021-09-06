package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func SyncData(w http.ResponseWriter, req *http.Request) {
	/*dg, cancel := getDgraphClient()
	defer cancel()
	SetDataToGraph(dg)
	json.NewEncoder(w).Encode("Data download")*/
	value, _ := io.ReadAll(req.Body)
	fmt.Println(string(value))
	fmt.Println("hola")
}
func ShowPurchaser(w http.ResponseWriter, req *http.Request) {
	dg, cancel := getDgraphClient()
	defer cancel()
	allClient := GetAllPurchaser(dg)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
	json.NewEncoder(w).Encode(allClient)
	/*client, err := json.Marshal(allClient)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	w.Write(client)*/
}
func Controller() {
	router := mux.NewRouter()
	router.HandleFunc("/v1", SyncData).Methods("POST")
	router.HandleFunc("/v1/purchaser", ShowPurchaser).Methods("GET")
	//log.Fatal(http.ListenAndServe(":3030", router))
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":3030", handler))

}

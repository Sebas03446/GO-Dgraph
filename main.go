package main

import (
	"CHALLENGE/handlers"
)

/*import (
	"CHALLENGE/resource"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)*/
/*
var people = resource.Buyer{}

func GetPeople(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(people)
}*/
func main() {
	handlers.Controller()
}

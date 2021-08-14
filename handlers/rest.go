package handlers

import (
	"CHALLENGE/resource"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

var people = resource.Buyer{}

func GetPeople(w http.ResponseWriter, req *http.Request) {
	people.Id = "prueba"
	people.Name = "jo"
	people.Age = 15
	people.TimeNow = time.Now().Unix()
	json.NewEncoder(w).Encode(people)
}
func Controller() {
	router := mux.NewRouter()
	dg, cancel := getDgraphClient()
	defer cancel()
	router.HandleFunc("/people", GetPeople).Methods("GET")
	Example_setObject(dg)
	log.Fatal(http.ListenAndServe(":3030", router))

}
func SetDataToGraph(dg *dgo.Dgraph) {
	var data = resource.GetBuyer()
	op := &api.Operation{}
	op.Schema = `type Buyer {
					id
					name
					age
					products
					time
				}
				type Product{
					id
					name
					price
					time
					buy
				}
				id: string @index(term).
				name: string @index(term).
				age: int @index(int).
				products:[uid].
				time: int @index(int).
				price: float @index(float). 
				buy: [uid].`
	ctx := context.Background()
	if err := dg.Alter(ctx, op); err != nil {
		log.Fatal(err)
	}

	mu := &api.Mutation{
		CommitNow: true,
	}
	pb, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	mu.SetJson = pb
	response, err := dg.NewTxn().Mutate(ctx, mu)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(response.Json))
	q := `{
		data(func: eq(name, "Beaner")) {
		  name
		  id
		  time
		  age
		  
		}
	  }
	  `
	//resp, err := dg.NewTxn().QueryWithVars(ctx, q, variables)
	resp, err := dg.NewTxn().Query(ctx, q)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(resp.Json))
	type Producto struct {
		Name  string  `dgraph:"name"`
		Price float64 `dgraph:"price"`
	}
	type Comprador struct {
		Name     string     `dgraph:"name"`
		Products []Producto `dgraph:"products"`
	}
	type Root struct {
		Data []resource.Buyer `dgraph:"data"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		log.Fatal(err)
	}

	out, _ := json.MarshalIndent(r, "", "\t")
	fmt.Printf("%s\n", out)

}
func Example_setObject(dg *dgo.Dgraph) {
}

func getDgraphClient() (*dgo.Dgraph, context.CancelFunc) {
	conn, err := grpc.Dial("localhost:9081", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	return dg, func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error while closing connection:%v", err)
		}
	}
}

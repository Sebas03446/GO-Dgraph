package handlers

import (
	"CHALLENGE/resource"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"google.golang.org/grpc"
)

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
func SetDataToGraph(dg *dgo.Dgraph) {
	var dataBuy = resource.GetBuyer()
	var dataProduct = resource.GetProduct()
	var dataTrans = resource.TransformTransaction()
	/*var dataPredicate = resource.GetTransaction(dataTrans)*/
	op := &api.Operation{}
	op.Schema = `type Buyer {
						id
						name
						age
						products
						time
					}
				type Product{
						p_id
						name
						price
						time
						buy
					}
				type Transaction{
					tr_id
					id
					ip
					device
					products
				}
				id: string @index(term).
				p_id: string @index(term).
				tr_id: string @index(term).
				buyer_id: string @index(term).
				ip: string @index(term).
				device: string @index(term).
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

	dBu, err := json.Marshal(dataBuy)
	if err != nil {
		log.Fatal(err)
	}
	dPr, err := json.Marshal(dataProduct)
	if err != nil {
		log.Fatal(err)
	}
	dTran, err := json.Marshal(dataTrans)
	if err != nil {
		log.Fatal(err)
	}
	mu.SetJson = dBu
	dg.NewTxn().Mutate(ctx, mu)
	mu.SetJson = dPr
	dg.NewTxn().Mutate(ctx, mu)
	mu.SetJson = dTran
	dg.NewTxn().Mutate(ctx, mu)
	CreateRelations(dg, dataTrans)

}
func CreateRelations(dg *dgo.Dgraph, dataTrans []resource.Transaction) {
	ctx := context.Background()
	var dataBuy = resource.AllTransactio(dataTrans, dg)
	dBu, err := json.Marshal(dataBuy)
	if err != nil {
		log.Fatal(err)
	}
	mu := &api.Mutation{
		CommitNow: true,
	}
	mu.SetJson = dBu
	dg.NewTxn().Mutate(ctx, mu)
}
func GetAllPurchaser(dg *dgo.Dgraph) {
	type UidTotal2 struct {
		Id   string `dgraph:"id"`
		Name string `dgraph:"name"`
	}
	type Root2 struct {
		Data []UidTotal2 `dgraph:"data{data}"`
	}
	ctx := context.Background()
	users := (`{
		data(func: has(age)) {
		 id
		 name
		}
	  }
	  `)
	resp2, err := dg.NewTxn().Query(ctx, users)
	if err != nil {
		log.Fatal(err)
	}
	var rUser Root2
	err = json.Unmarshal(resp2.Json, &rUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rUser.Data)

}

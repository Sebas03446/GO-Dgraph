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
	id := "a9fecc93"
	users := fmt.Sprintf(`{
		  data(func: eq(id, "%s")) {
		   uid
		   id
		  }
		}
		`, id)

	//resp, err := dg.NewTxn().QueryWithVars(ctx, q, variables)
	resp, err := dg.NewTxn().Query(ctx, users)
	if err != nil {
		log.Fatal(err)
	}
	/*resp2, err := dg.NewTxn().Query(ctx, users)
	if err != nil {
		log.Fatal(err)
	}*/
	//fmt.Println(string(resp.Json))
	/*type Producto struct {
		Name  string  `dgraph:"name"`
		Price float64 `dgraph:"price"`
	}
	type Comprador struct {
		Name     string     `dgraph:"name"`
		Products []Producto `dgraph:"products"`
	}*/
	type UidTotal struct {
		Id  string `dgraph:"id"`
		Uid string `dgraph:"uid"`
	}
	type Root struct {
		Data []UidTotal `dgraph:"data{data}"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(r.Data[0].P_id)
	var mapUid = make(map[string]string)
	for i := 0; i < len(r.Data); i++ {
		mapUid[r.Data[i].Id] = r.Data[i].Uid
	}
	fmt.Println(mapUid["a9fecc93"])
	/*err = json.Unmarshal(resp2.Json, &r)
	if err != nil {
		log.Fatal(err)
	}*/
	/*var mapUser = make(map[string]string)
	for i := 0; i < len(r.Data); i++ {
		mapUser[r.Data[i].Id] = r.Data[i].Uid
	}
	var dataTrans = resource.TransformTransaction()
	var dataPredicate = resource.GetTransaction(dataTrans, mapUid, mapUser)
	dTr, err := json.Marshal(dataPredicate)
	if err != nil {
		log.Fatal(err)
	}
	mu.SetJson = dTr
	dg.NewTxn().Mutate(ctx, mu)*/
	/*out, _ := json.MarshalIndent(r, "", "\t")
	fmt.Printf("%s\n", out)
	fmt.Println(len(mapUid))*/

}
func CreateRelation(dg *dgo.Dgraph) {
	ctx := context.Background()
	var dataBuy = resource.OneTransactio()
	dBu, err := json.Marshal(dataBuy)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(dBu))
	mu := &api.Mutation{
		CommitNow: true,
	}
	mu.SetJson = dBu
	dg.NewTxn().Mutate(ctx, mu)

}
func CreateRelations(dg *dgo.Dgraph) {
	ctx := context.Background()
	var dataTrans = resource.TransformTransaction()
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

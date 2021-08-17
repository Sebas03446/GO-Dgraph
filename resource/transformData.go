package resource

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgraph-io/dgo/v200"
)

type TransactionProduct struct {
	Uid string `json:"uid"`
}
type TransactionPredicate struct {
	Id         string               `json:"tr_id"`
	Ip         string               `json:"ip"`
	Device     string               `json:"device"`
	TimeNow    int64                `json:"time"`
	BuyerId    string               `json:"uid"`
	ProductIds []TransactionProduct `json:"products"`
}
type UidTotal2 struct {
	Id  string `dgraph:"id"`
	Uid string `dgraph:"uid"`
}
type Root2 struct {
	Data []UidTotal2 `dgraph:"data{data}"`
}
type UidTotal struct {
	P_id string `dgraph:"p_id"`
	Uid  string `dgraph:"uid"`
}
type Root struct {
	Data []UidTotal `dgraph:"data{data}"`
}

func AllTransactio(transaction []Transaction, dg *dgo.Dgraph) []TransactionPredicate {
	ctx := context.Background()
	var allT []TransactionPredicate
	for i := 0; i < len(transaction); i++ {
		id := transaction[i].BuyerId
		users := fmt.Sprintf(`{
		  data(func: eq(id, "%s")) {
		   uid
		   id
		  }
		}
		`, id)
		resp2, err := dg.NewTxn().Query(ctx, users)
		if err != nil {
			log.Fatal(err)
		}
		var rUser Root2
		err = json.Unmarshal(resp2.Json, &rUser)
		if err != nil {
			log.Fatal(err)
		}
		var mapUid2 = make(map[string]string)
		for k := 0; k < len(rUser.Data); k++ {
			mapUid2[rUser.Data[k].Id] = rUser.Data[k].Uid
		}
		var transactionUn TransactionPredicate
		uid := mapUid2[id]
		transactionUn.Id = transaction[i].Id
		transactionUn.Ip = transaction[i].Ip
		transactionUn.Device = transaction[i].Device
		transactionUn.TimeNow = time.Now().Unix()
		transactionUn.BuyerId = uid
		for j := 0; j < len(transaction[i].ProductIds); j++ {
			idProduct := transaction[i].ProductIds[j]
			q := fmt.Sprintf(`{
				data(func: eq(p_id, "%s")) {
				 uid
				 p_id
				}
			  }
			  `, idProduct)
			resp, err := dg.NewTxn().Query(ctx, q)
			if err != nil {
				log.Fatal(err)
			}
			var r Root
			err = json.Unmarshal(resp.Json, &r)
			if err != nil {
				log.Fatal(err)
			}
			var mapUid = make(map[string]string)
			for p := 0; p < len(r.Data); p++ {
				mapUid[r.Data[p].P_id] = r.Data[p].Uid
			}
			var transactionProd TransactionProduct
			transactionProd.Uid = mapUid[idProduct]
			transactionUn.ProductIds = append(transactionUn.ProductIds, transactionProd)
		}
		fmt.Println(transactionUn.ProductIds)
		allT = append(allT, transactionUn)
	}
	return allT
}

func GetProduct() []Product {
	cont, err := http.Get("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/products")
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(cont.Body)
	reader.LazyQuotes = true
	var product []Product
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			continue
		}
		line2 := strings.Split(line[0], "'")
		p, _ := strconv.ParseFloat(line2[2], 64)
		product = append(product, Product{
			Id:    line2[0],
			Name:  line2[1],
			Price: p,
		})
	}
	fmt.Println("Cantidad Productos")
	fmt.Println(len(product))
	return product
}
func GetBuyer() []Buyer {
	cont, err := http.Get("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/buyers")
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(cont.Body)
	var f []Buyer
	json.Unmarshal(body, &f)
	for i := 0; i < len(f); i++ {
		f[i].TimeNow = time.Now().Unix()
	}
	fmt.Println("Cantidad de compradores")
	fmt.Println(len(f))
	return f
}
func TransformTransaction() []Transaction {
	cont, err := http.Get("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/transactions")
	if err != nil {
		log.Fatal(err)
	}
	defer cont.Body.Close()
	bytesCont, _ := io.ReadAll(cont.Body)
	count := 0
	maxPos := 0
	var transaction []Transaction
	bytesCont = append(bytesCont, 35)
	for i := 0; i < len(bytesCont); i++ {
		if bytesCont[i] != 0 && count < 2 {
			count = 0
		} else if bytesCont[i] != 0 && count == 2 {
			count = 0
			var t Transaction
			countPos := 0
			minPos := maxPos
			for maxPos < i-1 {
				if bytesCont[maxPos] != 0 {
					maxPos++
					continue
				} else {
					switch countPos {
					case 0:
						t.Id = string(bytesCont[minPos:maxPos])
						maxPos++
						minPos = maxPos
						countPos++
					case 1:
						t.BuyerId = string(bytesCont[minPos:maxPos])
						maxPos++
						minPos = maxPos
						countPos++
					case 2:
						t.Ip = string(bytesCont[minPos:maxPos])
						maxPos++
						minPos = maxPos
						countPos++
					case 3:
						t.Device = string(bytesCont[minPos:maxPos])
						maxPos++
						minPos = maxPos
						countPos++
					case 4:
						prod := string(bytesCont[minPos+1 : maxPos-1])
						t.ProductIds = strings.Split(prod, ",")
						t.TimeNow = time.Now().Unix()
						maxPos += 2
					}
				}

			}
			transaction = append(transaction, t)
		} else {
			count++
		}

	}
	return transaction
}

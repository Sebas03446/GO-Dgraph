package resource

import (
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

func OneTransactio() TransactionPredicate {
	var newT TransactionPredicate
	newT.BuyerId = "0x439b1"
	var t TransactionProduct
	t.Uid = "0x43c84"
	var p TransactionProduct
	p.Uid = "0x43cb9"
	newT.Id = "#000061185900"
	newT.Ip = "118.42.53.99"
	newT.Device = "android"
	newT.TimeNow = time.Now().Unix()
	newT.ProductIds = append(newT.ProductIds, t)
	newT.ProductIds = append(newT.ProductIds, p)
	return newT
}
func AllTransactio(transaction []Transaction, mapUid map[string]string, mapUser map[string]string) []TransactionPredicate {
	var allT []TransactionPredicate
	for i := 0; i < len(transaction); i++ {
		//REIVAR MAP DE USUARIOS NO SE ESTAN ENVIANDO TODOS LOS USUARIOS, REVISAR EL QUERY (USAR EXPAND ALL ())
		var transactionUn TransactionPredicate
		uid := mapUser[transaction[i].BuyerId]
		//fmt.Println(transaction[i].BuyerId)
		fmt.Println(mapUser[transaction[i].BuyerId])
		transactionUn.Id = transaction[i].Id
		transactionUn.Ip = transaction[i].Ip
		transactionUn.Device = transaction[i].Device
		transactionUn.TimeNow = time.Now().Unix()
		transactionUn.BuyerId = uid
		for j := 0; j < len(transaction[i].ProductIds); j++ {
			var transactionProd TransactionProduct
			//fmt.Println(transaction[i].ProductIds[j])
			transactionProd.Uid = mapUid[transaction[i].ProductIds[j]]
			transactionUn.ProductIds = append(transactionUn.ProductIds, transactionProd)
		}
		//fmt.Println(transactionUn.ProductIds)
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
	//productJson, _ := json.Marshal(product)
	//fmt.Println("lista" + string(productJson))
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

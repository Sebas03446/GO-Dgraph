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

func GetProduct() {
	cont, err := http.Get("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/products")
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(cont.Body)
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
	productJson, _ := json.Marshal(product)
	fmt.Println("lista" + string(productJson))
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

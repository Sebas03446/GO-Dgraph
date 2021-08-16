package resource

type Transaction struct {
	Id         string   `json:"tr_id"`
	BuyerId    string   `json:"id"`
	Ip         string   `json:"ip"`
	Device     string   `json:"device"`
	ProductIds []string `json:"products_id"`
	TimeNow    int64    `json:"time"`
}

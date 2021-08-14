package resource

type Transaction struct {
	Id         string
	BuyerId    string
	Ip         float64
	Device     string
	ProductIds []string
	TimeNow    int64
}

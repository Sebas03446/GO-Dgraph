package resource

type Product struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	TimeNow int64   `json:"time"`
}

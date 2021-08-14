package resource

type Buyer struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	TimeNow int64  `json:"time"`
}

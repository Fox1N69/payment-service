package entity

type User struct {
	ID           int64         `json:"id"`
	Balance      int64         `json:"balance"`
	Transactions []Transaction `json:"transactions"`
}

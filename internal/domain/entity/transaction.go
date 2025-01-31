package entity

import "time"

type Transaction struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Amount      int64     `json:"amount"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

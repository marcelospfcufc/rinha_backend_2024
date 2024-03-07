package entity

import "time"

type Transaction struct {
	Id
	Value       int64
	Operation   string
	Description string
	CreatedAt   time.Time
}

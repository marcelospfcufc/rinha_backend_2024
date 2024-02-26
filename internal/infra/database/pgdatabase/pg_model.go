package pgdatabase

import "time"

type ClientModel struct {
	Id      uint
	Name    string
	Balance int64
	Credit  int64
}

type TransactionModel struct {
	Id          uint
	Value       int64
	Description string
	Operation   string
	Client_Id   uint
	CreatedAt   time.Time
}

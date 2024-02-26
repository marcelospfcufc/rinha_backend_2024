package entity

type Client struct {
	Id
	Name           string
	Credit         int64
	Transactions   []Transaction
	CurrentBalance int64
}

package entity

type Client struct {
	Id
	Name           string
	Credit         int64
	CurrentBalance int64
	Transactions   []Transaction
}

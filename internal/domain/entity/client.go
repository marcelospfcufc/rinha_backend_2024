package entity

type Client struct {
	Id
	Name         string
	Credit       int64
	Transactions []Transaction
}

func (user Client) Balance() int64 {

	var balance int64 = 0

	for _, transaction := range user.Transactions {
		value := transaction.Value
		if transaction.Operation == "d" {
			value = value * -1
		}
		balance += value
	}

	return balance
}

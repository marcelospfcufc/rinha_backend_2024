package entity

import "github.com/marcelospfcufc/rinha_backend_2024/internal/domain"

type Client struct {
	Id
	Name           string
	Credit         int64
	CurrentBalance int64
	Transactions   []Transaction
}

func (cli Client) calculateNewBalance(
	clientCredit int64,
	currentBalance int64,
	transactionValue int64,
	transactionOperation string,
) (clientNewCurrentBalance int64, err error) {
	newBalanceValue := currentBalance
	if transactionOperation == "d" {
		newBalanceValue -= transactionValue

		if newBalanceValue < clientCredit*-1 {
			return -1, domain.ErrClientWithoutBalance
		}
	} else {
		newBalanceValue += transactionValue
	}

	return newBalanceValue, nil
}

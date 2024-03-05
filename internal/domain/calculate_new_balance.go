package domain

func CalculateNewBalance(
	clientCredit int64,
	currentBalance int64,
	transactionValue int64,
	transactionOperation string,
) (clientNewCurrentBalance int64, err error) {
	newBalanceValue := currentBalance
	if transactionOperation == "d" {
		newBalanceValue -= transactionValue

		if newBalanceValue < clientCredit*-1 {
			return -1, ErrClientWithoutBalance
		}
	} else {
		newBalanceValue += transactionValue
	}

	return newBalanceValue, nil
}

package service

import (
	"context"

	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/repository"
)

type InputData struct {
	Value       int64
	Operation   string
	Description string
	ClientId    entity.Id
}

type OutputData struct {
	Credit  int64
	Balance int64
}

type AddTransactionService struct {
	repo repository.ClientRepository
}

func NewAddTransactionService(
	repo repository.ClientRepository,
) *AddTransactionService {
	service := AddTransactionService{
		repo: repo,
	}

	return &service
}

func calculateNewBalance(
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

func (service *AddTransactionService) Execute(ctx context.Context, inputData InputData) (output OutputData, err error) {

	clientCredit, currentNewBalance, err := service.repo.AddTransaction(
		ctx,
		inputData.ClientId,
		entity.Transaction{
			Value:       inputData.Value,
			Operation:   inputData.Operation,
			Description: inputData.Description,
		},
		calculateNewBalance,
	)

	if err != nil {
		return OutputData{}, err
	}

	output = OutputData{
		Credit:  clientCredit,
		Balance: currentNewBalance,
	}

	return
}

package service

import (
	"context"
	"time"

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
	clientRepository      repository.ClientRepository
	transactionRepository repository.TransactionRepository
}

func NewAddTransactionService(
	clientRepository repository.ClientRepository,
	transactionRepository repository.TransactionRepository,
) *AddTransactionService {
	service := AddTransactionService{
		clientRepository:      clientRepository,
		transactionRepository: transactionRepository,
	}

	return &service
}

func (service *AddTransactionService) Execute(inputData InputData) (output OutputData, err error) {

	/*
		var clientFound entity.Client
		clientFound, err = service.clientRepository.GetById(ctx, inputData.ClientId)

		if err != nil {
			return
		}

		clientBalance, err := service.transactionRepository.CalculateBalanceByClient(ctx, inputData.ClientId)
		if err != nil {
			return
		}

		balanceBeforeOperation := clientBalance
		var balanceAfterOperation int64 = 0

		value := inputData.Value

		if inputData.Operation == "d" {
			value = value * -1
		}

		balanceAfterOperation = balanceBeforeOperation + value

		if balanceAfterOperation+clientFound.Credit < 0 {
			err = domain.ErrClientWithoutBalance
			return
		}
	*/

	ctx := context.Background()

	utcTime := time.Now().UTC()

	var repoOutput repository.CreateTransactionOutputData

	repoOutput, err = service.transactionRepository.Create(
		ctx,
		inputData.ClientId,
		entity.Transaction{
			Value:       inputData.Value,
			Operation:   inputData.Operation,
			Description: inputData.Description,
			CreatedAt:   utcTime,
		},
	)

	if err != nil {
		return
	}

	return OutputData{
		Credit:  repoOutput.ClientCredit,
		Balance: repoOutput.CurrentBalance,
	}, nil
}

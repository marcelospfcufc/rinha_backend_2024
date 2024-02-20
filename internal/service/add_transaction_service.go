package service

import (
	"time"

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

	var clientFound entity.Client
	clientFound, err = service.clientRepository.GetById(inputData.ClientId)

	if err != nil {
		return
	}

	balanceBeforeOperation := clientFound.Balance()
	var balanceAfterOperation int64 = 0

	if inputData.Operation == "d" {
		balanceAfterOperation = balanceBeforeOperation - inputData.Value

		if balanceAfterOperation+clientFound.Credit < 0 {
			err = domain.ErrClientWithoutBalance
			return
		}
	} else {
		balanceAfterOperation = balanceBeforeOperation + inputData.Value
	}

	utcTime := time.Now().UTC()

	_, err = service.transactionRepository.Create(
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
		Credit:  clientFound.Credit,
		Balance: balanceAfterOperation,
	}, nil
}

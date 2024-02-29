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
	repo repository.Repository
}

func NewAddTransactionService(
	repo *repository.Repository,
) *AddTransactionService {
	service := AddTransactionService{
		repo: *repo,
	}

	return &service
}

func (service *AddTransactionService) Execute(ctx context.Context, inputData InputData) (output OutputData, err error) {
	var newBalanceValue int64

	domainClient, err := service.repo.GetSimplifiedClientById(ctx, inputData.ClientId)
	if err != nil {
		return OutputData{}, domain.ErrClientNotFound
	}

	newBalanceValue = domainClient.CurrentBalance
	if inputData.Operation == "d" {
		newBalanceValue -= inputData.Value

		if newBalanceValue < domainClient.Credit*-1 {
			return OutputData{}, domain.ErrClientWithoutBalance
		}
	} else {
		newBalanceValue += inputData.Value
	}

	err = service.repo.UpdateClientBalance(
		ctx,
		inputData.ClientId,
		newBalanceValue,
	)

	if err != nil {
		return OutputData{}, err
	}

	err = service.repo.AddTransaction(
		ctx,
		inputData.ClientId,
		entity.Transaction{
			Value:       inputData.Value,
			Operation:   inputData.Operation,
			Description: inputData.Description,
		},
	)

	if err != nil {
		return OutputData{}, err
	}

	return OutputData{
		Credit:  domainClient.Credit,
		Balance: newBalanceValue,
	}, nil
}

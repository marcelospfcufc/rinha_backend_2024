package service

import (
	"context"
	"errors"

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

func (service *AddTransactionService) Execute(ctx context.Context, inputData InputData) (output OutputData, err error) {
	err = errors.New("not implemented yet")
	output = OutputData{}
	return
}

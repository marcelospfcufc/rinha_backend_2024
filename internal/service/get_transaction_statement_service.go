package service

import (
	"context"
	"time"

	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/repository"
)

type GetTransactionStatementInputData struct {
	ClientId entity.Id
}

type GetTransactionStatementTransactionData struct {
	Value       int64
	Operation   string
	Description string
	CreatedAt   time.Time
}

type GetTransactionStatementOutputData struct {
	Balance      int64
	Credit       int64
	Transactions []GetTransactionStatementTransactionData
}

type GetTransactionStatementService struct {
	repo repository.Repository
}

func NewGetTransactionStatementService(
	repo *repository.Repository,
) *GetTransactionStatementService {
	service := GetTransactionStatementService{
		repo: *repo,
	}

	return &service
}

func (service *GetTransactionStatementService) Execute(
	ctx context.Context,
	inputData GetTransactionStatementInputData,
) (output GetTransactionStatementOutputData, err error) {

	var outputData GetTransactionStatementOutputData = GetTransactionStatementOutputData{}

	client, err := service.repo.GetSimplifiedClientById(ctx, inputData.ClientId)
	if err != nil {
		return outputData, domain.ErrClientNotFound
	}

	transactions, err := service.repo.GetTransactions(
		ctx,
		inputData.ClientId,
		10,
		repository.Desc,
	)

	if err != nil {
		return outputData, err
	}

	transactionsToReturn := make([]GetTransactionStatementTransactionData, len(transactions))

	for idx, transaction := range transactions {
		transactionsToReturn[idx] = GetTransactionStatementTransactionData{
			Value:       transaction.Value,
			CreatedAt:   transaction.CreatedAt,
			Operation:   transaction.Operation,
			Description: transaction.Description,
		}
	}

	outputData.Balance = client.CurrentBalance
	outputData.Credit = client.Credit
	outputData.Transactions = transactionsToReturn

	return outputData, err
}

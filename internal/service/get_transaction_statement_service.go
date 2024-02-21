package service

import (
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
	clientRepository      repository.ClientRepository
	transactionRepository repository.TransactionRepository
}

func NewGetTransactionStatementService(
	clientRepository repository.ClientRepository,
	transactionRepository repository.TransactionRepository,
) *GetTransactionStatementService {
	service := GetTransactionStatementService{
		clientRepository:      clientRepository,
		transactionRepository: transactionRepository,
	}

	return &service
}

func (service *GetTransactionStatementService) Execute(
	inputData GetTransactionStatementInputData,
) (output GetTransactionStatementOutputData, err error) {

	var outputData GetTransactionStatementOutputData
	hasClient := service.clientRepository.HasClientById(inputData.ClientId)

	if !hasClient {
		return outputData, domain.ErrClientNotFound
	}

	transactions, err := service.transactionRepository.GetAllByUser(inputData.ClientId, 10, repository.Desc)

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

	summaryBalance, err := service.transactionRepository.SummaryBalanceByClient(inputData.ClientId)
	if err != nil {
		return outputData, err
	}

	outputData.Balance = summaryBalance.Balance
	outputData.Credit = summaryBalance.Credit
	outputData.Transactions = transactionsToReturn

	return outputData, err
}

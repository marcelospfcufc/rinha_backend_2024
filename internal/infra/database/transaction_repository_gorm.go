package database

import (
	"context"

	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/repository"
	"gorm.io/gorm"
)

type TransactionRepositoryGorm struct {
	dbConnection *gorm.DB
}

func NewTransactionRepositoryGorm(dbConnection *gorm.DB) *TransactionRepositoryGorm {
	return &TransactionRepositoryGorm{dbConnection: dbConnection}
}

func (repo *TransactionRepositoryGorm) Create(ctx context.Context, clientId entity.Id, transaction entity.Transaction) (entity.Transaction, error) {
	transactionModel := Transaction{
		Operation:   transaction.Operation,
		Value:       transaction.Value,
		Description: transaction.Description,
		CreatedAt:   transaction.CreatedAt,
		ClientID:    clientId,
	}

	tx := repo.dbConnection.Begin()

	result := tx.WithContext(ctx).Create(&transactionModel)
	if result.Error != nil {
		tx.Rollback()
		return entity.Transaction{}, result.Error
	}

	tx.Commit()

	return entity.Transaction{
		Id:          transactionModel.ID,
		Value:       transaction.Value,
		Operation:   transaction.Operation,
		Description: transaction.Description,
		CreatedAt:   transaction.CreatedAt,
	}, result.Error
}

func (repo *TransactionRepositoryGorm) GetAllByUser(
	ctx context.Context,
	clientId entity.Id,
	limit int,
	orderBy repository.OrderBy,
) (
	[]entity.Transaction,
	error,
) {

	var transactions []Transaction

	result := repo.dbConnection.WithContext(ctx).Where("client_id = ?", clientId).Order("created_at " + orderBy.String()).Limit(limit).Find(&transactions)

	if result.Error != nil {
		return []entity.Transaction{}, result.Error
	}

	var transactionsToReturn []entity.Transaction
	for _, transaction := range transactions {
		transactionsToReturn = append(
			transactionsToReturn,
			entity.Transaction{
				Id:          transaction.ID,
				Value:       transaction.Value,
				Operation:   transaction.Operation,
				Description: transaction.Description,
				CreatedAt:   transaction.CreatedAt,
			},
		)
	}

	return transactionsToReturn, nil
}

func (repo *TransactionRepositoryGorm) SummaryBalanceByClient(ctx context.Context, clientId entity.Id) (repository.SummaryBalanceRepositoryData, error) {
	var clientFound Client
	var summaryBalanceData repository.SummaryBalanceRepositoryData
	var result *gorm.DB

	result = repo.dbConnection.WithContext(ctx).First(&clientFound, clientId)

	if result.Error != nil {
		return repository.SummaryBalanceRepositoryData{}, domain.ErrClientNotFound
	}

	var clientBalance int64
	result = repo.dbConnection.WithContext(ctx).Raw("SELECT COALESCE(SUM(CASE WHEN operation = 'c' THEN value ELSE -value END),0) FROM transactions WHERE client_id =?", clientId).Scan(&clientBalance)
	if result.Error != nil {
		return summaryBalanceData, result.Error
	}

	summaryBalanceData.Balance = clientBalance
	summaryBalanceData.Credit = clientFound.Credit

	return summaryBalanceData, nil

}

func (repo *TransactionRepositoryGorm) CalculateBalanceByClient(ctx context.Context, clientId entity.Id) (int64, error) {

	var clientBalance int64
	result := repo.dbConnection.WithContext(ctx).Raw("SELECT COALESCE(SUM(CASE WHEN operation = 'c' THEN value ELSE -value END),0) FROM transactions WHERE client_id =?", clientId).Scan(&clientBalance)
	if result.Error != nil {
		return clientBalance, result.Error
	}

	return clientBalance, nil
}

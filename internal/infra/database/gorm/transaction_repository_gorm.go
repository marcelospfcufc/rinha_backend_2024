package database

import (
	"errors"

	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
	"gorm.io/gorm"
)

type TransactionRepositoryGorm struct {
	dbConnection *gorm.DB
}

func NewTransactionRepositoryGorm(dbConnection *gorm.DB) *TransactionRepositoryGorm {
	return &TransactionRepositoryGorm{dbConnection: dbConnection}
}

func (repo *TransactionRepositoryGorm) Create(clientId entity.Id, transaction entity.Transaction) (entity.Transaction, error) {
	transactionModel := Transaction{
		Operation:   transaction.Operation,
		Value:       transaction.Value,
		Description: transaction.Description,
		CreatedAt:   transaction.CreatedAt,
		ClientID:    clientId,
	}

	result := repo.dbConnection.Create(&transactionModel)

	if result.Error != nil {
		return entity.Transaction{}, result.Error
	}

	return entity.Transaction{
		Id:          transactionModel.ID,
		Value:       transaction.Value,
		Operation:   transaction.Operation,
		Description: transaction.Description,
		CreatedAt:   transaction.CreatedAt,
	}, result.Error
}

func (repo *TransactionRepositoryGorm) GetAllByUser(clientId entity.Id) ([]entity.Transaction, error) {
	return []entity.Transaction{}, errors.New("not implemented")
}

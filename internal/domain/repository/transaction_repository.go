package repository

import (
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
)

type SummaryBalanceRepositoryData struct {
	Balance int64
	Credit  int64
}

type TransactionRepository interface {
	Create(clientId entity.Id, transaction entity.Transaction) (entity.Transaction, error)
	GetAllByUser(clientId entity.Id, limit int, orderBy OrderBy) ([]entity.Transaction, error)
	SummaryBalanceByClient(clientId entity.Id) (SummaryBalanceRepositoryData, error)
}

package repository

import (
	"context"

	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
)

type SummaryBalanceRepositoryData struct {
	Balance int64
	Credit  int64
}

type TransactionRepository interface {
	Create(ctx context.Context, clientId entity.Id, transaction entity.Transaction) (entity.Transaction, error)
	GetAllByUser(ctx context.Context, clientId entity.Id, limit int, orderBy OrderBy) ([]entity.Transaction, error)
	SummaryBalanceByClient(ctx context.Context, clientId entity.Id) (SummaryBalanceRepositoryData, error)
	CalculateBalanceByClient(ctx context.Context, clientId entity.Id) (int64, error)
}

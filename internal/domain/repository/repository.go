package repository

import (
	"context"

	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
)

type Repository interface {
	HasClientById(ctx context.Context, clientId entity.Id) bool
	Create(ctx context.Context, client entity.Client) (entity.Client, error)
	UpdateClientBalance(ctx context.Context, clientId entity.Id, newBalance int64) error
	AddTransaction(ctx context.Context, clientId entity.Id, transaction entity.Transaction) error
	GetTransactions(ctx context.Context, clientId entity.Id, limit int, orderBy OrderBy) ([]entity.Transaction, error)
	GetById(ctx context.Context, clientId entity.Id) (entity.Client, error)
	GetSimplifiedClientById(ctx context.Context, clientId entity.Id) (entity.Client, error)
	CalculateBalanceByClient(ctx context.Context, clientId entity.Id) (int64, error)
}

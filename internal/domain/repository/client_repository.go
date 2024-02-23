package repository

import (
	"context"

	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
)

type ClientRepository interface {
	HasClientById(ctx context.Context, clientId entity.Id) bool
	Create(client entity.Client) (entity.Client, error)
	GetById(ctx context.Context, clientId entity.Id) (entity.Client, error)
	GetTransactionsById(ctx context.Context, clientId entity.Id) ([]entity.Transaction, error)
}

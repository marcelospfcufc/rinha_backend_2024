package repository

import (
	"context"

	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
)

type ClientRepository interface {
	Create(
		ctx context.Context,
		client entity.Client,
	) (entity.Client, error)

	AddTransaction(
		ctx context.Context,
		clientId entity.Id,
		transaction entity.Transaction,
		calculateNewBalance func(
			clientCredit int64,
			currentBalance int64,
			transactionValue int64,
			transactionOperation string,
		) (clientNewCurrentBalance int64, err error),
	) (clientCredit int64, clientNewCurrentBalance int64, err error)

	GetTransactions(
		ctx context.Context,
		clientId entity.Id,
		limit int,
		orderBy OrderBy,
	) ([]entity.Transaction, error)

	GetClientById(
		ctx context.Context,
		clientId entity.Id,
	) (entity.Client, error)
}

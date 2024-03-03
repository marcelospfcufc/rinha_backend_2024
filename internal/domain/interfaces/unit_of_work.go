package interfaces

import (
	"context"

	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/repository"
)

type UnitOfWork interface {
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	RollBack(ctx context.Context) error
	GetRepository() repository.ClientRepository
}

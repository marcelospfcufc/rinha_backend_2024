package pgdatabase

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/repository"
)

type PgUnitOfWork struct {
	db   *pgxpool.Pool
	dbTx *pgx.Tx
}

func NewPgUnitOfWork(
	db *pgxpool.Pool,
) (*PgUnitOfWork, error) {

	return &PgUnitOfWork{
		db:   db,
		dbTx: nil,
	}, nil
}

func (unit *PgUnitOfWork) Begin(ctx context.Context) error {

	tx, err := unit.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return domain.ErrInternalDatabaseError
	}
	unit.dbTx = &tx
	return err
}

func (unit *PgUnitOfWork) Commit(ctx context.Context) error {
	return (*unit.dbTx).Commit(ctx)
}

func (unit *PgUnitOfWork) RollBack(ctx context.Context) error {
	return (*unit.dbTx).Rollback(ctx)
}

func (unit *PgUnitOfWork) GetRepository() repository.ClientRepository {
	repo := NewPgRepository(unit.db, unit.dbTx)
	return repo
}

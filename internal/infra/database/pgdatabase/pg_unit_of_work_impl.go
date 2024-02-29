package pgdatabase

import (
	"context"
	"database/sql"
	"log"

	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/repository"
)

type PgUnitOfWork struct {
	db          *sql.DB
	transaction *sql.Tx
}

func NewPgUnitOfWork(
	db *sql.DB,
) (*PgUnitOfWork, error) {

	return &PgUnitOfWork{
		db:          db,
		transaction: nil,
	}, nil
}

func (unit *PgUnitOfWork) Begin(ctx context.Context) error {
	tx, err := unit.db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	unit.transaction = tx

	return err
}

func (unit *PgUnitOfWork) Commit() error {
	return unit.transaction.Commit()
}

func (unit *PgUnitOfWork) RollBack() error {
	return unit.transaction.Rollback()
}

func (unit *PgUnitOfWork) GetRepository() *repository.Repository {
	repo := repository.Repository(NewPgRepository(unit.db, unit.transaction))
	return &repo
}

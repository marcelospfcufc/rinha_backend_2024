package pgdatabase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
)

type PgRepository struct {
	conn *sql.DB
}

func NewPgRepository(dbConnection *sql.DB) *PgRepository {
	return &PgRepository{conn: dbConnection}
}

func (pg *PgRepository) HasClientById(ctx context.Context, clientId entity.Id) bool {

	var clientIdFound uint

	query := "SELECT id FROM Clients WHERE id=$1"

	err := pg.conn.QueryRowContext(ctx, query, clientId).Scan(&clientIdFound)

	return err == nil
}

func (pg *PgRepository) Create(client entity.Client) (entity.Client, error) {
	/* clientModel := Client{
		Name:           client.Name,
		Credit:         client.Credit,
		CurrentBalance: client.CurrentBalance,
	}

	result := repo.dbConnection.Create(&clientModel)

	if result.Error != nil {
		return entity.Client{}, result.Error
	}

	return entity.Client{
		Id:             clientModel.ID,
		Name:           clientModel.Name,
		Credit:         clientModel.Credit,
		CurrentBalance: clientModel.CurrentBalance,
	}, result.Error */

	return entity.Client{}, errors.New("not implemented yet")
}

func (pg *PgRepository) GetById(ctx context.Context, clientId entity.Id) (entity.Client, error) {
	return entity.Client{}, errors.New("not implemented yet")
}

func (pg *PgRepository) GetTransactionsById(ctx context.Context, clientId entity.Id) ([]entity.Transaction, error) {
	return []entity.Transaction{}, errors.New("not implemented yet")

}

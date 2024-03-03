package pgdatabase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/repository"
)

type PgRepository struct {
	conn *pgxpool.Pool
	tx   pgx.Tx
}

func NewPgRepository(
	dbConnection *pgxpool.Pool,
	tx pgx.Tx,

) *PgRepository {
	return &PgRepository{
		conn: dbConnection,
		tx:   tx,
	}
}

func (pg *PgRepository) Create(
	ctx context.Context,
	client entity.Client,
) (entity.Client, error) {
	return entity.Client{}, errors.New("not implemented yet")
}

func (pg *PgRepository) AddTransaction(
	ctx context.Context,
	clientId entity.Id,
	transaction entity.Transaction,
	calculateNewBalance func(
		clientCredit int64,
		currentBalance int64,
		transactionValue int64,
		transactionOperation string,
	) (clientNewCurrentBalance int64, err error),
) (clientCredit int64, clientNewCurrentBalance int64, err error) {

	var dbClient entity.Client
	var row pgx.Row

	if pg.tx == nil {
		log.Fatal("dbTx is nil in AddTransacion")
	}

	queryToExecute := fmt.Sprintf(
		`
			SELECT 
			c.id, c.name, c.credit, c.balance			
			FROM Clients c 			
			WHERE c.id=%d 		
			FOR UPDATE	
		`,
		clientId,
	)

	row = pg.tx.QueryRow(ctx, queryToExecute)

	err = row.Scan(
		&dbClient.Id,
		&dbClient.Name,
		&dbClient.Credit,
		&dbClient.CurrentBalance,
	)

	if err != nil {
		return -1, -1, domain.ErrClientNotFound
	}

	newBalance, err := calculateNewBalance(
		dbClient.Credit,
		dbClient.CurrentBalance,
		transaction.Value,
		transaction.Operation,
	)

	if err != nil {
		return -1, -1, err
	}

	queryToExecute = fmt.Sprintf(
		`
			UPDATE clients 
			SET balance = %d
			WHERE id=%d			
    	`,
		newBalance,
		clientId,
	)

	_, err = pg.tx.Exec(ctx, queryToExecute)

	if err != nil {
		return -1, -1, err
	}

	queryToExecute = fmt.Sprintf(
		`
			INSERT INTO Transactions (value, description, operation, client_id)
			VALUES (%d, '%s', '%s', %d)
    	`,
		transaction.Value,
		transaction.Description,
		transaction.Operation,
		clientId,
	)

	_, err = pg.tx.Exec(ctx, queryToExecute)

	if err != nil {
		return -1, -1, err
	}

	return dbClient.Credit, newBalance, err
}

func (pg *PgRepository) GetTransactions(
	ctx context.Context,
	clientId entity.Id,
	limit int,
	orderBy repository.OrderBy,
) ([]entity.Transaction, error) {

	queryToExecute := fmt.Sprintf(
		`
			SELECT 
				id, 
				value, 
				description, 
				operation, 			
				created_at			
			FROM Transactions 
			WHERE 
				client_id=%d 
			ORDER BY 
				created_at %s, id %s 
			LIMIT %d 
		`,
		clientId,
		orderBy.String(),
		orderBy.String(),
		limit,
	)

	var rows pgx.Rows
	var err error

	if pg.tx != nil {
		rows, err = pg.tx.Query(ctx, queryToExecute)
	} else {
		rows, err = pg.conn.Query(ctx, queryToExecute)
	}

	if err != nil {
		return []entity.Transaction{}, err
	}
	defer rows.Close()

	var transactions []entity.Transaction
	for rows.Next() {
		var (
			transactionID int64
			value         int64
			description   string
			operation     string
			createdAt     time.Time
		)
		if err := rows.Scan(&transactionID, &value, &description, &operation, &createdAt); err != nil {
			return []entity.Transaction{}, err
		}

		transactions = append(
			transactions,
			entity.Transaction{
				Id:          entity.Id(transactionID),
				Value:       value,
				Description: description,
				Operation:   operation,
				CreatedAt:   createdAt,
			},
		)
	}

	if err := rows.Err(); err != nil {
		return []entity.Transaction{}, err
	}

	return transactions, nil
}

func (pg *PgRepository) GetClientById(
	ctx context.Context,
	clientId entity.Id,
) (entity.Client, error) {

	queryToExecute := fmt.Sprintf(
		`
			SELECT 
			c.id, c.name, c.credit, c.balance			
			FROM Clients c 			
			WHERE c.id=%d 					
		`,
		clientId,
	)

	var clientToReturn entity.Client = entity.Client{}
	var row pgx.Row
	var err error

	if pg.tx != nil {
		row = pg.tx.QueryRow(ctx, queryToExecute)
	} else {
		row = pg.conn.QueryRow(ctx, queryToExecute)
	}

	err = row.Scan(
		&clientToReturn.Id,
		&clientToReturn.Name,
		&clientToReturn.Credit,
		&clientToReturn.CurrentBalance,
	)

	if err != nil {
		return entity.Client{}, err
	}

	clientToReturn.Transactions = []entity.Transaction{}

	return clientToReturn, err
}

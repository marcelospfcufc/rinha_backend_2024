package pgdatabase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/repository"
)

type PgRepository struct {
	conn *sql.DB
	tx   *sql.Tx
	mut  sync.Mutex
}

func NewPgRepository(
	dbConnection *sql.DB,
	tx *sql.Tx,

) *PgRepository {
	return &PgRepository{
		conn: dbConnection,
		tx:   tx,
	}
}

func (pg *PgRepository) HasClientById(ctx context.Context, clientId entity.Id) bool {

	var clientIdFound uint

	queryToExecute := fmt.Sprintf(
		`
		SELECT 
			id 
		FROM 
			Clients 
		WHERE 
			id=%d
		`,
		clientId,
	)

	var row *sql.Row
	var err error

	if pg.tx != nil {
		row = pg.tx.QueryRowContext(ctx, queryToExecute)
	} else {
		row = pg.conn.QueryRowContext(ctx, queryToExecute)
	}

	row.Scan(&clientIdFound)

	return err == nil
}

func (pg *PgRepository) Create(ctx context.Context, client entity.Client) (entity.Client, error) {
	return entity.Client{}, errors.New("not implemented yet")
}

func (pg *PgRepository) UpdateClientBalance(
	ctx context.Context,
	clientId entity.Id,
	newBalance int64,
) error {

	queryToExecute := fmt.Sprintf(
		`
			UPDATE clients 
			SET balance = %d
			WHERE id=%d			
    	`,
		newBalance,
		clientId,
	)

	var err error

	if pg.tx != nil {
		_, err = pg.tx.ExecContext(ctx, queryToExecute)
	} else {
		_, err = pg.conn.ExecContext(ctx, queryToExecute)
	}

	log.Info("result: ", newBalance)

	return err
}

func (pg *PgRepository) AddTransaction(ctx context.Context, clientId entity.Id, transaction entity.Transaction) error {

	queryToExecute := fmt.Sprintf(
		`
			INSERT INTO Transactions (value, description, operation, client_id)
			VALUES (%d, '%s', '%s', %d)
    	`,
		transaction.Value,
		transaction.Description,
		transaction.Operation,
		clientId,
	)

	var err error

	if pg.tx != nil {
		_, err = pg.tx.ExecContext(ctx, queryToExecute)
	} else {
		_, err = pg.conn.ExecContext(ctx, queryToExecute)
	}

	return err
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

	var rows *sql.Rows
	var err error

	if pg.tx != nil {
		rows, err = pg.tx.QueryContext(ctx, queryToExecute)
	} else {
		rows, err = pg.conn.QueryContext(ctx, queryToExecute)
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

func (pg *PgRepository) GetById(ctx context.Context, clientId entity.Id) (entity.Client, error) {

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
	var row *sql.Row
	var err error

	if pg.tx != nil {
		row = pg.tx.QueryRowContext(
			ctx,
			queryToExecute,
		)
	} else {
		row = pg.conn.QueryRowContext(ctx, queryToExecute)
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

	queryToExecute = fmt.Sprintf(
		`
			SELECT 	t.id, t.value, t.operation, t.description, t.created_at
			FROM
			Transactions t
			Where t.client_id=%d
			ORDER BY created_at ASC, id ASC		
		`,
		clientId,
	)

	var rows *sql.Rows

	if pg.tx != nil {
		rows, err = pg.tx.QueryContext(
			ctx,
			queryToExecute,
		)
	} else {
		rows, err = pg.conn.QueryContext(
			ctx,
			queryToExecute,
		)
	}

	if err != nil {
		return entity.Client{}, err
	}

	defer rows.Close()

	var transactionsToReturn []entity.Transaction = []entity.Transaction{}

	for rows.Next() {

		transaction := entity.Transaction{}
		err = rows.Scan(
			&transaction.Id,
			&transaction.Value,
			&transaction.Operation,
			&transaction.Description,
			&transaction.CreatedAt,
		)
		if err != nil {
			return entity.Client{}, err
		}

		transactionsToReturn = append(transactionsToReturn, transaction)
	}

	clientToReturn.Transactions = transactionsToReturn

	return clientToReturn, err
}

func (pg *PgRepository) GetSimplifiedClientById(ctx context.Context, clientId entity.Id) (entity.Client, error) {

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
	var row *sql.Row
	var err error

	if pg.tx != nil {
		row = pg.tx.QueryRowContext(ctx, queryToExecute)
	} else {
		row = pg.conn.QueryRowContext(ctx, queryToExecute)
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

func (pg *PgRepository) CalculateBalanceByClient(ctx context.Context, clientId entity.Id) (int64, error) {

	queryToExecute := fmt.Sprintf(
		`
			SELECT 
			COALESCE(
				SUM(
					CASE WHEN operation = 'c' THEN value ELSE -value END
				)
				,0
			) AS total
			FROM 
				transactions 
			WHERE 
				client_id=%d
			ORDER BY 
				created_at ASC, id ASC
		`,
		clientId,
	)

	var calculatedBalance int64 = 0

	var row *sql.Row
	var err error

	if pg.tx != nil {
		row = pg.tx.QueryRowContext(ctx, queryToExecute)
	} else {
		row = pg.conn.QueryRowContext(ctx, queryToExecute)
	}

	err = row.Scan(
		&calculatedBalance,
	)

	return calculatedBalance, err
}

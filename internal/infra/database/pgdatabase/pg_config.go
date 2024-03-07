package pgdatabase

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateDatabase(ctx context.Context, dbPool *pgxpool.Pool) error {

	statements := []string{
		"DROP TABLE IF EXISTS Transactions;",
		"DROP TABLE IF EXISTS Clients;",
		`
		CREATE TABLE Clients (
			id BIGSERIAL PRIMARY KEY,
			name VARCHAR(255),
			credit INTEGER,
			balance INTEGER DEFAULT 0
		);
		`,
		`
		CREATE TABLE Transactions (
			id BIGSERIAL PRIMARY KEY,
			value INTEGER,
			description VARCHAR(10),
			operation CHAR(1),
			client_id INT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			FOREIGN KEY (client_id) REFERENCES Clients(id)
		);
		`,
		`
		INSERT INTO Clients (id, name, credit)
		VALUES
			(1,'o barato sai caro', 1000 * 100),
			(2,'zan corp ltda', 800 * 100),
			(3,'les cruders', 10000 * 100),
			(4,'padaria joia de cocaia', 100000 * 100),
			(5,'kid mais', 5000 * 100);
		`,
	}

	for _, statement := range statements {
		_, err := dbPool.Exec(context.Background(), statement)
		if err != nil {
			panic(err)
		}
	}

	log.Info("Tables created and data inserted successfully")

	return nil
}

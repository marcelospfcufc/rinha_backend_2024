package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/infra/database/pgdatabase"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/infra/rest"
)

type HandlerFunc func(c *fiber.Ctx) error

func main() {

	connStr := "postgres://postgres:qpalzm@172.24.0.4/rinha_db?sslmode=disable&timezone=UTC"
	dbPool, err := pgxpool.New(context.Background(), connStr)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	defer dbPool.Close()
	dbPool.Config().MaxConns = 4

	err = pgdatabase.CreateDatabase(context.Background(), dbPool)
	if err != nil {

		log.Fatal(err)
	}

	start := make(chan string, 10)
	finish := make(chan string, 10)

	processCtx := func(ctxs <-chan string) {
		for ctx := range ctxs {
			log.Info("Processing: ", ctx)
			log.Info("Finish: ", <-finish)
		}
	}

	go processCtx(start)

	queueGet := make(chan *fiber.Ctx, 20)
	go func() {
		for elem := range queueGet {
			log.Info("Processing request: ", elem.Method(), " - url:", elem.Path())
		}
	}()

	app := fiber.New()

	app.Post("/clientes/:id/transacoes", rest.PostWrapper(dbPool, start, finish))
	app.Get("/clientes/:id/extrato", rest.GetWrapper(dbPool, queueGet))

	app.Listen(":8081")
}

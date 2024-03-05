package main

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/infra/rest"
)

type HandlerFunc func(c *fiber.Ctx) error

func main() {

	connStr := "postgres://postgres:qpalzm@172.29.0.2/rinha_db?sslmode=disable&timezone=UTC"
	//connStr := os.Getenv("DB_URI")

	if strings.TrimSpace(connStr) == "" {
		log.Fatal("not found environment variable 'DB_URI'")
	}

	dbPool, err := pgxpool.New(context.Background(), connStr)

	if err != nil {
		log.Fatal(err.Error())
	}

	defer dbPool.Close()
	dbPool.Config().MaxConns = 100

	startPost := make(chan string, 10)
	finishPost := make(chan string, 10)

	startGet := make(chan string, 50)
	finishGet := make(chan string, 50)

	processCtx := func(ctxs <-chan string, finishes <-chan string) {
		for ctx := range ctxs {
			log.Info("Processing: ", ctx)
			log.Info("Finish: ", <-finishes)
		}
	}

	go processCtx(startPost, finishPost)
	go processCtx(startGet, finishGet)

	app := fiber.New()

	app.Post("/clientes/:id/transacoes", rest.PostWrapper(dbPool, startPost, finishPost))
	app.Get("/clientes/:id/extrato", rest.GetWrapper(dbPool, startGet, finishGet))

	app.Listen(":8080")
}

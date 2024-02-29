package main

import (
	"database/sql"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/infra/database/pgdatabase"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/infra/rest"
)

type HandlerFunc func(c *fiber.Ctx) error

func main() {

	connStr := "postgres://postgres:qpalzm@172.24.0.4/rinha_db?sslmode=disable&timezone=UTC"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	defer db.Close()

	db.SetMaxOpenConns(200) // Define o número máximo de conexões abertas
	db.SetMaxIdleConns(100)

	err = pgdatabase.CreateDatabase(db)
	if err != nil {
		log.Fatal(err)
	}

	countGet := 0
	countPost := 0

	queuePost := make(chan *fiber.Ctx, 10000)
	start := make(chan string)
	finish := make(chan string)

	processCtx := func(ctxs <-chan *fiber.Ctx) {
		for ctx := range ctxs {
			log.Info("Processando requisição para o caminho: %s\n", ctx.Path())
			countPost++
			start <- strconv.Itoa(countPost)
			log.Info("Finalizando: ", <-finish)
		}
	}

	go processCtx(queuePost)

	queueGet := make(chan *fiber.Ctx, 1)
	go func() {
		for elem := range queueGet {
			log.Info("Processing request: ", elem.Method(), " - url:", elem.Path())
			countGet++
		}
	}()

	app := fiber.New()

	app.Post("/clientes/:id/transacoes", rest.PostWrapper(db, queuePost, start, finish))
	app.Get("/clientes/:id/extrato", rest.GetWrapper(db, queueGet))

	app.Listen(":8081")
}

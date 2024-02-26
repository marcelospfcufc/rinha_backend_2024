package main

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/infra/database/pgdatabase"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/infra/rest"
)

type HandlerFunc func(c *fiber.Ctx) error

func main() {

	connStr := "postgres://postgres:qpalzm@172.24.0.3/rinha_db?sslmode=disable&timezone=UTC"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	defer db.Close()

	err = pgdatabase.CreateDatabase(db)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Post("/clientes/:id/transacoes", rest.PostWrapper(db))
	app.Get("/clientes/:id/extrato", rest.GetWrapper(db))

	app.Listen(":8081")
}

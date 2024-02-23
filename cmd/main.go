package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/infra/database"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/infra/rest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type HandlerFunc func(c *fiber.Ctx) error

func main() {

	//dsn := "rinha:qpalzm@tcp(localhost:3306)/rinha_db?charset=utf8mb4&parseTime=True&loc=UTC"
	dsn := "host=172.24.0.2 user=postgres password=qpalzm dbname=rinha_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxLifetime(time.Minute * 5)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(10)

	defer sqlDB.Close()

	db.AutoMigrate(&database.Client{}, &database.Transaction{})

	initialData := []*database.Client{
		{ID: 1, Name: "o barato sai caro", Credit: 1_000 * 100},
		{ID: 2, Name: "zan corp ltda", Credit: 800 * 100},
		{ID: 3, Name: "les cruders", Credit: 10_000 * 100},
		{ID: 4, Name: "padaria joia de cocaia", Credit: 100_000 * 100},
		{ID: 5, Name: "kid mais", Credit: 5_000 * 100},
	}

	testClient := database.Client{ID: 1}
	result := db.WithContext(context.Background()).First(&testClient)

	if result.Error != nil {
		result = db.Create(initialData)
		fmt.Println("Iniciando os dados da base: ", result.Error)
	}

	app := fiber.New()

	/* app.Use(func(c *fiber.Ctx) error {
		//time.Sleep(time.Millisecond * 100)
		log.Info("Request para: ", c.Path(), " - ", c.Method(), " - ", c.Params)
		return c.Next()
	}) */

	app.Post("/clientes/:id/transacoes", rest.PostWrapper(db))
	app.Get("/clientes/:id/extrato", rest.GetWrapper(db))

	app.Listen(":8081")
}

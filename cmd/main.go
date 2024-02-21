package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	database "github.com/marcelospfcufc/rinha_backend_2024/internal/infra/database/gorm"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/infra/rest"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dsn := "rinha:qpalzm@tcp(localhost:3306)/rinha_db?charset=utf8mb4&parseTime=True&loc=UTC"
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               dsn,
		DefaultStringSize: 256,
	}), &gorm.Config{})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	db.AutoMigrate(&database.Client{}, &database.Transaction{})

	initialData := []*database.Client{
		{ID: 1, Name: "o barato sai caro", Credit: 1_000 * 100},
		{ID: 2, Name: "zan corp ltda", Credit: 800 * 100},
		{ID: 3, Name: "les cruders", Credit: 10_000 * 100},
		{ID: 4, Name: "padaria joia de cocaia", Credit: 100_000 * 100},
		{ID: 5, Name: "kid mais", Credit: 5_000 * 100},
	}

	testClient := database.Client{ID: 1}
	result := db.First(&testClient)

	if result.Error != nil || result.RowsAffected == 0 {
		result = db.Create(initialData)
		fmt.Println("Iniciando os dados da base: ", result.Error)
	}

	app := fiber.New()

	transactionRoutes := rest.NewTransactionRoute(app, db)
	transactionRoutes.AddRoutes()
	app.Listen(":3000")
}

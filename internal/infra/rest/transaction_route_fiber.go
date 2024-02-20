package rest

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/controller"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
	database "github.com/marcelospfcufc/rinha_backend_2024/internal/infra/database/gorm"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/service"
	"gorm.io/gorm"
)

type TransactionRoute struct {
	app *fiber.App
	db  *gorm.DB
}

func NewTransactionRoute(app *fiber.App, db *gorm.DB) *TransactionRoute {
	return &TransactionRoute{
		app: app,
		db:  db,
	}
}

func (transactionRoute *TransactionRoute) AddRoutes() {

	transactionRoute.app.Post("/clientes/:id/transacoes", func(c *fiber.Ctx) error {

		transactionRepository := database.NewTransactionRepositoryGorm(transactionRoute.db)
		clientRepository := database.NewClientRepositoryGorm(transactionRoute.db)
		addTransactionService := service.NewAddTransactionService(
			clientRepository,
			transactionRepository,
		)
		ctrl := controller.NewAddTransactionController(*addTransactionService)

		id, _ := strconv.Atoi(c.Params("id"))
		var body controller.AddTransactionInputDto
		c.BodyParser(&body)

		outputController, err := ctrl.AddTransaction(controller.AddTransactionInputData{
			AddTransactionInputDto: body,
			ClientId:               entity.Id(id),
		})

		if err != nil {

			if err == domain.ErrClientNotFound {
				return c.SendStatus(fiber.StatusNotFound)
			} else if err == domain.ErrClientWithoutBalance {
				return c.SendStatus(fiber.StatusUnprocessableEntity)
			} else {
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		return c.Status(fiber.StatusOK).JSON(outputController)
	})
}

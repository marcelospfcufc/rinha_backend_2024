package rest

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/controller"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/infra/database"
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

func PostWrapper(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		transactionRepository := database.NewTransactionRepositoryGorm(db)
		clientRepository := database.NewClientRepositoryGorm(db)
		addTransactionService := service.NewAddTransactionService(
			clientRepository,
			transactionRepository,
		)
		ctrl := controller.NewAddTransactionController(*addTransactionService)

		id, _ := strconv.Atoi(c.Params("id"))
		var body controller.AddTransactionInputDto
		c.BodyParser(&body)

		if body.Operation != "c" && body.Operation != "d" {
			fmt.Println("Invalid Operation:", body.Operation)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if len(body.Description) < 1 || len(body.Description) > 10 {
			fmt.Println("Invalid Description:", body.Description)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		var outputController controller.AddTransactionOutputDto
		var err error

		outputController, err = ctrl.AddTransaction(controller.AddTransactionInputData{
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
	}
}

func GetWrapper(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		transactionRepository := database.NewTransactionRepositoryGorm(db)
		clientRepository := database.NewClientRepositoryGorm(db)

		getTransactionStatementService := service.NewGetTransactionStatementService(
			clientRepository,
			transactionRepository,
		)

		ctrl := controller.NewGetBankStatementController(*getTransactionStatementService)

		id, _ := strconv.Atoi(c.Params("id"))
		var body controller.AddTransactionInputDto
		c.BodyParser(&body)

		outputController, err := ctrl.GetBankStatement(controller.GetBankStatementInputDto{
			ClientId: entity.Id(id),
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
	}
}

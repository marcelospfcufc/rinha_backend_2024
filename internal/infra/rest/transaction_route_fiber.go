package rest

import (
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
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

func PostWrapper(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		return errors.New("not implemented yet")

		/* transactionRepository := database.NewTransactionRepositoryGorm(db)
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
			log.Info("Invalid Operation:", body.Operation)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if len(body.Description) < 1 || len(body.Description) > 10 {
			log.Info("Invalid Description:", body.Description)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		var outputController controller.AddTransactionOutputDto
		var err error

		outputController, err = ctrl.AddTransaction(
			controller.AddTransactionInputData{
				AddTransactionInputDto: body,
				ClientId:               entity.Id(id),
			},
		)

		if err != nil {

			if err == domain.ErrClientNotFound {
				return c.SendStatus(fiber.StatusNotFound)
			} else if err == domain.ErrClientWithoutBalance {
				return c.SendStatus(fiber.StatusUnprocessableEntity)
			} else {
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		return c.Status(fiber.StatusOK).JSON(outputController) */
	}
}

func GetWrapper(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		return errors.New("not implemented yet")

		/* transactionRepository :=
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
		*/
	}
}

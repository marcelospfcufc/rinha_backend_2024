package rest

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgxpool"
	addtransactionctrl "github.com/marcelospfcufc/rinha_backend_2024/internal/controller/add_transaction"
	bankstatementctrl "github.com/marcelospfcufc/rinha_backend_2024/internal/controller/get_bank_statement"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/infra/database/pgdatabase"
)

func PostWrapper(db *pgxpool.Pool, start chan<- string, finish chan<- string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var err error
		start <- c.Path()

		defer func() {
			finish <- c.Path()
		}()

		ctx, cancel := context.WithTimeout(c.Context(), time.Second*40)
		defer cancel()

		unitOfWork, err := pgdatabase.NewPgUnitOfWork(db)
		if err != nil {
			return err
		}

		ctrl := addtransactionctrl.NewAddTransactionController(
			unitOfWork,
		)

		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			log.Info("BadRequest id: ", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		var body addtransactionctrl.InputDto
		err = c.BodyParser(&body)
		if err != nil {
			log.Info("Failed to parse body", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if body.Operation != "c" && body.Operation != "d" {
			log.Info("Invalid Operation: ", "'", body.Operation, "'")
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		if len(body.Description) < 1 || len(body.Description) > 10 {
			log.Info("Invalid Description:", "'", body.Description, "'")
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		if body.Value < 0 {
			log.Info("Invalid Value:", "'", body.Value, "'")
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		var outputController addtransactionctrl.OutputDto

		outputController, err = ctrl.AddTransaction(
			ctx,
			addtransactionctrl.InputDto{
				Value:       body.Value,
				Operation:   body.Operation,
				Description: body.Description,
				ClientId:    entity.Id(id),
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

		return c.Status(fiber.StatusOK).JSON(outputController)
	}
}

func GetWrapper(db *pgxpool.Pool, start chan<- string, finish chan<- string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var err error
		start <- c.Path()

		defer func() {
			finish <- c.Path()
		}()

		unitOfWork, err := pgdatabase.NewPgUnitOfWork(db)
		if err != nil {
			return err
		}

		ctrl := bankstatementctrl.NewGetBankStatementController(
			unitOfWork,
		)

		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			log.Info("BadRequest id: ", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		outputController, err := ctrl.GetBankStatement(bankstatementctrl.InputDto{
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

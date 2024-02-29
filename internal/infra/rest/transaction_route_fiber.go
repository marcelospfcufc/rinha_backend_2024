package rest

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/controller"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/infra/database/pgdatabase"
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

func PostWrapper(db *sql.DB, queue chan *fiber.Ctx, start <-chan string, finish chan<- string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		queue <- c
		count := <-start
		log.Info(count)
		var err error

		defer func() {
			finish <- "Finalizando"
		}()

		ctx, cancel := context.WithTimeout(c.Context(), time.Second*40)
		defer cancel()

		unitOfWork, err := pgdatabase.NewPgUnitOfWork(db)
		if err != nil {
			return err
		}

		err = unitOfWork.Begin(ctx)
		if err != nil {
			return err
		}

		defer func() {
			if err != nil {
				err = unitOfWork.RollBack()
			} else {
				err = unitOfWork.Commit()
			}
		}()

		serviceAddTransaction := service.NewAddTransactionService(unitOfWork.GetRepository())
		ctrl := controller.NewAddTransactionController(unitOfWork, serviceAddTransaction)

		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			log.Info("BadRequest id: ", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		var body controller.AddTransactionInputDto
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

		var outputController controller.AddTransactionOutputDto

		outputController, err = ctrl.AddTransaction(
			ctx,
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

		return c.Status(fiber.StatusOK).JSON(outputController)
	}
}

func GetWrapper(db *sql.DB, queue chan *fiber.Ctx) fiber.Handler {
	return func(c *fiber.Ctx) error {

		queue <- c

		var err error

		unitOfWork, err := pgdatabase.NewPgUnitOfWork(db)
		if err != nil {
			return err
		}

		getTransactionStatementService := service.NewGetTransactionStatementService(
			unitOfWork.GetRepository(),
		)

		ctrl := controller.NewGetBankStatementController(
			unitOfWork,
			*getTransactionStatementService,
		)

		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			log.Info("BadRequest id: ", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

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

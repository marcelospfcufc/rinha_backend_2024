package rest

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/repository"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/infra/database/pgdatabase"
)

func UserExistsMiddlewareWrapper(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var err error

		unitOfWork, err := pgdatabase.NewPgUnitOfWork(db)
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(c.Params("id"))

		if err != nil {
			return domain.ErrClientNotFound
		}

		repo := repository.Repository(*unitOfWork.GetRepository())

		hasClient := repo.HasClientById(context.Background(), entity.Id(id))

		if !hasClient {
			return domain.ErrClientNotFound
		}

		return c.Next()
	}
}

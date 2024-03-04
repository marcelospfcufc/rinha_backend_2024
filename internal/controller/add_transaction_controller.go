package controller

import (
	"context"

	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/interfaces"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/service"
)

type AddTransactionInputDto struct {
	Value       int64  `json:"valor"`
	Operation   string `json:"tipo"`
	Description string `json:"descricao"`
}

type AddTransactionInputData struct {
	AddTransactionInputDto
	ClientId entity.Id
}

type AddTransactionOutputDto struct {
	Credit  int64 `json:"limite"`
	Balance int64 `json:"saldo"`
}

type AddTransactionController struct {
	UnitOfWork interfaces.UnitOfWork
	service    *service.AddTransactionService
}

func NewAddTransactionController(
	unitOfWork interfaces.UnitOfWork,
	service *service.AddTransactionService,
) *AddTransactionController {

	return &AddTransactionController{
		UnitOfWork: unitOfWork,
		service:    service,
	}
}

func (ctrl *AddTransactionController) AddTransaction(
	ctx context.Context,
	input AddTransactionInputData,
) (AddTransactionOutputDto, error) {

	var err error
	ctrl.UnitOfWork.Begin(ctx)
	defer func() {
		if err != nil {
			err = ctrl.UnitOfWork.RollBack(ctx)
		}
	}()

	serviceOutputData, err := ctrl.service.Execute(
		ctx,
		service.InputData{
			Value:       input.Value,
			ClientId:    input.ClientId,
			Operation:   input.Operation,
			Description: input.Operation,
		},
	)

	if err != nil {
		return AddTransactionOutputDto{}, err
	}

	ctrl.UnitOfWork.Commit(ctx)

	return AddTransactionOutputDto{
		Credit:  serviceOutputData.Credit,
		Balance: serviceOutputData.Balance,
	}, err
}

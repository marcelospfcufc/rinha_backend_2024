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
	UnitOfWork            interfaces.UnitOfWork
	AddTransactionService *service.AddTransactionService
}

func NewAddTransactionController(
	unitOfWork interfaces.UnitOfWork,
	addTransactionService *service.AddTransactionService,
) *AddTransactionController {

	return &AddTransactionController{
		UnitOfWork:            unitOfWork,
		AddTransactionService: addTransactionService,
	}
}

func (ctrl *AddTransactionController) AddTransaction(
	ctx context.Context,
	input AddTransactionInputData,
) (AddTransactionOutputDto, error) {

	out, err := ctrl.AddTransactionService.Execute(
		ctx,
		service.InputData{
			ClientId:    input.ClientId,
			Value:       input.Value,
			Operation:   input.Operation,
			Description: input.Description,
		},
	)

	if err != nil {
		return AddTransactionOutputDto{}, err
	}

	return AddTransactionOutputDto{
		Credit:  out.Credit,
		Balance: out.Balance,
	}, err
}

package controller

import (
	"context"

	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/interfaces"
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
}

func NewAddTransactionController(
	unitOfWork interfaces.UnitOfWork,
) *AddTransactionController {

	return &AddTransactionController{
		UnitOfWork: unitOfWork,
	}
}

func calculateNewBalance(
	clientCredit int64,
	currentBalance int64,
	transactionValue int64,
	transactionOperation string,
) (clientNewCurrentBalance int64, err error) {
	newBalanceValue := currentBalance
	if transactionOperation == "d" {
		newBalanceValue -= transactionValue

		if newBalanceValue < clientCredit*-1 {
			return -1, domain.ErrClientWithoutBalance
		}
	} else {
		newBalanceValue += transactionValue
	}

	return newBalanceValue, nil
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

	repository := ctrl.UnitOfWork.GetRepository()

	clientCredit, currentNewBalance, err := repository.AddTransaction(
		ctx,
		input.ClientId,
		entity.Transaction{
			Value:       input.Value,
			Operation:   input.Operation,
			Description: input.Description,
		},
		calculateNewBalance,
	)

	if err != nil {
		return AddTransactionOutputDto{}, err
	}

	ctrl.UnitOfWork.Commit(ctx)

	return AddTransactionOutputDto{
		Credit:  clientCredit,
		Balance: currentNewBalance,
	}, err
}

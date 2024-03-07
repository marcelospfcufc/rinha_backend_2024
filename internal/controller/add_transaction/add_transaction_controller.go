package addtransactionctrl

import (
	"context"

	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/interfaces"
)

type InputDto struct {
	Value       int64  `json:"valor"`
	Operation   string `json:"tipo"`
	Description string `json:"descricao"`
	ClientId    entity.Id
}

type OutputDto struct {
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

func (ctrl *AddTransactionController) AddTransaction(
	ctx context.Context,
	input InputDto,
) (OutputDto, error) {

	var err error
	ctrl.UnitOfWork.Begin(ctx)
	defer func() {
		if err != nil {
			err = ctrl.UnitOfWork.RollBack(ctx)
		}
	}()

	clientRepository := ctrl.UnitOfWork.GetClientRepository()

	clientCredit, currentNewBalance, err := clientRepository.AddTransaction(
		ctx,
		input.ClientId,
		entity.Transaction{
			Value:       input.Value,
			Operation:   input.Operation,
			Description: input.Description,
		},
		domain.CalculateNewBalance,
	)

	if err != nil {
		return OutputDto{}, err
	}

	ctrl.UnitOfWork.Commit(ctx)

	return OutputDto{
		Credit:  clientCredit,
		Balance: currentNewBalance,
	}, err
}

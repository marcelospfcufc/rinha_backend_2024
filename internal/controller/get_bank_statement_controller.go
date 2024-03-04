package controller

import (
	"context"
	"time"

	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/interfaces"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/service"
)

type GetBankStatementInputDto struct {
	ClientId entity.Id
}

type GetBankStatementOutputDto struct {
	Balance          ClientBalanceData `json:"saldo"`
	LastTransactions []TransactionData `json:"ultimas_transacoes"`
}

type TransactionData struct {
	Value       int64  `json:"valor"`
	Operation   string `json:"tipo"`
	Description string `json:"descricao"`
	CreatedAt   string `json:"realizada_em"`
}

type ClientBalanceData struct {
	Balance     int64  `json:"total"`
	RequestDate string `json:"data_extrato"`
	Credit      int64  `json:"limite"`
}

type GetBankStatementController struct {
	UnitOfWork                     interfaces.UnitOfWork
	getTransactionStatementService *service.GetTransactionStatementService
}

func NewGetBankStatementController(
	unitOfWork interfaces.UnitOfWork,
	service *service.GetTransactionStatementService,
) *GetBankStatementController {
	controller := GetBankStatementController{
		getTransactionStatementService: service,
	}

	return &controller
}

func (controller *GetBankStatementController) GetBankStatement(input GetBankStatementInputDto) (GetBankStatementOutputDto, error) {
	var transactionStatementOutputDto GetBankStatementOutputDto

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*40,
	)

	defer cancel()

	serviceOutput, err := controller.getTransactionStatementService.Execute(
		ctx,
		service.GetTransactionStatementInputData{
			ClientId: input.ClientId,
		},
	)

	if err != nil {
		return transactionStatementOutputDto, err
	}

	lastTransactionToReturn := make(
		[]TransactionData,
		len(serviceOutput.Transactions),
	)

	for idx, transaction := range serviceOutput.Transactions {
		lastTransactionToReturn[idx] = TransactionData{
			Value:       transaction.Value,
			Operation:   transaction.Operation,
			Description: transaction.Description,
			CreatedAt:   transaction.CreatedAt.Format("2006-01-02T15:04:05.999999Z"),
		}
	}

	transactionStatementOutputDto.Balance = ClientBalanceData{
		Credit:      serviceOutput.Credit,
		Balance:     serviceOutput.Balance,
		RequestDate: time.Now().UTC().Format("2006-01-02T15:04:05.999999Z"),
	}

	transactionStatementOutputDto.LastTransactions = lastTransactionToReturn

	return transactionStatementOutputDto, err
}

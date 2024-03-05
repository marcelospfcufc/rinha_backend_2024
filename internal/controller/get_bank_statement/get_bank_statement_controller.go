package bankstatementctrl

import (
	"context"
	"time"

	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/interfaces"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/repository"
)

type InputDto struct {
	ClientId entity.Id
}

type OutputDto struct {
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
	unitOfWork interfaces.UnitOfWork
}

func NewGetBankStatementController(
	unitOfWork interfaces.UnitOfWork,
) *GetBankStatementController {
	controller := GetBankStatementController{
		unitOfWork: unitOfWork,
	}

	return &controller
}

func (ctrl *GetBankStatementController) GetBankStatement(
	input InputDto,
) (
	OutputDto,
	error,
) {
	var transactionStatementOutputDto OutputDto

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*40,
	)

	defer cancel()

	clientRepository := ctrl.unitOfWork.GetClientRepository()

	client, err := clientRepository.GetClientById(
		ctx,
		input.ClientId,
	)

	if err != nil {
		return transactionStatementOutputDto, domain.ErrClientNotFound
	}

	lastTentransactions, err := clientRepository.GetTransactions(
		ctx,
		input.ClientId,
		10,
		repository.Desc,
	)

	if err != nil {
		return transactionStatementOutputDto, err
	}

	lastTransactionToReturn := make(
		[]TransactionData,
		len(lastTentransactions),
	)

	for idx, transaction := range lastTentransactions {
		lastTransactionToReturn[idx] = TransactionData{
			Value:       transaction.Value,
			Operation:   transaction.Operation,
			Description: transaction.Description,
			CreatedAt:   transaction.CreatedAt.Format("2006-01-02T15:04:05.999999Z"),
		}
	}

	transactionStatementOutputDto.Balance = ClientBalanceData{
		Credit:      client.Credit,
		Balance:     client.CurrentBalance,
		RequestDate: time.Now().UTC().Format("2006-01-02T15:04:05.999999Z"),
	}

	transactionStatementOutputDto.LastTransactions = lastTransactionToReturn
	return transactionStatementOutputDto, err
}

package database

import (
	"context"

	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain"
	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
	"gorm.io/gorm"
)

type ClientRepositoryGorm struct {
	dbConnection *gorm.DB
}

func NewClientRepositoryGorm(dbConnection *gorm.DB) *ClientRepositoryGorm {
	return &ClientRepositoryGorm{dbConnection: dbConnection}
}

func (repo *ClientRepositoryGorm) HasClientById(ctx context.Context, clientId entity.Id) bool {
	var clientFound Client
	result := repo.dbConnection.WithContext(ctx).First(&clientFound, clientId)

	if result.Error == nil && result.RowsAffected > 0 {
		return true
	}

	return false
}

func (repo *ClientRepositoryGorm) Create(client entity.Client) (entity.Client, error) {
	clientModel := Client{
		Name:           client.Name,
		Credit:         client.Credit,
		CurrentBalance: client.CurrentBalance,
	}

	result := repo.dbConnection.Create(&clientModel)

	if result.Error != nil {
		return entity.Client{}, result.Error
	}

	return entity.Client{
		Id:             clientModel.ID,
		Name:           clientModel.Name,
		Credit:         clientModel.Credit,
		CurrentBalance: clientModel.CurrentBalance,
	}, result.Error
}

func (repo ClientRepositoryGorm) GetById(ctx context.Context, clientId entity.Id) (entity.Client, error) {
	var clientFound Client

	/*ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	defer func() {
		if result.Error != nil {
			log.Info("GetById - ", clientId, " - ", result.Error)
		}
	}()*/

	//result := repo.dbConnection.WithContext(ctx).Preload("Transactions").First(&clientFound, clientId)
	result := repo.dbConnection.WithContext(ctx).First(&clientFound, clientId)
	if result.Error != nil {
		return entity.Client{}, domain.ErrClientNotFound
	}

	/*var transactionToReturn []entity.Transaction = make([]entity.Transaction, len(clientFound.Transactions)+1)

	for _, dbTransaction := range clientFound.Transactions {
		transactionToReturn = append(
			transactionToReturn,
			entity.Transaction{
				Id:          dbTransaction.ID,
				Value:       dbTransaction.Value,
				Operation:   dbTransaction.Operation,
				Description: dbTransaction.Description,
				CreatedAt:   dbTransaction.CreatedAt,
			},
		)
	}*/

	return entity.Client{
		Id:           clientFound.ID,
		Name:         clientFound.Name,
		Credit:       clientFound.Credit,
		Transactions: []entity.Transaction{},
	}, nil
}

func (repo ClientRepositoryGorm) GetTransactionsById(ctx context.Context, clientId entity.Id) ([]entity.Transaction, error) {

	var dbTransactionsByClient []Transaction = []Transaction{}
	result := repo.dbConnection.WithContext(ctx).Where("client_id = ?", clientId).Find(&dbTransactionsByClient)
	if result.Error != nil {
		return []entity.Transaction{}, result.Error
	}

	var transactionToReturn []entity.Transaction = make([]entity.Transaction, len(dbTransactionsByClient)+1)

	for _, dbTransaction := range dbTransactionsByClient {
		transactionToReturn = append(
			transactionToReturn,
			entity.Transaction{
				Id:          dbTransaction.ID,
				Value:       dbTransaction.Value,
				Operation:   dbTransaction.Operation,
				Description: dbTransaction.Description,
				CreatedAt:   dbTransaction.CreatedAt,
			},
		)
	}

	return transactionToReturn, nil

}

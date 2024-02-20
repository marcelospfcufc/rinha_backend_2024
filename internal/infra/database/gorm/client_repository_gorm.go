package database

import (
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

func (repo *ClientRepositoryGorm) Create(client entity.Client) (entity.Client, error) {
	clientModel := Client{
		Name:   client.Name,
		Credit: client.Credit,
	}

	result := repo.dbConnection.Create(&clientModel)

	if result.Error != nil {
		return entity.Client{}, result.Error
	}

	return entity.Client{
		Id:     clientModel.ID,
		Name:   clientModel.Name,
		Credit: clientModel.Credit,
	}, result.Error
}

func (repo ClientRepositoryGorm) GetById(clientId entity.Id) (entity.Client, error) {

	var clientFound Client
	result := repo.dbConnection.Preload("Transactions").First(&clientFound, clientId)

	if result.Error != nil {
		return entity.Client{}, domain.ErrClientNotFound
	}

	var transactionToReturn []entity.Transaction = make([]entity.Transaction, len(clientFound.Transactions)+1)

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
	}

	return entity.Client{
		Id:           clientFound.ID,
		Name:         clientFound.Name,
		Credit:       clientFound.Credit,
		Transactions: transactionToReturn,
	}, nil
}

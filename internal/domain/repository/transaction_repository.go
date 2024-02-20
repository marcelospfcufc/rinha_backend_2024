package repository

import "github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"

type TransactionRepositoryData struct {
	TransactionId TransactionIdType
	Value         ValueType
	Operation     OperationType
	RealizedIn    RealizedInType
	Description   DescriptionType
	ClientId      any
}

type TransactionRepository interface {
	Create(clientId entity.Id, transaction entity.Transaction) (entity.Transaction, error)
	GetAllByUser(clientId entity.Id) ([]entity.Transaction, error)
}

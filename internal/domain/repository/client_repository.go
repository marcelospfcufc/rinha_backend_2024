package repository

import "github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"

type ClientRepository interface {
	HasClientById(clientId entity.Id) bool
	Create(client entity.Client) (entity.Client, error)
	GetById(clientId entity.Id) (entity.Client, error)
}

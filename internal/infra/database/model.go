package database

import (
	"time"

	"github.com/marcelospfcufc/rinha_backend_2024/internal/domain/entity"
)

type Client struct {
	ID             entity.Id `gorm:"primaryKey"`
	CreatedAt      time.Time
	Name           string
	Credit         int64
	Transactions   []Transaction `gorm:"foreignkey:ClientID"`
	CurrentBalance int64
}

type Transaction struct {
	ID          entity.Id `gorm:"primaryKey"`
	Operation   string
	Value       int64
	CreatedAt   time.Time
	ClientID    entity.Id
	Description string
}

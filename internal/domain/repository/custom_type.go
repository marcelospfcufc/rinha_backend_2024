package repository

import (
	"errors"
	"time"
)

type UserIdType uint8
type TransactionIdType uint64
type ValueType int32

type RealizedInType *time.Time
type DescriptionType string
type CtxDbKey string

type OrderBy int

const (
	Asc OrderBy = iota
	Desc
)

const (
	DbKey CtxDbKey = CtxDbKey("DbConn")
)

const (
	debit  = "d"
	credit = "c"
)

type OperationType string

func (op OperationType) IsValid() bool {
	return op == debit || op == credit
}

func ParseOperationType(value string) (OperationType, error) {
	switch value {
	case debit, credit:
		return OperationType(value), nil
	default:
		return "", errors.New("invalid value to OperationType")
	}
}

func (order OrderBy) String() string {
	orders := [...]string{"asc", "desc"}
	if order < Asc || order > Desc {
		panic("OrderBy with invalid value")
	}
	return orders[order]
}

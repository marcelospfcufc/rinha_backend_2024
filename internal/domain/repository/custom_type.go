package repository

import "time"

type UserIdType uint8
type TransactionIdType uint64
type ValueType int32
type OperationType string
type RealizedInType *time.Time
type DescriptionType string

type OrderBy int

const (
	Asc OrderBy = iota
	Desc
)

func (order OrderBy) String() string {
	orders := [...]string{"asc", "desc"}
	if order < Asc || order > Desc {
		panic("OrderBy with invalid value")
	}
	return orders[order]
}

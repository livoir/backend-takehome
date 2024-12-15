package domain

import "database/sql"

type Transaction interface {
	Commit() error
	Rollback() error
	GetTx() *sql.Tx
}

type Transactor interface {
	Begin() (Transaction, error)
}
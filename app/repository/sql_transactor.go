package repository

import (
	"app/domain"
	"database/sql"
)

type SQLTransaction struct {
	tx *sql.Tx
}

func (t *SQLTransaction) Commit() error {
	return t.tx.Commit()
}

func (t *SQLTransaction) Rollback() error {
	return t.tx.Rollback()
}

func (t *SQLTransaction) GetTx() *sql.Tx {
	return t.tx
}

type SQLTransactor struct {
	db *sql.DB
}

func NewSQLTransactor(db *sql.DB) domain.Transactor {
	return &SQLTransactor{db: db}
}

func (t *SQLTransactor) Begin() (domain.Transaction, error) {
	tx, err := t.db.Begin()
	if err != nil {
		return nil, err
	}
	return &SQLTransaction{tx: tx}, nil
}

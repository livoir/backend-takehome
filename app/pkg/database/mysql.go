package database

import (
	"app/pkg/logger"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

func NewMysqlConnection(host, port, dbName, user, password string) (*sql.DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", user, password, host, port, dbName)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		logger.Log.Error("Error opening database connection: ", zap.Error(err))
		return nil, err
	}
	return db, nil
}

func RunMigration(db *sql.DB) error {
	goose.SetDialect("mysql")
	err := goose.Up(db, "db/migrations")
	if err != nil {
		logger.Log.Error("Error running migration: ", zap.Error(err))
		return err
	}
	return nil
}

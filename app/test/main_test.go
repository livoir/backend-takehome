package test

import (
	"app/delivery/http"
	"app/pkg/database"
	"app/pkg/logger"
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

var db *sql.DB
var router *gin.Engine
var mysqlContainer testcontainers.Container

func TestMain(m *testing.M) {
	// Setup
	// Start MySQL container
	err := logger.Init()
	if err != nil {
		panic(err)
	}
	gin.SetMode(gin.TestMode)
	ctx := context.Background()
	mysqlContainer, err = mysql.Run(ctx, "mysql:8.0")
	if err != nil {
		panic(err)
	}
	defer mysqlContainer.Terminate(ctx)
	mysqlHost, err := mysqlContainer.Host(ctx)
	if err != nil {
		panic(err)
	}
	mysqlPort, err := mysqlContainer.MappedPort(ctx, "3306")
	if err != nil {
		panic(err)
	}
	db, err = database.NewMysqlConnection(mysqlHost, mysqlPort.Port(), "test", "root", "test")
	if err != nil {
		panic(err)
	}
	err = database.RunMigration(db, "../db/migrations")
	if err != nil {
		panic(err)
	}

	// Read private and public key
	privateKey, err := os.ReadFile("../cert/backend_takehome_rsa")
	if err != nil {
		panic(err)
	}
	publicKey, err := os.ReadFile("../cert/backend_takehome_rsa.pub")
	if err != nil {
		panic(err)
	}

	// Setup router
	router, err = http.SetupRouter(db, string(privateKey), string(publicKey))
	if err != nil {
		panic(err)
	}
	code := m.Run()

	// Teardown
	os.Exit(code)
}

package main

import (
	"app/delivery/http"
	"app/pkg/database"
	"app/pkg/logger"
	"fmt"
	"os"
)

func main() {
	if err := logger.Init(); err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		return
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Printf("Error syncing logger: %v\n", err)
		}
	}()
	mysqlUser := os.Getenv("BACKEND_TAKE_HOME_MYSQL_USER")
	mysqlPassword := os.Getenv("BACKEND_TAKE_HOME_MYSQL_PASSWORD")
	mysqlHost := os.Getenv("BACKEND_TAKE_HOME_MYSQL_HOST")
	mysqlPort := os.Getenv("BACKEND_TAKE_HOME_MYSQL_PORT")
	mysqlDatabase := os.Getenv("BACKEND_TAKE_HOME_MYSQL_DATABASE")

	jwtPrivateKeyPath := os.Getenv("BACKEND_TAKE_HOME_JWT_PRIVATE_KEY_PATH")
	jwtPublicKeyPath := os.Getenv("BACKEND_TAKE_HOME_JWT_PUBLIC_KEY_PATH")

	db, err := database.NewMysqlConnection(mysqlHost, mysqlPort, mysqlDatabase, mysqlUser, mysqlPassword)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	if err := database.RunMigration(db, "db/migrations"); err != nil {
		logger.Log.Error(err.Error())
		return
	}

	privateKey, err := os.ReadFile(jwtPrivateKeyPath)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	publicKey, err := os.ReadFile(jwtPublicKeyPath)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	router, err := http.SetupRouter(db, string(privateKey), string(publicKey))
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	if err := router.Run(":8080"); err != nil {
		logger.Log.Error(err.Error())
		return
	}
}

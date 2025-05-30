package main

import (
	"fmt"
	config "github.com/MarlonG1/delivery-backend/configs"
	"github.com/MarlonG1/delivery-backend/configs/database"
	"github.com/MarlonG1/delivery-backend/internal/bootstrap"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/server"
	infraDB "github.com/MarlonG1/delivery-backend/internal/infrastructure/database"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/websocket"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"

	_ "github.com/MarlonG1/delivery-backend/docs/swagger"
)

func main() {
	envConfig, err := config.NewEnvConfig()
	if err != nil {
		fmt.Println("Error loading environment variables" + err.Error())
		return
	}

	err = logs.InitLogger(envConfig)
	if err != nil {
		logs.Fatal("Error initializing logger", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	dbConnection, err := initDatabaseConfigurations(envConfig)
	if err != nil {
		return
	}

	wsHub := websocket.NewHub()
	go wsHub.Run()

	container := bootstrap.NewContainer(dbConnection.Db, envConfig, wsHub)
	apiSv := server.NewAPIServer(container, envConfig)

	if err = apiSv.Start(); err != nil {
		logs.Fatal("Error starting server", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	defer func(dbConnection *database.DbConnection) {
		err := dbConnection.Close()
		if err != nil {
			logs.Error("Main method, on close connection", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}(dbConnection)
}

func initDatabaseConfigurations(envConfig *config.EnvConfig) (*database.DbConnection, error) {
	dbDriver := selectDatabaseDriver(envConfig)
	dbConnection := database.NewDatabaseConnection(dbDriver)
	if err := dbConnection.Open(); err != nil {
		return nil, err
	}

	if err := infraDB.RunMigrations(dbConnection.Db); err != nil {
		logs.Fatal("Error running migrations", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	return dbConnection, nil
}

func selectDatabaseDriver(envConfig *config.EnvConfig) database.DriverConfig {
	switch envConfig.Database.Driver {
	case "mysql":
		return database.NewMysqlDriver(envConfig)
	case "postgres":
		return database.NewPostgresDriver(envConfig)
	default:
		logs.Fatal("Invalid database driver type", map[string]interface{}{
			"database_type": envConfig.Database.Driver,
		})
		return nil
	}
}

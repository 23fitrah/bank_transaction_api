package main

import (
	"fmt"
	"log"
	"time"

	"database/sql"
	"transaction_api/config"
	"transaction_api/handler"
	"transaction_api/repository"
	"transaction_api/routes"
	"transaction_api/service"
	"transaction_api/utils"

	"github.com/gin-gonic/gin"
)

type databaseConnections struct {
	mssql *sql.DB
}

type appServices struct {
	accountService       *service.AccountService
	transactionService   *service.TransactionService
	logMonitoringService *service.LogMonitoringService
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbs, err := initDatabases(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize databases: %v", err)
	}
	defer closeDatabases(dbs)

	initElasticsearch(cfg)

	services := initServices(dbs)

	startHTTPServer(cfg, services)
}

func initDatabases(cfg *config.Config) (*databaseConnections, error) {
	dbs := &databaseConnections{}

	dbMssql, err := config.NewMSSQLConnection(cfg.DatabaseMsql)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQL Server DB Bank: %w", err)
	}
	dbs.mssql = dbMssql
	log.Println("Connected to SQL Server DB Bank")

	return dbs, nil
}

func closeDatabases(dbs *databaseConnections) {

	if dbs.mssql != nil {
		dbs.mssql.Close()
	}

}

func initElasticsearch(cfg *config.Config) {
	if cfg.Logging.ElasticSearch.Enabled {
		_, err := config.NewElasticsearchConnection(config.ESConfig{
			URL:      cfg.Logging.ElasticSearch.URL,
			Username: cfg.Logging.ElasticSearch.Username,
			Password: cfg.Logging.ElasticSearch.Psswd,
			Timeout:  10 * time.Second,
		})
		if err != nil {
			log.Printf("Warning: Failed to connect to Elasticsearch: %v", err)
		} else {
			log.Println("Connected to Elasticsearch")
		}
	}
}

func initServices(dbs *databaseConnections) *appServices {
	accountRepo := repository.NewAccountRepository(dbs.mssql)
	transactionRepo := repository.NewTransactionRepository(dbs.mssql)
	logRepo := repository.NewLogMonitoringRepository(dbs.mssql)

	return &appServices{
		accountService:       service.NewAccountService(accountRepo),
		transactionService:   service.NewTransactionService(transactionRepo),
		logMonitoringService: service.NewLogMonitoringService(logRepo),
	}
}

func startHTTPServer(cfg *config.Config, services *appServices) {
	r := gin.New()

	r.Use(utils.LogMonitoringMiddleware(services.logMonitoringService))

	accountHandler := handler.NewAccountHandler(services.accountService)
	transactionHandler := handler.NewTransactionHandler(services.transactionService, services.accountService)

	routes.RegisterAccountRoutes(r, accountHandler)
	routes.RegisterTransactionRoutes(r, transactionHandler)

	log.Printf("Server starting on port %s...", cfg.Server.Port)
	if err := r.Run(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

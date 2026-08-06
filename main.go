package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
	"transaction_api/config"
	"transaction_api/handler"
	"transaction_api/repository"
	"transaction_api/routes"
	"transaction_api/service"
	"transaction_api/utils"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/redis/go-redis/v9"
)

type databaseConnections struct {
	mssql *sql.DB
}

type redisConnections struct {
	redisdb *redis.Client
}
type elasticConnection struct {
	elasticdb *elastic.Client
}

type rabbitMqConnection struct {
	rmqclient *config.RabbitMQ
}
type appServices struct {
	accountService       *service.AccountService
	transactionService   *service.TransactionService
	logMonitoringService *service.LogMonitoringService
	userService          *service.UserService
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

	esdb, err := initElasticsearch(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize elastic: %v", err)
	}

	rds, err := initRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	rmq, err := initRabbitMq(cfg)
	if err != nil {
		log.Printf("Failed to initialize Rabbit MQ: %v", err)
		//log.Fatalf("Failed to initialize Rabbit MQ: %v", err)
	}

	services := initServices(dbs, rds, esdb)

	utils.NewOutboxWorker(dbs.mssql, rmq.rmqclient, cfg).Start(context.Background())

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

func initElasticsearch(cfg *config.Config) (*elasticConnection, error) {
	es := &elasticConnection{}
	var err error

	if cfg.Logging.ElasticSearch.Enabled {
		es.elasticdb, err = config.NewElasticsearchConnection(config.ESConfig{
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
	return es, err
}

func initRedis(cfg *config.Config) (*redisConnections, error) {
	var err error
	rdb := &redisConnections{}
	rdb.redisdb, err = config.ConnectRedis()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return rdb, err
}

func initRabbitMq(cfg *config.Config) (*rabbitMqConnection, error) {
	var err error
	rbmq := &rabbitMqConnection{}
	rbmq.rmqclient, err = config.NewRabbitMQ(cfg.RabbitMQ)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Rabbit MQ: %w", err)
	}
	return rbmq, err
}

func initServices(dbs *databaseConnections, rdb *redisConnections, esdb *elasticConnection) *appServices {
	accountRepo := repository.NewAccountRepository(dbs.mssql)
	transactionRepo := repository.NewTransactionRepository(dbs.mssql)
	logRepo := repository.NewLogMonitoringRepository(dbs.mssql, rdb.redisdb, esdb.elasticdb)
	userRepo := repository.NewUserRepository(dbs.mssql)

	return &appServices{
		accountService:       service.NewAccountService(accountRepo),
		transactionService:   service.NewTransactionService(transactionRepo),
		logMonitoringService: service.NewLogMonitoringService(logRepo),
		userService:          service.NewUserService(userRepo),
	}
}

func startHTTPServer(cfg *config.Config, services *appServices) {
	r := gin.New()

	userHandler := handler.NewUserHandler(services.userService)
	rAuth := r.Group("/api/v1/auth")
	rAuth.Use(utils.LogMonitoringMiddleware(services.logMonitoringService))
	{
		routes.RegisterUserRoutes(rAuth, userHandler)
	}

	accountHandler := handler.NewAccountHandler(services.accountService)
	transactionHandler := handler.NewTransactionHandler(services.transactionService, services.accountService)
	logHandler := handler.NewLogHandler(services.logMonitoringService)
	r.Use(utils.AuthMiddlewareGin(services.userService), utils.LogMonitoringMiddleware(services.logMonitoringService))
	{
		routes.RegisterAccountRoutes(r, accountHandler)
		routes.RegisterTransactionRoutes(r, transactionHandler)
		routes.RegisterLogMonitoringRoutes(r, logHandler)
	}

	log.Printf("Server starting on port %s...", cfg.Server.Port)
	if err := r.Run(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

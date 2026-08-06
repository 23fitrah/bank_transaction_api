package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server        ServerConfig   `json:"server"`
	DatabaseMysql DatabaseConfig `json:"database"`
	DatabaseMsql  DatabaseConfig `json:"databaseMsql"`
	Logging       LoggingConfig  `json:"logging"`
	RedisDB       RedisConfig    `json:"redisdb"`
	RabbitMQ      RabbitMQConfig `json:"rabbitmq"`
}

type LoggingConfig struct {
	ElasticSearch ElasticearchConfig `json:"elasticSearch"`
	Level         string             `json:"level"`
}

type ElasticearchConfig struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Psswd    string `json:"password"`
	Enabled  bool   `json:"enabled"`
}

type ServerConfig struct {
	Port            string `json:"port"`
	ReadTimeout     int    `json:"readTimeout"`
	WriteTimeout    int    `json:"writeTimeout"`
	ShutdownTimeout int    `json:"shutdownTimeout"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type RedisConfig struct {
	RedisHost     string `json:"redishost"`
	RedisPort     string `json:"redisport"`
	RedisUsername string `json:"redisusername"`
	RedisPassword string `json:"redispassword"`
	RedisDatabase string `json:"redisdatabase"`
}

type RabbitMQConfig struct {
	UrlRMQ   string `json:"url"`
	Exchange string `json:"exchange"`
	Routing  string `json:"exchange"`
	Queue    string `json:"Queue"`
}

func Load() (*Config, error) {
	_ = godotenv.Load() // Load .env file

	readTimeout, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
	writeTimeout, _ := strconv.Atoi(os.Getenv("SERVER_WRITE_TIMEOUT"))
	shutdownTimeout, _ := strconv.Atoi(os.Getenv("SERVER_SHUTDOWN_TIMEOUT"))

	return &Config{
		Server: ServerConfig{
			Port:            os.Getenv("SERVER_PORT"),
			ReadTimeout:     readTimeout,
			WriteTimeout:    writeTimeout,
			ShutdownTimeout: shutdownTimeout,
		},
		DatabaseMsql: DatabaseConfig{
			Host:     os.Getenv("MSSQL_DB_HOST"),
			Port:     os.Getenv("MSSQL_DB_PORT"),
			Username: os.Getenv("MSSQL_DB_USER"),
			Password: os.Getenv("MSSQL_DB_PASS"),
			Database: os.Getenv("MSSQL_DB_NAME"),
		},
		Logging: LoggingConfig{
			ElasticSearch: ElasticearchConfig{
				URL:      os.Getenv("ES_URL"),
				Username: os.Getenv("ES_USER"),
				Psswd:    os.Getenv("ES_PASS"),
				Enabled:  os.Getenv("ES_ENABLED") == "true",
			},
			Level: os.Getenv("LOG_LEVEL"),
		},
		RedisDB: RedisConfig{
			RedisHost:     os.Getenv("REDIS_HOST"),
			RedisPort:     os.Getenv("REDIS_PORT"),
			RedisUsername: "",
			RedisPassword: os.Getenv("REDIS_PASSWORD"),
			RedisDatabase: os.Getenv("REDIS_DATABASE"),
		},
		RabbitMQ: RabbitMQConfig{
			UrlRMQ:   os.Getenv("RABBITMQ_URL"),
			Exchange: os.Getenv("RABBITMQ_EXCHANGE"),
			Routing:  os.Getenv("RABBITMQ_ROUTING"),
			Queue:    os.Getenv("RABBITMQ_QUEUE"),
		},
	}, nil
}

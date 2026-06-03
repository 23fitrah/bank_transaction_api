package config

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

func NewMSSQLConnection(config DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func GetParameterES(jenis string) (string, error) {
	cfg, err := Load()
	if err != nil {
		return "", fmt.Errorf("Failed to load config: %v", err)
	}
	dbMssql, err := NewMSSQLConnection(cfg.DatabaseMsql)
	if err != nil {
		return "", fmt.Errorf("failed to connect to SQL Server "+cfg.DatabaseMsql.Database+": %w", err.Error())
	}

	const query = `
        SELECT Status
        FROM [W_PARAMETER]
        WHERE Jenis = @p1
    `
	var Status string

	if err := dbMssql.QueryRow(query, jenis).Scan(&Status); err != nil {

		return "", fmt.Errorf("failed to check Type W_PARAMETER: %w - %w", jenis, err)
	}

	if len(Status) > 0 {
		return Status, nil
	} else {
		return "", fmt.Errorf("Type not found : %w", jenis)
	}

}

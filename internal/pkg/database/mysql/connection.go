package mysql

import (
	"database/sql"
	"fmt"
	"time"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

type Connection struct {
	*sql.DB
}

func New(config Config) (*Connection, error) {
	db, err := sql.Open("mysql", config.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("mysql failed to connect: %w", err)
	}

	db.SetConnMaxLifetime(config.ConnectionMaxLifetime)
	if config.ConnectionMaxIdleTime > time.Second*0 {
		db.SetConnMaxIdleTime(config.ConnectionMaxIdleTime)
	}
	db.SetMaxIdleConns(config.MaxIdleConnections)
	db.SetMaxOpenConns(config.MaxOpenConnections)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("mysql failed to ping: %w", err)
	}

	return &Connection{
		DB: db,
	}, nil
}

package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func New(addr string, maxOpenConns, maxIdleConns int, connMaxLifetime string, connMaxIdleTime string) *sql.DB {
	sql, err := sql.Open("postgres", addr)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	sql.SetMaxOpenConns(maxOpenConns)
	sql.SetMaxIdleConns(maxIdleConns)
	connMaxLifetimeDuration, err := time.ParseDuration(connMaxLifetime)
	if err != nil {
		fmt.Println("Invalid connection max lifetime duration:", err)
	}
	connMaxIdleTimeDuration, err := time.ParseDuration(connMaxIdleTime)
	if err != nil {
		fmt.Println("Invalid connection max idle time duration:", err)
	}
	sql.SetConnMaxLifetime(connMaxLifetimeDuration)
	sql.SetConnMaxIdleTime(connMaxIdleTimeDuration)

	if err := sql.PingContext(context.Background()); err != nil {
		panic("Failed to ping database: " + err.Error())
	}
	return sql
}

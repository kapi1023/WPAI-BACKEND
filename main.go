package main

import (
	"airplane/service"
	"context"
	"log"
	"os"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func main() {
	dbConn := getEnv("SQL", true)

	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	db, err = sqlx.ConnectContext(ctx, "mssql", dbConn)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(5)
	defer db.Close()

	service.Start(db)
}

func getEnv(key string, must ...bool) string {
	value := os.Getenv(key)
	if value == "" {
		if len(must) > 0 && !must[0] {
			return ""
		}
		log.Fatalf("Environment variable %s is required", key)
	}
	return value
}

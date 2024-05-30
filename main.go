package main

import (
	"airplane/service"
	"database/sql"
	"log"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", "user=youruser dbname=yourdb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	service.Start(db)
}

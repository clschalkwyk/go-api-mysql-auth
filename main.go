package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("Start")
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3309)/authdb")
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	var version string
	err = db.QueryRow("select version()").Scan(&version)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(version)
}

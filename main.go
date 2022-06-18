package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var mydb = ConnectDb()

func ConnectDb() *sql.DB {

	fmt.Println("Start")
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASS"),
			os.Getenv("MYSQL_HOST"),
			os.Getenv("MYSQL_PORT"),
			os.Getenv("MYSQL_DB"),
		))

	if err != nil {
		log.Fatal(err)
	}
	return db
}

func Ping(c *gin.Context) {
	var version string
	err := mydb.QueryRow("select version()").Scan(&version)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(version)

	c.JSON(200, gin.H{
		"message": version,
	})
}

// register user
// login
// generate jwt
// validate jwt

func main() {

	r := gin.Default()
	r.GET("/ping", Ping)
	r.Run("0.0.0.0:8060")
}

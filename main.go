package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type RegisterUserDto struct {
	Email           string `json: "email"`
	Password        string `json: "password"`
	ConfirmPassword string `json: "confirmPassword"`
}

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

func hashPassword(pwd []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	return hash, err
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

func ExistingUser(email string) bool {
	var cnt int
	if err := mydb.QueryRow("select count(1) cnt from users where email = ?", email).Scan(&cnt); err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatal(err)
	}
	return cnt > 0
}

func RegisterUser(c *gin.Context) {
	data := RegisterUserDto{}
	if err := c.BindJSON(&data); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if data.Password != data.ConfirmPassword {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if ExistingUser(data.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"result":  "Error",
			"message": "User exists",
		})
		return
	}

	hashed, err := hashPassword([]byte(data.Password))
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusAccepted, gin.H{
			"result":  "Failed",
			"message": "Unable to register user",
		})
	}
	data.Password = string(hashed)

	q := "insert into users (email, password) values(?, ?)"
	insert, err := mydb.Prepare(q)
	if err != nil {
		log.Fatal(err)
	}
	_, err = insert.Exec(data.Email, data.Password)
	insert.Close()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(data)
	c.JSON(http.StatusAccepted, gin.H{
		"result":  "Success",
		"message": "User registered",
	})
}

// register user
// login
// generate jwt
// validate jwt

func main() {

	r := gin.Default()
	r.GET("/ping", Ping)
	r.POST("/register", RegisterUser)
	r.Run("0.0.0.0:8060")
}

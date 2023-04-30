package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func ConnectDataBase() *sql.DB {
	var DB *sql.DB
	Dbdriver, ok := os.LookupEnv("DB_DRIVER")
	if !ok {
		log.Fatal("please specify DB_DRIVER")
	}
	DbHost, ok := os.LookupEnv("DB_HOST")
	if !ok {
		log.Fatal("please specify DB_HOST")
	}
	DbUser, ok := os.LookupEnv("DB_USER")
	if !ok {
		log.Fatal("please specify DB_USER")
	}
	DbPassword, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		log.Fatal("please specify DB_PASSWORD")
	}
	DbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		log.Fatal("please specify DB_NAME")
	}
	DbPort, ok := os.LookupEnv("DB_PORT")
	if !ok {
		log.Fatal("please specify DB_PORT")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", DbHost, DbUser, DbPassword, DbName, DbPort)
	var err error
	for i := 0; i < 5; i++ {
		DB, err = sql.Open(Dbdriver, dsn)
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
		log.Printf("Error connecting to database: %v, retrying...", err)
	}

	if err != nil {
		log.Println("Cannot connect to database ", Dbdriver)
		log.Fatal("connection error:", err)
	} else {
		log.Println("We are connected to the database ", Dbdriver)
	}

	defer DB.Close()
	fmt.Println("Welcome!")
	return DB
}

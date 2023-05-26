package database

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
	return DB
}

func ClearData(db *sql.DB) error {
	_, err := db.Query("TRUNCATE TABLE  users, games, genres, publishers CASCADE;")
	if err != nil {
		return err
	}
	_, err = db.Query("INSERT INTO users(id, email, username, role, hashedPassword, badgeColor) values('699a5565-4c7e-4d18-be3e-ea04eb4f5e4d', 'admin@a.a', 'admin', 'admin', '$argon2id$v=19$m=65536,t=3,p=4$iZAGQSDhbte1+l0oF+rD/g$QWVeHaFYiR8iUA1BvK9+Pkua9EV3K/6y6CMTqSSet4Y', '' ) ON CONFLICT DO NOTHING;")
	return err
}

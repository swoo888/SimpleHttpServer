package dbconn

import (
	"fmt"
	"log"
	_ "github.com/lib/pq"
	"database/sql"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "BlueSky8@3"
	//DB_PASSWORD = "RedPocket123"
	DB_NAME     = "postgres"
)

func NewDbConnection() *sql.DB {
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Println("error opening dbconn. ", err)
		return nil
	}
	return db
}

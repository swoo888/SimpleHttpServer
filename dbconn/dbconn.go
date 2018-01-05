package dbconn

import (
	"fmt"
	"log"
	_ "github.com/lib/pq"
	"database/sql"
	"time"
)

const (
	HOST_NAME   = "postgresdb"
	DB_USER     = "postgres"
	DB_PASSWORD = "BlueSky8@3"
	//DB_PASSWORD = "RedPocket123"
	DB_NAME = "postgres"
)

func NewDbConnection() (db *sql.DB) {
	dbInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		HOST_NAME, DB_USER, DB_PASSWORD, DB_NAME, )

	maxWait := 10
	var err error
	db = nil
	for {
		db, err = sql.Open("postgres", dbInfo)
		if err != nil {
			// wait for postgres to be ready
			if (maxWait >= 0) {
				log.Println("wating for postgres, sleeping...")
				time.Sleep(1 * time.Second)
				maxWait -= 1
			} else {
				return nil
			}
		} else {
			return db
		}
	}
	return nil
}

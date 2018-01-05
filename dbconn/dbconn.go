package dbconn

import (
	"fmt"
	"log"
	"database/sql"
	"time"
	"flag"
)

var hostName = flag.String("host", "localhost", "postgres hostname")
var dbName = flag.String("db", "", "postgres database")
var dbUserName = flag.String("user", "", "postgres user")
var dbPassword = flag.String("pass", "", "postgres password");

var dbConn *sql.DB

// GetMainDbConn returns the main db connection
func GetMainDbConn() *sql.DB {
	if dbConn == nil {
		waitForDb()
		dbConn = newDbConnection()
	}
	return dbConn
}

// newDbConnection will return a new postgres database connection.  it will try multiple tries
// to obtain a db connection
func newDbConnection() (db *sql.DB) {
	if *dbName == "" || *hostName == "" || *dbUserName == "" || *dbPassword == "" {
		log.Println("postgres credentials not set")
		flag.PrintDefaults()
		return nil
	}
	dbInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		*hostName, *dbUserName, *dbPassword, *dbName, )

	maxWait := 3
	var err error
	for {
		db, err = sql.Open("postgres", dbInfo)
		if err != nil {
			// wait for postgres to be ready
			if (maxWait >= 0) {
				log.Println("wating for postgres, sleeping...")
				time.Sleep(100 * time.Millisecond)
				maxWait--
			} else {
				return nil
			}
		} else {
			return db
		}
	}
}

// waitForDb waits for database to be ready.  database launched from docker takes time to be ready
func waitForDb() {
	const SuccessCnt = 2
	successCntRemaining := SuccessCnt
	failedCntRemaining := 5
	for {
		db := newDbConnection()
		if db != nil {
			successCntRemaining--
			if successCntRemaining <= 0 {
				break
			}
			time.Sleep(2 * time.Second)
		} else {
			failedCntRemaining--
			if failedCntRemaining <= 0 {
				panic("database is not up")
			}
			successCntRemaining = SuccessCnt
			time.Sleep(2 * time.Second)
		}
	}
}

// CloseMainDbConn closes the main database connection
func CloseMainDbConn() {
	if dbConn == nil {
		return
	}
	err := dbConn.Close()
	if err != nil {
		log.Println("error closing db.", err)
	}
}

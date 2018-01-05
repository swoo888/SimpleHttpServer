package main

import (
	"fmt"
	"SpaceX/api"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Println("SpaceX golang starting")
	time.Sleep(8 * time.Second) //wait for postgres
	api.CreateSpaceBnBDb()

	//db := dbconn.NewDbConnection()
	//defer db.Close()
	//var lastInsertId int
	//err:= db.QueryRow("INSERT INTO booking(\"user\", title, description, expiration, " +
	//	"location_x, location_y) " +
	//	"VALUES($1,$2,$3,$4,$5,$6) returning id;",
	//	"user 1", "title 1", "description 1",
	//		"2012-12-09 13:10:23", 1.0, 2.0).Scan(&lastInsertId)
	//checkErr(err)
	//fmt.Println("last inserted id =", lastInsertId)

	r := mux.NewRouter()
	r.HandleFunc("/api/listings", api.Listings)
	r.HandleFunc("/api/listings/{id:[0-9]+}", api.SingleListing)
	log.Fatal(http.ListenAndServe(":9090", r))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
package main

import (
	"SimpleHttpServer/api"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"flag"
	"SimpleHttpServer/dbconn"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("SimpleHttpServer starting")
	flag.Parse()
	dbconn.WaitForDb()
	api.CreateSpaceBnBDb()

	r := mux.NewRouter()
	r.HandleFunc("/api/listings", api.Listings)
	r.HandleFunc("/api/listings/{id:[0-9]+}", api.SingleListing)
	log.Fatal(http.ListenAndServe(":9090", r))
}

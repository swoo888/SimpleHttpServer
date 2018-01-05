package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"flag"
	"SimpleHttpServer/dbconn"
	_ "github.com/lib/pq"
	"SimpleHttpServer/api/booking"
	bookingModel "SimpleHttpServer/model/booking"
)

func main() {
	log.Println("SimpleHttpServer starting")
	flag.Parse()
	_ = dbconn.GetMainDbConn()
	defer dbconn.CloseMainDbConn()
	bookingModel.CreateBookingTable()

	r := mux.NewRouter()
	r.HandleFunc("/api/listings", booking.Bookings)
	r.HandleFunc("/api/listings/{id:[0-9]+}", booking.Booking)
	log.Fatal(http.ListenAndServe(":9090", r))
}

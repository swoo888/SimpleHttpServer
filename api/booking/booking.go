package booking

import (
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"log"
	"SimpleHttpServer/model/booking"
	"SimpleHttpServer/util"
)

// Booking is handler function for url /api/listings/id"
func Booking(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	switch req.Method {
	case "GET":
		bookingGet(resp, id)
	case "DELETE":
		bookingDelete(resp, id)
	case "PUT":
		bookingPut(resp, req, id)
	default:
		log.Println("request method: ", req.Method)
	}
}

func bookingGet(resp http.ResponseWriter, id int64) {
	b, err := booking.GetBooking(id)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusNotFound)
		return
	}
	util.WriteJSONRespOrInternalServerError(resp, *b)
}

func bookingDelete(resp http.ResponseWriter, id int64) {
	res, err := booking.DeleteBooking(id)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	writeResultJSONResponse(resp, res)
}

func bookingPut(resp http.ResponseWriter, req *http.Request, id int64) {
	b := booking.Booking{}
	if err := util.DecodeJSONBodyOrBadRequest(resp, req.Body, &b); err != nil {
		return
	}
	res, err := b.UpdateBooking(id)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	writeResultJSONResponse(resp, res)
}

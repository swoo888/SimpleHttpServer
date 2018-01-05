package booking

import (
	"net/http"
	"log"
	"strconv"
	"database/sql"
	"SimpleHttpServer/util"
	"SimpleHttpServer/model/booking"
)

// Bookings is handler function for url /api/listings
func Bookings(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		bookingsGet(resp, req)
	case "DELETE":
		bookingsDelete(resp)
	case "POST":
		bookingsPost(resp, req)
	default:
		log.Println("request method: ", req.Method)
	}
}

func bookingsGet(resp http.ResponseWriter, req *http.Request) {
	active := req.URL.Query().Get("active")
	length := req.URL.Query().Get("length")
	page := req.URL.Query().Get("page")

	rows, err := booking.GetAllBookingRows(active)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	min, max, hasPagination := getPaginationMinMax(page, length)
	cur := 0
	listings := make([]booking.Booking, 0)
	for rows.Next() {
		if !hasPagination || (cur >= min && cur < max) {
			b, err := booking.GetBookingFromCurRow(rows)
			if err != nil {
				http.Error(resp, err.Error(), http.StatusInternalServerError)
				return
			}
			listings = append(listings, *b)
		}
		if hasPagination && cur >= max {
			break;
		}
		cur ++
	}
	util.WriteJSONRespOrInternalServerError(resp, listings)
}

func bookingsPost(resp http.ResponseWriter, req *http.Request) {
	var b booking.Booking
	if util.DecodeJSONBodyOrBadRequest(resp, req.Body, &b) != nil {
		return
	}
	lastInsertID, err := b.InsertBooking()
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	addedRec := struct{ ID int }{ID: lastInsertID}
	util.WriteJSONRespOrInternalServerError(resp, addedRec)
}

func getPaginationMinMax(page string, length string) (min int, max int, hasPagination bool) {
	min = -1
	max = -1
	hasPagination = false
	var err error
	if page == "" || length == "" {
		return
	}
	var pageNum int
	if pageNum, err = strconv.Atoi(page); err != nil {
		log.Println("Invalid page")
		return
	}
	var pageLen int
	if pageLen, err = strconv.Atoi(length); err != nil {
		log.Println("Invalid length")
		return
	}
	if pageNum <= 0 || pageLen <= 0 {
		log.Println("Invalid page or length")
		return
	}
	hasPagination = true
	min = (pageNum - 1) * pageLen
	max = (pageNum) * pageLen
	return
}

func bookingsDelete(resp http.ResponseWriter) {
	res, err := booking.DeleteAll()
	if err != nil {
		return
	}
	writeResultJSONResponse(resp, res)
}

func writeResultJSONResponse(resp http.ResponseWriter, res sql.Result) {
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	d := struct{ RowsAffected int64 }{RowsAffected: rowsAffected}
	util.WriteJSONRespOrInternalServerError(resp, d)
}

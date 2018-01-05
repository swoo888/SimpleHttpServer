package api

import (
	"net/http"
	"encoding/json"
	"log"
	"fmt"
	"SimpleHttpServer/dbconn"
	"strconv"
	"github.com/gorilla/mux"
	"database/sql"
	"time"
)

// Listings is handler function for url /api/listings
func Listings(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		listingsGet(resp, req)
	case "DELETE":
		listingsDelete(resp, req)
	case "POST":
		listingsPost(resp, req)
	default:
		log.Println("request method: ", req.Method)
	}
}

func listingsGet(resp http.ResponseWriter, req *http.Request) {
	active := req.URL.Query().Get("active")
	length := req.URL.Query().Get("length")
	page := req.URL.Query().Get("page")

	db := dbconn.NewDbConnection()
	defer dbconn.CloseDb(db)
	q := "SELECT * FROM booking "
	if active == "1" {
		q += "WHERE  expiration > now() "
	}
	q += "order by id "

	rows, err := db.Query(q)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	min, max, hasPagination := getPaginationMinMax(page, length)
	cur := 0
	listings := []spaceBnB{}
	for rows.Next() {
		if !hasPagination || (cur >= min && cur < max) {
			b := spaceBnB{}
			err = rows.Scan(&b.ID, &b.User, &b.Title, &b.Description, &b.Expiration,
				&b.Location.X, &b.Location.Y)
			if err != nil {
				http.Error(resp, err.Error(), http.StatusInternalServerError)
				return
			}
			listings = append(listings, b)
		}
		if hasPagination && cur >= max {
			break;
		}
		cur ++
	}
	writeJSONRespListings(listings, resp)
}

func writeJSONRespListings(listings []spaceBnB, resp http.ResponseWriter) {
	js, err := json.Marshal(listings)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	_, err = resp.Write(js)
	if err != nil {
		log.Println("unable to write json data to response.", err)
		return
	}
}

func closeReqBody(req *http.Request){
	err := req.Body.Close()
	if err != nil {
		log.Println("error closing request body.", err)
	}
}

func listingsPost(resp http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var b spaceBnB
	err := decoder.Decode(&b)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	defer closeReqBody(req)

	db := dbconn.NewDbConnection()
	defer dbconn.CloseDb(db)
	var lastInsertID int
	err = db.QueryRow("INSERT INTO booking(\"user\", title, description, expiration, "+
		"location_x, location_y) "+
		"VALUES($1,$2,$3,$4,$5,$6) returning id;",
		b.User, b.Title, b.Description,
		getDbExpiration(b.Expiration), b.Location.X, b.Location.Y).
		Scan(&lastInsertID)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("last inserted id =", lastInsertID)

	addedRec := struct{ ID int }{ID: lastInsertID}
	js, err := json.Marshal(addedRec)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSONResp(resp, js)
}

func getDbExpiration(expiration string) string {
	t, err := time.Parse("2006-01-02T15:04:05", expiration)
	if err != nil {
		return "now()"
	}
	return t.Format("2006-01-02 15:04:05")
}

func getPaginationMinMax(page string, length string) (min int, max int, hasPagination bool) {
	min = -1
	max = -1
	hasPagination = false
	var err error
	if page == "" || length == "" {
		log.Println("Invalid page or length")
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

func listingsDelete(resp http.ResponseWriter, req *http.Request) {
	db := dbconn.NewDbConnection()
	defer dbconn.CloseDb(db)

	log.Println("# Deleting")
	res, err := db.Exec("delete from booking")
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	writeResultJSONResponse(resp, res)
}

// SingleListing is handler function for url /api/listings/id"
func SingleListing(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	switch req.Method {
	case "GET":
		singleListingGet(resp, id)
	case "DELETE":
		singleListingDelete(resp, id)
	case "PUT":
		singleListingPut(resp, req, id)
	default:
		log.Println("request method: ", req.Method)
	}
}

func singleListingGet(resp http.ResponseWriter, id int) {
	db := dbconn.NewDbConnection()
	defer dbconn.CloseDb(db)
	q := "SELECT * FROM booking WHERE id=$1"
	stmt, err := db.Prepare(q)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	rows, err := stmt.Query(id)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	b := spaceBnB{}
	if rows.Next() {
		err = rows.Scan(&b.ID, &b.User, &b.Title, &b.Description, &b.Expiration,
			&b.Location.X, &b.Location.Y)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	js, err := json.Marshal(b)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSONResp(resp, js)
}

func singleListingDelete(resp http.ResponseWriter, id int) {
	db := dbconn.NewDbConnection()
	defer dbconn.CloseDb(db)
	q := "DELETE FROM booking WHERE id=$1"
	stmt, err := db.Prepare(q)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := stmt.Exec(id)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	writeResultJSONResponse(resp, res)
}

func singleListingPut(resp http.ResponseWriter, req *http.Request, id int) {
	db := dbconn.NewDbConnection()
	defer dbconn.CloseDb(db)
	decoder := json.NewDecoder(req.Body)
	var b spaceBnB
	err := decoder.Decode(&b)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	defer closeReqBody(req)
	q := "UPDATE booking set \"user\"=$1, title=$2, " +
		"description=$3, expiration=$4, " +
		"location_x=$5, location_y=$6 WHERE id=$7 "
	stmt, err := db.Prepare(q)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := stmt.Exec(b.User, b.Title, b.Description,
		getDbExpiration(b.Expiration),
		b.Location.X, b.Location.Y, id)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
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
	js, err := json.Marshal(d)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSONResp(resp, js)
}

func writeJSONResp(resp http.ResponseWriter, data []byte) {
	resp.Header().Set("Content-Type", "application/json")
	_, err:=resp.Write(data)
	if err != nil{
		log.Println("error writing json data.", err)
	}
}

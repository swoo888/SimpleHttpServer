package api

import (
	"net/http"
	"encoding/json"
	"log"
	"fmt"
	"SpaceX/dbconn"
	"strconv"
	"errors"
	"github.com/gorilla/mux"
	"database/sql"
	"time"
)

type Test struct {
	Name    string
	Dummies []string
}

func Listings(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		ListingsGet(resp, req)
	case "DELETE":
		ListingsDelete(resp, req)
	case "POST":
		ListingsPost(resp, req)
	default:
		log.Println("request method: ", req.Method)
	}
}

func ListingsGet(resp http.ResponseWriter, req *http.Request) {
	active := req.URL.Query().Get("active")
	length := req.URL.Query().Get("length")
	page := req.URL.Query().Get("page")

	db := dbconn.NewDbConnection()
	defer db.Close()
	fmt.Println("# Querying")
	q := "SELECT * FROM booking "
	if active == "1" {
		q += "WHERE  expiration > now() "
	}
	q += "order by id "

	rows, err := db.Query(q)
	checkErr(err)

	min := 0
	max := 0
	listings := []SpaceBnB{}
	hasPage := false
	pageNum, pageLen, err := getPageNumLen(page, length)
	if err == nil {
		hasPage = true
		min = (pageNum - 1) * pageLen
		max = (pageNum) * pageLen
	}
	cur := 0
	for rows.Next() {
		if !hasPage || (cur >= min && cur < max) {
			b := SpaceBnB{}
			err = rows.Scan(&b.ID, &b.User, &b.Title, &b.Description, &b.Expiration,
				&b.Location.X, &b.Location.Y)
			checkErr(err)
			fmt.Printf("%v\n", b)
			listings = append(listings, b)
		}
		if (hasPage && cur >= max) {
			break;
		}
		cur += 1
	}
	writeJsonRespListings(listings, resp)
	return
}

func writeJsonRespListings(listings []SpaceBnB, resp http.ResponseWriter) {
	js, err := json.Marshal(listings)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(js)
}

func ListingsPost(resp http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var b SpaceBnB
	err := decoder.Decode(&b)
	checkErr(err)
	defer req.Body.Close()

	db := dbconn.NewDbConnection()
	defer db.Close()
	var lastInsertId int
	err = db.QueryRow("INSERT INTO booking(\"user\", title, description, expiration, "+
		"location_x, location_y) "+
		"VALUES($1,$2,$3,$4,$5,$6) returning id;",
		b.User, b.Title, b.Description,
		getDbExpiration(b.Expiration), b.Location.X, b.Location.Y).
		Scan(&lastInsertId)
	checkErr(err)
	fmt.Println("last inserted id =", lastInsertId)

	addedRec := struct{ ID int }{ID: lastInsertId}
	js, err := json.Marshal(addedRec)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(js)
}

func getDbExpiration(expiration string) string {
	t, err := time.Parse("2006-01-02T15:04:05", expiration)
	if err != nil {
		return "now()"
	}
	return t.Format("2006-01-02 15:04:05")
}

func getPageNumLen(page string, length string) (int, int, error) {
	if page == "" || length == "" {
		return -1, -1, errors.New("Invalid page or length")
	}
	pageNum := 0
	var err error
	if pageNum, err = strconv.Atoi(page); err != nil {
		return -1, -1, errors.New("Invalid page")
	}
	pageLen := 0
	if pageLen, err = strconv.Atoi(length); err != nil {
		return -1, -1, errors.New("Invalid length")
	}
	if pageNum <= 0 || pageLen <= 0 {
		return -1, -1, errors.New("Invalid page or length")
	}
	return pageNum, pageLen, nil
}

func ListingsDelete(resp http.ResponseWriter, req *http.Request) {
	db := dbconn.NewDbConnection()
	defer db.Close()

	fmt.Println("# Deleting")
	res, err := db.Exec("delete from booking")
	checkErr(err)

	rowsAffected, err := res.RowsAffected()
	checkErr(err)
	d := struct{ RowsAffected int64 }{RowsAffected: rowsAffected}
	js, err := json.Marshal(d)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(js)
}

func SingleListing(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	log.Println("ID: ", id)
	db := dbconn.NewDbConnection()
	defer db.Close()
	switch req.Method {
	case "GET":
		q := "SELECT * FROM booking WHERE id=$1"
		stmt, err := db.Prepare(q)
		rows, err := stmt.Query(id)
		checkErr(err)
		b := SpaceBnB{}
		if rows.Next() {
			err = rows.Scan(&b.ID, &b.User, &b.Title, &b.Description, &b.Expiration,
				&b.Location.X, &b.Location.Y)
			checkErr(err)
		}

		js, err := json.Marshal(b)
		writeJsonResp(resp, js)
		return
	case "DELETE":
		q := "DELETE FROM booking WHERE id=$1"
		stmt, err := db.Prepare(q)
		checkErr(err)
		res, err := stmt.Exec(id)
		checkErr(err)
		writeResultJsonResponse(resp, res)
	case "PUT":
		decoder := json.NewDecoder(req.Body)
		var b SpaceBnB
		err := decoder.Decode(&b)
		checkErr(err)
		defer req.Body.Close()
		q := "UPDATE booking set \"user\"=$1, title=$2, " +
			"description=$3, expiration=$4, " +
			"location_x=$5, location_y=$6 WHERE id=$7 "
		stmt, err := db.Prepare(q)
		checkErr(err)
		res, err := stmt.Exec(b.User, b.Title, b.Description,
			getDbExpiration(b.Expiration),
			b.Location.X, b.Location.Y, id)
		checkErr(err)
		writeResultJsonResponse(resp, res)
	default:
		log.Println("request method: ", req.Method)
	}
}

func writeResultJsonResponse(resp http.ResponseWriter, res sql.Result) {
	rowsAffected, err := res.RowsAffected()
	checkErr(err)
	d := struct{ RowsAffected int64 }{RowsAffected: rowsAffected}
	js, err := json.Marshal(d)
	writeJsonResp(resp, js)
}

func writeJsonResp(resp http.ResponseWriter, data []byte) {
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(data)
}

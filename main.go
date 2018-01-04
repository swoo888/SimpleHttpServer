package main

import (
	"fmt"
	"SpaceX/dbconn"
	"time"
)

func main() {
	fmt.Println("SpaceX golang starting")
	//mux := http.NewServeMux()
	//hello := func(resp http.ResponseWriter, req *http.Request) {
	//	resp.Header().Add("Content-Type", "text/html")
	//	resp.WriteHeader(http.StatusOK)
	//	fmt.Fprint(resp, "Hello from Above!")
	//}
	//goodbye := func(resp http.ResponseWriter, req *http.Request) {
	//	resp.Header().Add("Content-Type", "text/html")
	//	resp.WriteHeader(http.StatusOK)
	//	fmt.Fprint(resp, "Goodbye, it's been real!")
	//}
	//mux.HandleFunc("/hello", hello)
	//mux.HandleFunc("/goodbye", goodbye)
	//http.ListenAndServe(":9080", mux)
	//test1()

	time.Sleep(1000*time.Second)
}

func test1() {
	createUserInfo := `
		CREATE TABLE userinfo
    	(
			id serial NOT NULL,
			username character varying(100) NOT NULL,
			departname character varying(500) NOT NULL,
			Created date,
			CONSTRAINT userinfo_pkey PRIMARY KEY (id)
		)`
	db := dbconn.NewDbConnection()
	defer db.Close()

	_, err:= db.Exec(createUserInfo)
	checkErr(err)
	var lastInsertId int
	err = db.QueryRow("INSERT INTO userinfo(username,departname,created) " +
		"VALUES($1,$2,$3) returning id;",
		"astaxie", "研发部门", "2012-12-09").Scan(&lastInsertId)
	checkErr(err)
	fmt.Println("last inserted id =", lastInsertId)

	fmt.Println("# Updating")
	stmt, err := db.Prepare("update userinfo set username=$1 where id=$2")
	checkErr(err)
	res, err := stmt.Exec("astaxieupdate", lastInsertId)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect, "rows changed")

	fmt.Println("# Querying")
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)

	for rows.Next() {
		var id int
		var username string
		var department string
		var created time.Time
		err = rows.Scan(&id, &username, &department, &created)
		checkErr(err)
		fmt.Println("id | username | department | created ")
		fmt.Printf("%3v | %8v | %6v | %6v\n", id, username, department, created)
	}

	fmt.Println("# Deleting")
	stmt, err = db.Prepare("delete from userinfo where id=$1")
	checkErr(err)

	res, err = stmt.Exec(lastInsertId)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect, "rows changed")

	time.Sleep(1000*time.Second)
}


func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
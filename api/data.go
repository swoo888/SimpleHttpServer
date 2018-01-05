package api

import (
	"SpaceX/dbconn"
)

type SpaceBnB struct {
	ID          int          `json:"-"`
	User        string       `json:"user"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Expiration  string       `json:"expiration"`
	Location    LocationType `json:"location"`
}

type LocationType struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

func NewDummySpaceBnB() SpaceBnB {
	return SpaceBnB{
		ID:          1,
		User:        "A User",
		Title:       "A Title",
		Description: "A Description",
		Expiration:  "2006-01-02T15:04:05",
		Location: LocationType{
			X: 2.1,
			Y: 3.1,
		},
	}
}

func CreateSpaceBnBDb() {
	stmtDrop := `
		DROP TABLE if EXISTS booking
		`
	stmt := `
		CREATE TABLE booking
    	(
			id serial NOT NULL,
			"user" character varying(120) NOT NULL,
			title character varying(140) NOT NULL,
			description character varying(1024),
			expiration timestamp,
			location_x real,
			location_y real,
			CONSTRAINT booking_pkey PRIMARY KEY (id)
		)`
	db := dbconn.NewDbConnection()
	defer db.Close()

	_, err := db.Exec(stmtDrop)
	checkErr(err)

	_, err = db.Exec(stmt)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

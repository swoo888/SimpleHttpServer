package api

import (
	"SimpleHttpServer/dbconn"
)

type spaceBnB struct {
	ID          int          `json:"-"`
	User        string       `json:"user"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Expiration  string       `json:"expiration"`
	Location    locationType `json:"location"`
}

type locationType struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

//func newDummySpaceBnB() spaceBnB {
//	return spaceBnB{
//		ID:          1,
//		User:        "A User",
//		Title:       "A Title",
//		Description: "A Description",
//		Expiration:  "2006-01-02T15:04:05",
//		Location: locationType{
//			X: 2.1,
//			Y: 3.1,
//		},
//	}
//}

// CreateSpaceBnBDb creates the database tables to be used by our application
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
	defer dbconn.CloseDb(db)

	_, err := db.Exec(stmtDrop)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(stmt)
	if err != nil {
		panic(err)
	}
}

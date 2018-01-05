package booking

import (
	"SimpleHttpServer/dbconn"
	"database/sql"
)

// Booking is the data structure for table Booking
// A simple booking demo table
type Booking struct {
	ID          int64        `json:"id,omitempty"`
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

const selectAll = `SELECT id, "user", title, description,
		to_char(expiration at time zone 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS'),
		location_x, location_y FROM booking `

// GetAllBookingRows returns all bookings
// if active is true, then only rows not expired is returned
func GetAllBookingRows(active string) (*sql.Rows, error) {
	q := selectAll
	if active == "1" {
		q += "WHERE  expiration > now() "
	}
	q += "order by id "

	rows, err := dbconn.GetMainDbConn().Query(q)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// GetBookingFromCurRow Scan current booking row into Booking data
func GetBookingFromCurRow(rows *sql.Rows) (*Booking, error) {
	b := Booking{}
	err := rows.Scan(&b.ID, &b.User, &b.Title, &b.Description, &b.Expiration,
		&b.Location.X, &b.Location.Y)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// DeleteAll deletes all rows in booking
func DeleteAll() (sql.Result, error) {
	res, err := dbconn.GetMainDbConn().Exec("delete from booking")
	if err != nil {
		return nil, err
	}
	return res, nil
}

// InsertBooking inserts current record into booking
func (b *Booking) InsertBooking() (lastInsertID int, err error) {
	lastInsertID = -1
	err = dbconn.GetMainDbConn().QueryRow(
		`INSERT INTO booking("user", title, description, expiration,
		location_x, location_y)
		VALUES($1,$2,$3,$4,$5,$6) returning id;`,
		b.User, b.Title, b.Description,
		b.Expiration, b.Location.X, b.Location.Y).Scan(&lastInsertID)
	return lastInsertID, err
}

// GetBooking returns the booking record for current booking id
func GetBooking(id int64) (*Booking, error) {
	q := selectAll + " WHERE id=$1"
	stmt, err := dbconn.GetMainDbConn().Prepare(q)
	if err != nil {
		return nil, err
	}
	b := Booking{}
	err = stmt.QueryRow(id).Scan(&b.ID, &b.User, &b.Title, &b.Description, &b.Expiration,
		&b.Location.X, &b.Location.Y)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// DeleteBooking deletes booking record with id
func DeleteBooking(id int64) (sql.Result, error) {
	q := "DELETE FROM booking WHERE id=$1"
	stmt, err := dbconn.GetMainDbConn().Prepare(q)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateBooking updates booking record with id with booking data
func (b *Booking) UpdateBooking(id int64) (sql.Result, error) {
	q := `UPDATE booking set "user"=$1, title=$2,
		description=$3, expiration=$4,
		location_x=$5, location_y=$6 WHERE id=$7`
	stmt, err := dbconn.GetMainDbConn().Prepare(q)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(b.User, b.Title, b.Description, b.Expiration,
		b.Location.X, b.Location.Y, id)
	if err != nil {
		return nil, err
	}
	return res, nil
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

// CreateBookingTable creates the database tables to be used by our application
func CreateBookingTable() {
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
	_, err := dbconn.GetMainDbConn().Exec(stmtDrop)
	if err != nil {
		panic(err)
	}

	_, err = dbconn.GetMainDbConn().Exec(stmt)
	if err != nil {
		panic(err)
	}
}

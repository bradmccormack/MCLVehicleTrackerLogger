package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	_ "github.com/mattn/go-sqlite3"
)

type GPSRecord struct {
        latitude  string
        longitude string
        message   string
        speed     int
        heading   float64
        fix       bool
        date      time.Time
        ID        string
}

var ip = flag.String("ip", "internal.myclublink.com.au", "ip address to send gps co-ordinates to")
var db = flag.String("database", "backend.db", "database to open gps records from")

//TODO - The Query should be concated as a single string. This enables us to easily send the sentence to the backend to replay easily.


var query = flag.String("query", "select * from GPSRecords", "query to obtain gps records -eg select * from GPSRecords")
var db *sql.DB

func main() {
	flag.Parse()
	var err err
	db, err = sql.Open("sqlite3", db)
	if err != nil {
		fmt.Printf("Failed to open %s", db)
		os.Exit(1)
	}

	defer db.Close

	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("Failed to execute query %s", query)
		os.Exit(1)
	}
	
	for rows.Next() {
		//TODO scan in the string
		//send the string over tcp
		//sleep a second		
	}



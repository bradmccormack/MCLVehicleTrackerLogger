package main

import (
	"database/sql"
	"net"
	"log"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"time"
)


var ip = flag.String("ip", "127.0.0.1:6969", "ip address to send gps co-ordinates to")
var dbname = flag.String("database", "backend.db", "database to open gps records from")
var query = flag.String("query", "select id, Latitude, Longitude, Speed, Heading, Fix, BusID from GPSRecords where ID > 1000", "query to obtain gps records -eg select * from GPSRecords")
var db *sql.DB

type GPS struct {
	ID int
	Message string
	Latitude string
	Longitude string
	Speed float64
	Heading float64
	Fix bool
	DateTime time.Time
	BusID string
}
func main() {
	flag.Parse()
	var err error
	db, err = sql.Open("sqlite3", *dbname)
	if err != nil {
		fmt.Printf("Failed to open %s", db)
		os.Exit(1)
	}

	rows, err := db.Query(*query)
	if err != nil {
		fmt.Printf("Failed to execute query %s", query)
		os.Exit(1)
	}

	var cords []GPS	

	for rows.Next() {
		var row GPS
		err = rows.Scan(&row.ID, &row.Latitude, &row.Longitude, &row.Speed, &row.Heading, &row.Fix, &row.BusID)		
		if(err != nil) {
			fmt.Printf("Error reading row - %s\n", err.Error())	
		}
		row.Message = ""
		row.DateTime = time.Now()
		cords = append(cords, row)
	}
	db.Close()
 
	fmt.Printf("Sending data to %s", *ip)
	//for now just spit across the latitude and longitude
	conn, err := net.Dial("tcp", *ip)
	if err != nil {
		log.Fatal("Cannot do tcp connection - %s", err.Error()) 
	} else {
		fmt.Printf("Connection made successfully \n")
	}



	var msg string
	for _, cord := range cords {
		var Fix string
		if(cord.Fix) { 
			Fix = "true" 
		} else { 
			Fix = "false"
		}

		//T signifies testing. The server will not log replayed co-ordinates
		msg = "T" + cord.Message + ",L"
                msg += cord.Latitude + ","
                msg += cord.Longitude + ","
                msg += "S" + fmt.Sprint(cord.Speed) + ","
		msg += "H" + fmt.Sprint(cord.Heading) + ","
                msg += "D" + cord.DateTime.Format(time.RFC3339) + ","
                msg += "F" + Fix + ","
                msg += cord.BusID

		//convert string to bytes.. send to server then wait a second and loop
		//bytes := []byte(msg)
		fmt.Printf("sentence is %s\n", msg)       
		fmt.Fprintf(conn, msg)
		time.Sleep(50 * time.Millisecond)

	
		
	}

}

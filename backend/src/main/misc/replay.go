package main

import (
	"database/sql"
	"net"
	"log"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"encoding/gob"
    "encoding/json"
	"os"
	"time"
)

type Response map[string]interface{}

var ip = flag.String("ip", "127.0.0.1:6969", "ip address to send gps co-ordinates to")
var dbname = flag.String("database", "backend.db", "database to open gps records from")
var query = flag.String("query", "select id, Latitude, Longitude, Speed, Heading, Fix, BusID from GPSRecords where ID > 1000", "query to obtain gps records -eg select * from GPSRecords")
var vid = flag.String("vid", "", "The id of the vehicle to use (override db value")

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

type GPSRecord struct {
	Latitude  string
	Longitude string
	Message   string
	Speed     float64
	Heading   float64
	Fix       bool
	Date      time.Time
	ID        string
}

type DiagnosticRecord struct {
	CPUTemp	float64
	CPUVolt float64
	CPUFreq float64
	MemFree	uint64
}



func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

func init() {
	gob.Register(GPS{})

}


func main() {
	flag.Parse()
	var err error
	db, err = sql.Open("sqlite3", *dbname)
	if err != nil {
		fmt.Printf("Failed to open %s", db)
		os.Exit(1)
	}

	fmt.Printf("Using Query %s\n", *query)

	rows, err := db.Query(*query)
	if err != nil {
		fmt.Printf("Failed to execute query (%s)\n", err.Error())
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

 	recCnt := len(cords)
 	if(recCnt == 0) {
 		fmt.Printf("Found no records that matched your query !\n")
 		os.Exit(0)
 	}

	fmt.Printf("Sending data to %s\n", *ip)
	conn, err := net.Dial("tcp", *ip)
	if err != nil {
		log.Fatal("Cannot do tcp connection - %s", err.Error()) 
	} else {
		fmt.Printf("Connection made successfully \n")
	}

	fmt.Printf("Found the following amount of records %d\n", recCnt)

	var msg,diag string
	for _, cord := range cords {
		var Fix string
		if(cord.Fix) { 
			Fix = "true" 
		} else { 
			Fix = "false"
		}

		var vehicleID string
		if(*vid != "") {
			vehicleID = *vid
		} else {
			vehicleID = cord.BusID
		}

		//T signifies testing. The server will not log replayed co-ordinates
				msg = "T" + cord.Message
                msg += "L" + cord.Latitude + "," + cord.Longitude + ","
                msg += "S" + fmt.Sprint(cord.Speed) + ","
				msg += "H" + fmt.Sprint(cord.Heading) + ","
                msg += "D" + cord.DateTime.Format(time.RFC3339) + ","
                msg += "F" + Fix + ","
                msg += "I" + vehicleID


		diag = "CT0.0,CV0.0,CF0.0,MF0.0"


		//send JSON over
		fmt.Printf("sentence is %s\n", msg)       

		fmt.Fprint(conn, Response{ "sentence" : msg, "diagnostics" : diag})
		time.Sleep(250 * time.Millisecond)

	
		
	}

}

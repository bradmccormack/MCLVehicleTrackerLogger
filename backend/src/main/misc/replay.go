package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"time"
)


var ip = flag.String("ip", "dev.myclublink.com.au", "ip address to send gps co-ordinates to")
var dbname = flag.String("database", "backend.db", "database to open gps records from")
var query = flag.String("query", "select * from GPSRecords", "query to obtain gps records -eg select * from GPSRecords")
var db *sql.DB

type GPS struct {
	ID int
	Message string
	Latitude string
	Longitude string
	Speed int
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
		err = rows.Scan(&row.ID, &row.Message, &row.Latitude, &row.Longitude, &row.Speed, &row.Heading, &row.Fix, &row.DateTime, &row.BusID)		
		if(err != nil) {
			fmt.Printf("Error reading row - %s", err.Error())	
		}
		cords = append(cords, row)
	}
	db.Close()

	//for now just spit across the latitude and longitude
	var msg string
	for _, cord := range cords {
		var Fix string
		if(cord.Fix) { 
			Fix = "true" 
		} else { 
			Fix = "false"
		}

		msg = cord.Message + ",L"
                msg += cord.Latitude + ","
                msg += cord.Longitude + ","
                msg += "S" + string(cord.Speed) + ","
		msg += "H" + fmt.Sprint(cord.Heading, ',')
                msg += "D" + cord.DateTime.Format(time.RFC3339) + ","
                msg += "F" + Fix + ","
                msg += cord.BusID;
                //Byte[] sendBytes = Encoding.ASCII.GetBytes(data);
                /*
                Byte[] sendBytes = Encoding.ASCII.GetBytes("PHi there buddy,L" + lat +
                ",150.81094,S0.00,H147.2,D" + dt + ",Ftrue,BRADSBUS");
                */

                //udpClient.Send(sendBytes, sendBytes.Length);
                //System.Threading.Thread.Sleep(1000);
       
		
	}

}

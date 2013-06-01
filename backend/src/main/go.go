package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
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

func main() {

	db, err := sql.Open("sqlite3", "./backend.db")
	if err != nil {
		fmt.Printf("Cannot open backend.db . Exiting")
		os.Exit(1)
	}
	defer db.Close()

	service := ":215"
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	if err != nil {
		fmt.Printf("Failed to resolve UDP address")
		os.Exit(1)
	}
	fmt.Printf("Listening on UDP Port %s", service)

	for {
		con, err := net.ListenUDP("udp", udpAddr)
		if err != nil {
			fmt.Printf("Failed to create udp connection - %s", err)
			os.Exit(1)
		}
		handleClient(db, con)
	}
}

func notifyHTTP(entry *GPSRecord) {
	//http://html5labs.interoperabilitybridges.com/prototypes/websockets/websockets/info    -- WOOT it's possible to use web sockets on the client
}

func logEntry(db *sql.DB, entry *GPSRecord) {

	_, err := db.Exec("INSERT INTO GPSRecords (Message, Latitude, Longitude, Speed, Heading, Fix, DateTime, BusID) VALUES ( ? , ?, ? , ? , ? ,? ,? , ?)",
		entry.message,
		entry.latitude,
		entry.longitude,
		entry.speed,
		entry.heading,
		entry.fix,
		entry.date,
		entry.ID)

	if err != nil {
		fmt.Printf("Failed to insert row %s", err)
	}

	//daytime := time.Now().String()
}

//palm off reading and writing to a go routine
func handleClient(db *sql.DB, conn *net.UDPConn) {
	defer conn.Close()
	var buff [512]byte
	var entry GPSRecord

	n, addr, err := conn.ReadFromUDP(buff[:])
	if err != nil {
		fmt.Printf("Error reading from UDP")
	}

	gpsfields := strings.Split(string(buff[:n]), ",")
	if len(gpsfields) != 8 {
		fmt.Printf("Error. GPS fields length is incorrect. Is %d should be %d", len(gpsfields), 8)
		fmt.Printf("The source string was %s\n", string(buff[:n]))
		os.Exit(1)
	}
	//All data is validated on the logger end so I'm going to assume for now that Parsing will be fine. Perhaps a network error could occur and I'll fix that up later

	entry.message = gpsfields[0][1:]
	entry.latitude = gpsfields[1][1:]
	entry.longitude = gpsfields[2]
	entry.speed, _ = strconv.Atoi(gpsfields[3][1:])
	entry.heading, _ = strconv.ParseFloat(gpsfields[4][1:], 32)
	entry.date, _ = time.Parse(time.RFC3339, gpsfields[5][1:]) //todo pull out just the date component and format
	entry.fix = gpsfields[6][1:] == "true"
	entry.ID = gpsfields[7]

	fmt.Printf("Message %s Lat %s Long %s speed %d heading %f fix %t date %s time %s id %s\n",
		entry.message,
		entry.latitude,
		entry.longitude,
		entry.speed,
		entry.heading,
		entry.fix,
		entry.date,
		entry.ID)

	go logEntry(db, &entry) //save to database
	go notifyHTTP(&entry)   //notify any HTTP observers

	conn.WriteToUDP([]byte("OK"), addr)
	fmt.Printf("Responded to %s\n", addr)
	conn.Close()

}

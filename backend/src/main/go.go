package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/garyburd/go-websocket/websocket"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net"
	"net/http"
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

var service = flag.String("service", ":6969", "udp port to bind to")
var addr = flag.String("addr", ":8080", "http(s)) service address")

var con *websocket.Conn //fix me up this is dirty (use channels and stuff later)

func handleHTTP() {
	fmt.Printf("Listening for HTTP on %s\n", *addr)
	//go h.run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			http.Error(w, "Method not allowed", 405)
			return
		}

		/*
			if r.Header.Get("Origin") != "http://"+r.Host {
				http.Error(w, "Origin not allowed", 403)
				return
			}
		*/

		var err error
		con, err = websocket.Upgrade(w, r.Header, nil, 1024, 1024)
		if _, ok := err.(websocket.HandshakeError); ok {
			http.Error(w, "Not a websocket handshake", 400)
			return
		} else if err != nil {

			log.Println(err)
			return
		}

		//defer con.Close()

		//fmt.Printf("Received a request via web socket route. YAY")
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		fmt.Printf("Failed to listen for http on %s", *addr)
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {

	flag.Parse()

	db, err := sql.Open("sqlite3", "./backend.db")
	if err != nil {
		fmt.Printf("Cannot open backend.db . Exiting")
		os.Exit(1)
	}
	defer db.Close()

	udpAddr, err := net.ResolveUDPAddr("udp4", *service)
	if err != nil {
		fmt.Printf("Failed to resolve UDP address")
		os.Exit(1)
	}

	fmt.Printf("Listening on UDP Port %s\n", *service)

	go handleHTTP()

	for {
		con, err := net.ListenUDP("udp", udpAddr)
		if err != nil {
			fmt.Printf("Failed to create udp connection - %s", err)
			os.Exit(1)
		}
		handleClient(db, con)
	}
}

func updateClient(entry *GPSRecord) {
	if con != nil {
		wswriter, _ := con.NextWriter(websocket.OpText)

		io.WriteString(wswriter, string(entry.latitude+","+entry.longitude)) //we want to write some JSON instead of text for now just do a dodgy string
		//http://html5labs.interoperabilitybridges.com/prototypes/websockets/websockets/info    -- WOOT it's possible to use web sockets on the client
	} else {
		fmt.Printf("Web socket is closed")
	}
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

	fmt.Printf("The date that I was sent was %s\n", gpsfields[5][1:])

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
	updateClient(&entry)    //notify any HTTP observers //make this a goroutine later

	conn.WriteToUDP([]byte("OK"), addr)
	fmt.Printf("Responded to %s\n", addr)
	conn.Close()

}

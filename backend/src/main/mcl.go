package main

import (
	"./actions"
	"./views"
	"database/sql"
	"flag"
	"fmt"
	"github.com/garyburd/go-websocket/websocket"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"io"
	"log"
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
var connections []*websocket.Conn //slice of Websocket connections
var Db *sql.DB

func handleWebSocketInit(w http.ResponseWriter, r *http.Request) {
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
	var (
		connection *websocket.Conn
		err        error
	)

	connection, err = websocket.Upgrade(w, r.Header, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	connections = append(connections, connection)
	fmt.Printf("Amount of clients listening is %d", len(connections))
}

//TODO look at implementing trinity mvc framework
func handleHTTP() {

	Router := mux.NewRouter()

	viewRouter := Router.Methods("GET").Subrouter()
	actionRouter := Router.Methods("POST").Subrouter()

	//Handle web socket traffic specially
	Router.HandleFunc("/ws", handleWebSocketInit)

	//TODO - Look at moving non-websocket traffic to fastcgi protocol

	//View Routes
	viewRouter.HandleFunc("/settings", ViewSettings)
	viewRouter.HandleFunc("/settings")
	viewRouter.HandleFunc("/", func(w http.ResponseWriter, request *http.Request) {
		http.Error(w, "Invalid view", 403)
	})

	//Action Routes
	actionRouter.HandleFunc("/login", ActionLogin)
	actionRouter.HandleFunc("/settings", ActionSettings)
	actionRouter.HandleFunc("/", func(w http.ResponseWriter, request *http.Request) {
		http.Error(w, "Invalid action", 403)
	})

	//Use the router
	http.Handle("/", Router)

	fmt.Printf("Listening for HTTP on %s\n", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		fmt.Printf("Failed to listen for http on %s", *addr)
		log.Fatal("ListenAndServe: ", err)
	}

}

func main() {

	flag.Parse()

	var err error

	Db, err = sql.Open("sqlite3", "./backend.db")
	if err != nil {
		fmt.Printf("Cannot open backend.db . Exiting")
		os.Exit(1)
	}
	defer Db.Close()

	udpAddr, err := net.ResolveUDPAddr("udp4", *service)
	if err != nil {
		fmt.Printf("Failed to resolve UDP address")
		os.Exit(1)
	}
	fmt.Printf("Listening on UDP Port %s\n", *service)

	//handle web requests in a seperate go-routine
	go handleHTTP()

	//wait around for UDP requests and handle them when they come in
	for {
		udpcon, err := net.ListenUDP("udp", udpAddr)
		if err != nil {
			fmt.Printf("Failed to create udp connection - %s", err)
			os.Exit(1)
		}
		handleClient(Db, udpcon)
	}
}

func updateClient(entry *GPSRecord) {

	if connections == nil {
		fmt.Printf("No clients listening.. not reporting")
		return
	}

	fmt.Printf("Responding to %d listening clients\n", len(connections))
	for _, client := range connections {
		//get a websocket writer
		wswriter, _ := client.NextWriter(websocket.OpText)

		if wswriter != nil {
			io.WriteString(wswriter, string(entry.latitude+","+entry.longitude)) //we want to write some JSON instead of text for now just do a dodgy string
		} else {
			fmt.Printf("No ws writer available\n") //this web socket was abruptly closed so we need to close that client and remove it from the connections slice
			client.Close()
			//connections[index] = nil

		}

	}
}

func logEntry(entry *GPSRecord) {

	_, err := Db.Exec("INSERT INTO GPSRecords (Message, Latitude, Longitude, Speed, Heading, Fix, DateTime, BusID) VALUES ( ? , ?, ? , ? , ? ,? ,? , ?)",
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
func handleClient(Db *sql.DB, conn *net.UDPConn) {

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

	go logEntry(&entry)  //save to database
	updateClient(&entry) //notify any HTTP observers //make this a goroutine later

	conn.WriteToUDP([]byte("OK"), addr)

}

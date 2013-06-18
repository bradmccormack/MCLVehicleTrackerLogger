package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/garyburd/go-websocket/websocket"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
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

var service = flag.String("service", ":6969", "tcp port to bind to")
var addr = flag.String("addr", ":8080", "http(s)) service address")
var connections []*websocket.Conn //slice of Websocket connections
var db *sql.DB

func ActionSettings(w http.ResponseWriter, r *http.Request) {

}

func ViewSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	var err error
	t := template.New("Settings")
	t, err = template.ParseFiles("templates/settings.html")
	if err != nil {
		fmt.Printf("Failed to parse the template file!\n")
		return
	}

	userID := 1 //this should come from the request form

	row := db.QueryRow("SELECT S.MapAPI, U.FirstName, U.LastName FROM Settings S, User U WHERE S.UserID = ?", userID)

	var settings = map[string]string{
		"MapAPI":    "",
		"FirstName": "",
		"LastName":  "",
	}

	var MapAPI, FirstName, LastName string
	row.Scan(&MapAPI, &FirstName, &LastName)

	settings["MapAPI"] = MapAPI
	settings["FirstName"] = FirstName
	settings["LastName"] = LastName

	t.Execute(w, settings) //second param is the data interface. It's the equiv of Bondi's AddDataSet I believe'

}

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

	Router.HandleFunc("/ws", handleWebSocketInit)
	Router.HandleFunc("/settings", ViewSettings)

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

	db, err = sql.Open("sqlite3", "./backend.db")
	if err != nil {
		fmt.Printf("Cannot open backend.db . Exiting")
		os.Exit(1)
	}
	defer db.Close()

	//tcpAddr, err := net.ResolveTCPAddr("tcp", *service)
	lnk, err := net.Listen("tcp", *service)
	if err != nil {
		fmt.Printf("Failed to get tcp listener")
		os.Exit(1)
	}
	fmt.Printf("Listening on TCP Port %s\n", *service)

	//handle web requests in a seperate go-routine
	go handleHTTP()

	//wait around for tcp requests and handle them when they come in
	for {
		//tcpcon, err := net.ListenTCP("tcp", tcpAddr)
		tcpcon, err := lnk.Accept()
		if err != nil {
			fmt.Printf("Failed to create tcp connection - %s", err)
			os.Exit(1)
		}
		//note to self, the part after tcpcon. is called type assertion. TODO find out how it relates to casting in other languages
		handleClient(db, tcpcon.(*net.TCPConn))
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
func handleClient(db *sql.DB, conn *net.TCPConn) {

	defer conn.Close()
	var buff [512]byte
	var entry GPSRecord

	n, err := conn.Read(buff[:])
	if err != nil {
		fmt.Printf("Error reading from TCP")
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

	conn.Write([]byte("OK"))

}

package main

import (
	"database/sql"
	"encoding/json"
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

type User struct {
	id          int
	firstname   string
	lastname    string
	password    string
	company     string
	accesslevel int
}

type Session struct {
	user User
}

var service = flag.String("service", ":6969", "tcp port to bind to")
var addr = flag.String("addr", ":8080", "http(s)) service address")
var connections []*websocket.Conn //slice of Websocket connections
var Db *sql.DB

var actions = map[string]interface{}{
	"ActionInvalid": func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid Action", 403)
	},

	"ActionLogin": func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		var user User
		name := r.FormValue("name")
		password := r.FormValue("password")

		err := Db.QueryRow("SELECT ID, U.FirstName, U.LastName, U.AccessLevel, C.Name, C.MaxUsers, C.Expiry FROM User U, Company C WHERE U.FirstName = ? AND U.Password = ? AND C.ID = U.CompanyID", name, password)

		switch {
		case err == sql.ErrNoRows:
			log.Printf("No user with that ID.")
			//decrease the retries and send failure back with reduced retries
		case err != nil:
			log.Fatal(err)
		default:
			//read the row and make a session cookie and send it back with some JSON
			row.Scan(&user.id)

		}

		json, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		fmt.Printf("In Action Login")
	},

	"ActionSettings": func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Printf("In Action Settings")
	},
}

//Note - Template caching needs to be implemented http://golang.org/doc/articles/wiki/ There is an inefficiency in this code: renderTemplate calls ParseFiles every time a page is rendered.
var views = map[string]interface{}{

	"ViewInvalid": func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid view", 403)
	},

	"ViewLogin": func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("In ViewLogin")
		var err error
		t := template.New("Login")
		t, err = template.ParseFiles("templates/login.html")
		if err != nil {
			fmt.Printf("Failed to parse the template file!\n")
			return
		}
		t.Execute(w, nil)
	},

	"ViewSettings": func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("In ViewSettings")
		w.Header().Add("Content-Type", "text/html")

		var err error
		t := template.New("Settings")
		t, err = template.ParseFiles("templates/settings.html")
		if err != nil {
			fmt.Printf("Failed to parse the template file!\n")
			return
		}

		userID := 1 //this should come from the request form

		row := Db.QueryRow("SELECT S.MapAPI, U.FirstName, U.LastName FROM Settings S, User U WHERE S.UserID = ?", userID)

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

		t.Execute(w, settings)
	},
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

	viewRouter := Router.Methods("GET").Subrouter()
	actionRouter := Router.Methods("POST").Subrouter()

	//Handle web socket traffic specially
	Router.HandleFunc("/ws", handleWebSocketInit)

	//TODO - Look at moving non-websocket traffic to fastcgi protocol

	//View Routes
	viewRouter.HandleFunc("/system/settings", views["ViewSettings"].(func(http.ResponseWriter, *http.Request)))
	viewRouter.HandleFunc("/system/login", views["ViewLogin"].(func(http.ResponseWriter, *http.Request)))
	viewRouter.HandleFunc("/", views["ViewInvalid"].(func(http.ResponseWriter, *http.Request)))

	//Action Routes
	actionRouter.HandleFunc("/system/login", actions["ActionLogin"].(func(http.ResponseWriter, *http.Request)))
	actionRouter.HandleFunc("/system/settings", actions["ActionSettings"].(func(http.ResponseWriter, *http.Request)))
	actionRouter.HandleFunc("/", actions["ActionInvalid"].(func(http.ResponseWriter, *http.Request)))

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
		handleClient(Db, tcpcon.(*net.TCPConn))
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
func handleClient(Db *sql.DB, conn *net.TCPConn) {

	defer conn.Close()
	var buff [512]byte
	var entry GPSRecord

	n, err := conn.Read(buff[:])
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

	conn.Write([]byte("OK"))

}

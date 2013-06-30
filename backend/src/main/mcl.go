package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/garyburd/go-websocket/websocket"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
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

type Company struct {
	Name     string
	Maxusers int
	Expiry   string
}

type User struct {
	ID          int
	Firstname   string
	Lastname    string
	Password    string
	Accesslevel int
}

type Session struct {
	User    *User
	Company *Company
}

type Response map[string]interface{}

//set the domain based upon the path the executable was run from
var domain string = "dev.myclublink.com.au"

var service = flag.String("service", ":6969", "tcp port to bind to")
var addr = flag.String("addr", ":8080", "http(s)) service address")
var connections []*websocket.Conn //slice of Websocket connections

//for generating secure cookies
var hashKey = securecookie.GenerateRandomKey(32)
var blockKey = securecookie.GenerateRandomKey(32)
var s = securecookie.New(hashKey, nil) //don't supply the blockkey for now and not encrypt the data. hashkey is used to identify, block key is used to hash
var Db *sql.DB

/*TODO make the following a map and use a cookiejar too to keep track of cookies per user */
var user User
var company Company
var session Session


func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

var actions = map[string]interface{}{
	"ActionInvalid": func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid Action", 403)
	},

	"ActionLogin": func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		//var user User
		//var company Company
		//var session Session

		name := r.FormValue("name")
		password := r.FormValue("password")

		if Db == nil {
			log.Fatal(Db)
		}

		result := Db.QueryRow("SELECT U.ID, U.FirstName, U.LastName, U.AccessLevel, C.Name, C.MaxUsers, C.Expiry FROM User U, Company C WHERE U.FirstName = ? AND U.Password = ? AND C.ID = U.CompanyID",
			name, password).Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Accesslevel, &company.Name, &company.Maxusers, &company.Expiry)

		switch {
		case result == sql.ErrNoRows:
			fmt.Fprint(w, Response{"success": false, "message": "IncorrectLogin", "retries": 0})
			return
		case result != nil:
			log.Fatal(result)
		default:
			session.User = &user
			session.Company = &company

			//TODO think about multiple users and cookie jar implementation?
			encoded, err := s.Encode("Session", session)
			if err != nil {
				fmt.Fprint(w, Response{"success": false, "message": "Failed to create Session cookie"})
				return
			}

			expire := time.Now().AddDate(0, 0, 1)
			cookie := http.Cookie{Name: "Session", Value: encoded, Path: "/",
				Domain: domain, Expires: expire, RawExpires: expire.Format(time.UnixDate), MaxAge: 86400, Secure: true, HttpOnly: true}
			http.SetCookie(w, &cookie)
			fmt.Fprint(w, Response{"success": true, "message": "All good", "session": session})
		}

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
		var err error
		t := template.New("Login")
		t, err = template.ParseFiles("templates/login.html")
		if err != nil {
			fmt.Printf("Failed to parse the template file!\n")
			return
		}

		var LoginInfo = map[string]bool{
			"LoggedOut": false,
		}
		t.Execute(w, LoginInfo)
	},
	"ViewLicense": func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")

		var err error
		t := template.New("License")
		t, err = template.ParseFiles("templates/license.html")
		if err != nil {
			log.Fatal("Failed to read the template file for license. Fix it")
		}
		//TODO fix this fucker
		fmt.Printf("Session user is %s", session.User.Firstname)
		t.Execute(w, session)

	},

	"ViewSettings": func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")

		var err error
		t := template.New("Settings")
		t, err = template.ParseFiles("templates/settings.html")
		if err != nil {
			log.Fatal("Failed to parse the template file for settings. Fix it")
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

func createDb() {

	file, err := os.Open("./backend.db")
	if err != nil {
		fmt.Printf("./backend.db didn't exist. Creating it ! \n")
	} else {
		file.Close()
		return
	}

	fHandle, err := os.Create("./backend.db")
	if err != nil {
		log.Fatal("Cannot create ./backend.db ! Bailing from running server\n")
	}
	fHandle.Close()

	//TODO add indexes
	statements := []string{

		"BEGIN EXCLUSIVE TRANSACTION;",

		//Use a string array of raw string literals

		`CREATE TABLE GPSRecords (
         id integer primary key autoincrement, Message text,
         Latitude text not null,
         Longitude text not null,
         Speed integer not null,
         Heading float not null,
         Fix boolean not null,
         DateTime date not null default current_timestamp,
        BusID text not null);`,

		`CREATE TABLE Errors (
        id integer primary key autoincrement,
        GPSRecordID integer not null,
        Error text,
        DateTime date not null default current_timestamp,
        FOREIGN KEY (GPSRecordID) REFERENCES GPSrecords(id)
	);`,

		`CREATE TABLE Network (
        id integer primary key autoincrement,
        GPSRecordID integer not null,
        Acknowledge boolean not null default 0,
        FOREIGN KEY (GPSRecordID) REFERENCES GPSRecords(id)
	);`,

		`CREATE TABLE Company (
        ID integer primary key autoincrement,
        Name text not null,
        Expiry date not null default current_timestamp,
        MaxUsers integer not null default 0
	);`,

		`CREATE TABLE User (
        ID integer primary key autoincrement,
        FirstName text not null,
        LastName text not null,
        CompanyID integer not null,
        Password text not null,
        AccessLevel integer not null default 0,
        FOREIGN KEY (CompanyID) REFERENCES Company(ID)
	);`,

		`CREATE TABLE Settings (
        ID integer primary key autoincrement,
        UserID integer not null,
        MapAPI text not null default 'GoogleMaps',
        FOREIGN KEY (UserID) REFERENCES User(ID)
	);`,

		"INSERT INTO Company (Name, MaxUsers) VALUES ('myClubLink' , 1);",
		"INSERT INTO User (FirstName, LastName, CompanyID, Password, AccessLevel) VALUES ('guest','user', 1, 'guest', 0);",
		"INSERT INTO Settings (UserID, MapAPI) VALUES (1, 'GoogleMaps');",
		"COMMIT TRANSACTION;",
	}
	Db, err = sql.Open("sqlite3", "./backend.db")

	for _, statement := range statements {
		_, err := Db.Exec(statement)
		if err != nil {
			log.Fatal(err)
		}
	}

	Db.Close()
	fmt.Printf("Finished creating\n")
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
	viewRouter.HandleFunc("/system/license", views["ViewLicense"].(func(http.ResponseWriter, *http.Request)))
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

	createDb()

	Db, err = sql.Open("sqlite3", "./backend.db")
	if err != nil {
		fmt.Printf("Cannot open backend.db . Exiting")
		//os.Exit(1)
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
	var buff = make([]byte, 512)
	var entry GPSRecord

	conn.SetReadBuffer(512)

	var n int
	var err error
	
	var fuck = 69
	for(fuck == 69) {
		n, err = conn.Read(buff)
		if err != nil {
			fmt.Printf("Error reading from UDP")
		}
		fmt.Printf("Read %s",string(buff[:n]))
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

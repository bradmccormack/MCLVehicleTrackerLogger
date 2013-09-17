package main

import (
//	"bytes"
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/garyburd/go-websocket/websocket"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

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
	Email	    string
}

type Settings struct {
	MapAPI 	    string
	Interpolate	int
	SnaptoRoad	int
	CameraPanTrigger int
	RadioCommunication int
	DataCommunication int
	SecurityRemoteAdmin int
	SecurityConsoleAccess int
	SecurityAdminPasswordReset int
	MobileSmartPhoneAccess int
	MobileShowBusLocation int
}


type Response map[string]interface{}

//set the domain based upon the path the executable was run from
var domain string = "dev.myclublink.com.au"

var service = flag.String("service", ":6969", "tcp port to bind to")
var addr = flag.String("addr", ":8080", "http(s)) service address")
var connections []*websocket.Conn                                       //slice of Websocket connections
var random *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano())) //new random with unix time nano seconds as seed
//Session information
var store = sessions.NewCookieStore([]byte("emtec789"))
var Db *sql.DB

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

		var user User
		var company Company
		var settings Settings

		name := r.FormValue("name")
		password := r.FormValue("password")

		if Db == nil {
			log.Fatal(Db)
		}

		result := Db.QueryRow(`
			SELECT U.ID, U.FirstName, U.LastName, U.Password, U.AccessLevel, U.Email, C.Name, C.MaxUsers, C.Expiry, S.MapAPI, S.Interpolate, S.SnaptoRoad, S.CameraPanTrigger,
			S.RadioCommunication, S.3GCommunication, S.SecurityRemoteAdmin, S.SecurityConsoleAccess, S.SecurityAdminPasswordReset, S.MobileSmartPhoneAccess, S.MobileShowBusLocation 
			FROM User U, Company C, Settings S
			WHERE UPPER(U.FirstName) = ? AND U.Password = ? AND C.ID = U.CompanyID AND S.UserID = U.ID`,
			strings.ToUpper(name), password).Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Accesslevel, &user.Email, &company.Name, &company.Maxusers, &company.Expiry,
			&settings.MapAPI, &settings.Interpolate, &settings.SnaptoRoad, &settings.CameraPanTrigger, &settings.RadioCommunication, &settings.DataCommunication, &settings.SecurityRemoteAdmin,
			&settings.SecurityConsoleAccess, &settings.SecurityAdminPasswordReset, &settings.MobileSmartPhoneAccess, &settings.MobileShowBusLocation)

		switch {
		case result == sql.ErrNoRows:
			fmt.Fprint(w, Response{"success": false, "message": "Incorrect Username or Password", "retries": 10})
			return
		case result != nil:
			log.Fatal(result)
		default:

			//TODO check expiry of license
			//TODO check amount of loggined in users



			session, _ := store.Get(r, "session")

			session.Values["User"] = user
			session.Values["Company"] = company
			session.Options = &sessions.Options{
				Path:   "/",
				MaxAge: 86400, //1 day
			}

			if err := session.Save(r, w); err != nil {
				fmt.Printf("Can't save session data (%s)\n", err.Error())

			}
			fmt.Fprint(w, Response{"success": true, "message": "Login OK", "user": user.Firstname + " " + user.Lastname})
		}

	},

	"ActionSettings": func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Printf("In Action Settings")
	},
	"ActionHistorialRoute": func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "application/json")
		//map with string key and slice of string slices
		var Route = make(map[string][][]string)

		dateFrom := r.FormValue("dateFrom")
		dateTo := r.FormValue("dateTo")

		rows, err := Db.Query("SELECT BusID, Latitude, Longitude, Speed, Heading, Fix, DateTime FROM GPSRecords WHERE datetime >=? AND datetime <=? GROUP BY id ORDER BY datetime asc", dateFrom, dateTo)
		if err != nil {
			log.Fatal(err)
		}

		var ID, Lat, Long, Speed, Fix, Heading, Date string

		//build up the map here
		for rows.Next() {
			if err := rows.Scan(&ID, &Lat, &Long, &Speed, &Heading, &Fix, &Date); err != nil {
				log.Fatal(err)
			}
			Route[ID] = append(Route[ID], []string{Lat, Long, Speed, Fix, Heading, Date})
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		fmt.Fprint(w, Response{"success": true, "data": Route})

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
	"ViewSupport": func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		session, err := store.Get(r, "session")
		if session == nil {
			fmt.Printf("Session is nil \n")
		}

		if err != nil {
			fmt.Printf("Error loading session information %s", err.Error())
		}
		t, err := template.ParseFiles("templates/support.html")
		if err != nil {
			log.Fatal("Failed to read the template file for support. Fix it")
		}
		t.Execute(w, session.Values)
	},
	"ViewReport": func(w http.ResponseWriter, r *http.Request) {
		//TODO redirect to root if not logged in

		w.Header().Add("Content-Type", "application/json")

		//session, _ := store.Get(r, "session")

		var percentAvailable int = random.Intn(75)
		availability := [...]int{percentAvailable, 100 - percentAvailable}

		var kmPerDay [7]int
		for i := 0; i < 7; i++ {
			kmPerDay[i] = random.Intn(60)
		}

		fmt.Fprint(w, Response{"Availability": availability, "KMPerDay": kmPerDay})
		//t.Execute(w, session.Values)
	},

	"ViewMap": func(w http.ResponseWriter, r *http.Request) {
		//TODO redirect to root if not logged in
		w.Header().Add("Content-Type", "text/html")

		session, _ := store.Get(r, "session")
		fmt.Printf("Session is %s", Response{"session": session})

		var err error
		t := template.New("Map")
		t, err = template.ParseFiles("templates/map.html")
		if err != nil {
			log.Fatal("Failed to read the template file for map. Fix it")
		}
		t.Execute(w, session.Values)
	},

	"ViewLicense": func(w http.ResponseWriter, r *http.Request) {
		//TODO redirect to root if not logged in
		w.Header().Add("Content-Type", "text/html")

		session, _ := store.Get(r, "session")

		var err error
		t := template.New("License")
		t, err = template.ParseFiles("templates/license.html")
		if err != nil {
			log.Fatal("Failed to read the template file for license. Fix it")
		}
		t.Execute(w, session.Values)
	},

	"ViewSettings": func(w http.ResponseWriter, r *http.Request) {
		//TODO redirect to root if not logged in
		w.Header().Add("Content-Type", "text/html")

		session, _ := store.Get(r, "session")

		var mapAPI string
		var interpolate, snaptoroad bool
		var user User = session.Values["User"].(User)

		result := Db.QueryRow(`
                        SELECT S.MapAPI, S.Interpolate, S.SnaptoRoad
                        FROM Settings S, User U
			WHERE S.UserID = U.ID 
			AND U.ID = ?`, user.ID).Scan(&mapAPI, &interpolate, &snaptoroad)

		switch {
		case result != nil:
			log.Fatal(result)
		default:
			session, _ := store.Get(r, "session")

			//TODO add the receive data settings here and get from the DB
			session.Values["Settings"] = map[string]interface{}{
				"MapAPI":      mapAPI,
				"Interpolate": interpolate,
				"SnaptoRoad":  snaptoroad,
			}
			session.Save(r, w)
		}

		var err error
		t := template.New("Settings")
		t, err = template.ParseFiles("templates/settings.html")
		if err != nil {
			fmt.Printf(err.Error())
			log.Fatal("\nFailed to parse the template file for settings. Fix it")
		}

		/*TODO change accesslevel to text, Guest/Admin etc so it is more friendly */
		t.Execute(w, session.Values)
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
         Fix integer not null,
         DateTime date not null default current_timestamp,
        BusID text not null);`,
		
	`CREATE TABLE Support (
	SupportID integer primary key autoincrement,
	UserID integer not null,
	Subject text not null,
	Body text not null,
	DateTime date not null default current_timestamp,
	FOREIGN KEY (UserID) REFERENCES User(ID)
	);`,


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
        Acknowledge integer not null default 0,
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
	Email text not null,	
        FOREIGN KEY (CompanyID) REFERENCES Company(ID)
	);`,

		`CREATE TABLE Settings (
        ID integer primary key autoincrement,
        UserID integer not null,
        MapAPI text not null default 'GoogleMaps',
	Interpolate integer not null default 1,
	SnaptoRoad integer not null default 1,
	CameraPanTrigger integer not null default 10,
        RadioCommunication integer not null default 1,
        DataCommunication integer not null default 1,
        SecurityRemoteAdmin integer not null default 0,
        SecurityConsoleAccess integer not null default 0,
        SecurityAdminPasswordReset integer not null default 0,
        MobileSmartPhoneAccess integer not null default 0,
        MobileShowBusLocation integer not null default 0,
        FOREIGN KEY (UserID) REFERENCES User(ID)
	);`,

		"INSERT INTO Company (Name, MaxUsers) VALUES ('myClubLink' , 1);",
		"INSERT INTO User (FirstName, LastName, CompanyID, Password, AccessLevel, Email) VALUES ('guest','user', 1, 'guest', 0, 'guest@myclublink.com.au');",
		"INSERT INTO Settings (UserID, MapAPI) VALUES (1, 'GoogleMaps');",

		"INSERT INTO Company (Name, MaxUsers, Expiry) VALUES ('Sussex Inlet RSL Group', 5, '2013-07-20 12:00:00');",
		"INSERT INTO User (FirstName, LastNAme, CompanyID, Password, AccessLevel, Email) VALUES ('Craig', 'Smith', 2, 'craig', 10, 'craig@sussexinlet.com.au');",
		"INSERT INTO Settings (UserID, MapAPI) VALUES (2, 'GoogleMaps');",
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
	fmt.Printf("Web socket requested from %s", r.RemoteAddr)
	if r.Method != "GET" {
		fmt.Printf("GET method request for socket. Not allowed\n")
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
		fmt.Printf("Not a websocket handshake \n")
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	connections = append(connections, connection)
	fmt.Printf("Amount of clients listening is %d\n", len(connections))
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
	viewRouter.HandleFunc("/system/support", views["ViewSupport"].(func(http.ResponseWriter, *http.Request)))
	viewRouter.HandleFunc("/system/map", views["ViewMap"].(func(http.ResponseWriter, *http.Request)))
	viewRouter.HandleFunc("/system/report", views["ViewReport"].(func(http.ResponseWriter, *http.Request)))
	viewRouter.HandleFunc("/", views["ViewInvalid"].(func(http.ResponseWriter, *http.Request)))

	//Action Routes
	actionRouter.HandleFunc("/system/login", actions["ActionLogin"].(func(http.ResponseWriter, *http.Request)))
	actionRouter.HandleFunc("/system/settings", actions["ActionSettings"].(func(http.ResponseWriter, *http.Request)))
	actionRouter.HandleFunc("/system/historicalroute", actions["ActionHistorialRoute"].(func(http.ResponseWriter, *http.Request)))
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

func init() {
	gob.Register(User{})
	gob.Register(Company{})
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

	//handle web requests in a seperate go-routine
	go handleHTTP()

	var recreateConnection bool = true
	var tcpcon net.Conn

	//wait around for tcp requests and handle them when they come in
	for {
		if recreateConnection {
			lnk, err := net.Listen("tcp", *service)
			if err != nil {
				fmt.Printf("Failed to get tcp listener - %s", err.Error())
				os.Exit(1)
			}
			fmt.Printf("Listening on TCP Port %s\n", *service)

			tcpcon, err = lnk.Accept()
				
			fmt.Printf("Link Accepted")
			if err != nil {
				fmt.Printf("Failed to create tcp connection - %s", err)
				os.Exit(1)
			}
			recreateConnection = handleClient(Db, tcpcon.(*net.TCPConn))
			if recreateConnection {
				lnk.Close()
			}

		}
	}
}

func updateClient(entry *GPSRecord) {

	if connections == nil {
		//fmt.Printf("No clients listening.. not reporting")
		return
	}

	//fmt.Printf("Responding to %d listening clients\n", len(connections))
	for index, client := range connections {
		//get a websocket writer
		wswriter, _ := client.NextWriter(websocket.OpText)

		if wswriter != nil {
			io.WriteString(wswriter, Response{"Entry": entry}.String())
		} else {
			//fmt.Printf("No ws writer available\n") //this web socket was abruptly closed so we need to close that client and remove it from the connections slice
			client.Close()
			//remove from slice
			connections = append(connections[:index], connections[index+1:]...)
		}

	}
}

func logEntry(entry *GPSRecord) {

	_, err := Db.Exec("INSERT INTO GPSRecords (Message, Latitude, Longitude, Speed, Heading, Fix, DateTime, BusID) VALUES ( ? , ?, ? , ? , ? ,? ,? , ?)",
		entry.Message,
		entry.Latitude,
		entry.Longitude,
		entry.Speed,
		entry.Heading,
		entry.Fix,
		entry.Date,
		entry.ID)

	if err != nil {
		fmt.Printf("Failed to insert row %s", err)
	}

	//daytime := time.Now().String()
}

//palm off reading and writing to a go routine
func handleClient(Db *sql.DB, conn *net.TCPConn) bool {
	var buff = make([]byte, 512)
	var entry GPSRecord
	
	conn.SetDeadline(time.Now().Add(time.Second + time.Second + time.Second + time.Second))
	conn.SetReadBuffer(512)
	var n int
	var err error
	var data bool = true
	for data {
		n, err = conn.Read(buff)
		if err != nil {
			fmt.Printf("Error occured - %s", err.Error())
			fmt.Printf("Error reading from TCP - Will recreate the connection \n")
			return true
		}
		conn.SetDeadline(time.Now().Add(time.Second + time.Second + time.Second + time.Second))
		fmt.Printf("Sentence was %s", string(buff))
		gpsfields := strings.Split(string(buff[:n]), ",")
		if len(gpsfields) != 8 {
			fmt.Printf("Error. GPS fields length is incorrect. Is %d should be %d", len(gpsfields), 8)
			fmt.Printf("The source string was %s\n", string(buff[:n]))
			continue
		}
		//All data is validated on the logger end so I'm going to assume for now that Parsing will be fine. Perhaps a network error could occur and I'll fix that up later

		entry.Message = gpsfields[0][1:]
		entry.Latitude = gpsfields[1][1:]
		entry.Longitude = gpsfields[2]
		entry.Speed, _ = strconv.ParseFloat(gpsfields[3][1:], 32)
		entry.Heading, _ = strconv.ParseFloat(gpsfields[4][1:], 32)
		entry.Date, _ = time.Parse(time.RFC3339, gpsfields[5][1:]) //todo pull out just the date component and format
		entry.Fix = gpsfields[6][1:] == "true"
		entry.ID = gpsfields[7][1:]

		fmt.Printf("Message %s Lat %s Long %s speed %f heading %f fix %t date %s id %s\n",
			entry.Message,
			entry.Latitude,
			entry.Longitude,
			entry.Speed,
			entry.Heading,
			entry.Fix,
			entry.Date,
			entry.ID)

		if string(buff[0:1]) != "T" {
			go logEntry(&entry) //save to database
		} else {
			fmt.Printf("Replayed packets. Not saving to DB\n")
		}

		updateClient(&entry) //notify any HTTP observers //make this a goroutine later
		conn.Write([]byte("OK\n"))
	}
	return false
}

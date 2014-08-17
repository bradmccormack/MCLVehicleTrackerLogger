package main

import (
	"./socket"
	"./types"
	"./utility"
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//set the domain based upon the path the executable was run from
var domain string = "dev.myclublink.com.au"
var service = flag.String("service", ":6969", "tcp port to bind to")
var addr = flag.String("addr", ":8080", "http(s) service address")

var random *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano())) //new random with unix time nano seconds as seed
//Session information
var store = sessions.NewCookieStore([]byte("emtec789")) //this needs to be randomized something from /dev/random

var Db *sql.DB

var actions = map[string]interface{}{
	"ActionInvalid": func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid Action", 403)
	},
	"ActionLogout": func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "application/json")

		if Db == nil {
			log.Fatal(Db)
		}

		session, _ := store.Get(r, "data")

		var user types.User = session.Values["User"].(types.User)

		//Update DB
		Db.Exec("UPDATE ApplicationLogin SET LoggedOut = CURRENT_TIMESTAMP WHERE UserID = ? AND LoggedOut IS NULL", user.ID)

		var hash = utility.GetSocketHash(r, user.FirstName, user.LastName)
		socket.WebSocketClose(hash)

		//clear cookie
		session.Values["User"] = ""
		session.Values["Company"] = ""
		session.Values["Settings"] = ""

		if err := session.Save(r, w); err != nil {
			fmt.Printf("Can't save session data (%s)\n", err.Error())
		}

		fmt.Fprint(w, types.JSONResponse{"success": true, "message": "Log out ok"})

	},

	"ActionLogin": func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		Db.Exec("ATTACH DATABASE 'license.key' AS L")

		var user types.User
		var company types.Company
		var settings types.Settings

		name := r.FormValue("name")
		password := r.FormValue("password")

		if Db == nil {
			log.Fatal(Db)
		}

		result := Db.QueryRow(`	
			SELECT U.ID, U.FirstName, U.LastName, U.Password, U.AccessLevel, U.Email, C.Name, C.MaxUsers, C.Expiry, C.LogoPath, S.MapAPI, S.Interpolate, S.SnaptoRoad, S.CameraPanTrigger,
			CS.RadioCommunication, CS.DataCommunication, CS.SecurityRemoteAdmin, CS.SecurityConsoleAccess, CS.SecurityAdminPasswordReset, CS.MobileSmartPhoneAccess, CS.MobileShowBusLocation,
			CS.MinZoom, CS.MaxZoom, CS.ClubBoundaryKM
			FROM User U
			LEFT JOIN COMPANY AS C on C.ID = U.CompanyID
			LEFT JOIN Settings AS S on S.UserID = U.ID
		    LEFT JOIN CompanySettings AS CS on CS.CompanyID = C.ID
			WHERE UPPER(U.FirstName) = ? AND U.Password = ?`,
			strings.ToUpper(name), password).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Password, &user.Accesslevel, &user.Email, &company.Name, &company.Maxusers, &company.Expiry,
			&company.LogoPath, &settings.MapAPI, &settings.Interpolate, &settings.SnaptoRoad, &settings.CameraPanTrigger, &settings.RadioCommunication, &settings.DataCommunication, &settings.SecurityRemoteAdmin,
			&settings.SecurityConsoleAccess, &settings.SecurityAdminPasswordReset, &settings.MobileSmartPhoneAccess, &settings.MobileShowBusLocation, &settings.MinZoom, &settings.MaxZoom, &settings.ClubBoundaryKM)

		switch {
		case result == sql.ErrNoRows:
			fmt.Fprint(w, types.JSONResponse{"success": false, "errors": []string{"Incorrect User/Password specified"}})

		case result != nil:
			log.Fatal(result)
		default:
			var Errors []string

			var LoggedInCount int

			var result error
			result = Db.QueryRow("SELECT COUNT(1) FROM L.ApplicationLogin WHERE LoggedOut IS NULL AND UserID = ?", user.ID).Scan(&LoggedInCount)
			if result != nil {
				log.Fatal(result)
			}

			if LoggedInCount == company.Maxusers {
				Errors = append(Errors, "Amount of users logged in ("+strconv.Itoa(LoggedInCount)+") matches your license limit ("+strconv.Itoa(company.Maxusers)+")")
			}

			var ExpiryDate time.Time

			const layout = "2006-01-2 15:4:5" //http://golang.org/src/pkg/time/format.go
			ExpiryDate, _ = time.Parse(layout, company.Expiry)

			if ExpiryDate.Unix() < time.Now().Unix() {
				Errors = append(Errors, "Your license has expired. Please contact myClublink support to renew your License")
			}

			if len(Errors) == 0 {
				Db.Exec("INSERT INTO L.ApplicationLogin (UserID) VALUES ( ?)", user.ID)
				session, _ := store.Get(r, "data")
				session.Values["User"] = user
				session.Values["Company"] = company
				session.Values["Settings"] = settings
				session.Options = &sessions.Options{
					Path:   "/",
					MaxAge: 86400, //1 day
				}

				if err := session.Save(r, w); err != nil {
					fmt.Printf("Can't save session data (%s)\n", err.Error())
				}

				fmt.Fprint(w, types.JSONResponse{"success": true, "message": "Login ok", "user": user, "company": company, "settings": settings})
			} else {
				fmt.Fprint(w, types.JSONResponse{"success": false, "message": "Login failed", "errors": Errors})
			}

		}
	},

	"ActionSettingsPassword": func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		session, _ := store.Get(r, "data")
		var user types.User = session.Values["User"].(types.User)
		var settings types.Settings = session.Values["Settings"].(types.Settings)

		var f map[string]interface{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&f)
		if err != nil {
			log.Fatal(err)
		}
		var password string

		_ = Db.QueryRow("SELECT Password FROM License.User WHERE ID = ?", user.ID).Scan(&password)
		if password == f["passwordold"] {
			//If only Allow admins to reset password is NOT set then update the users password
			if settings.SecurityAdminPasswordReset == 0 {
				Db.Exec("UPDATE License.User SET Password = ? WHERE ID = ?", f["password"], user.ID)
			} else {
				if user.Accesslevel == 10 {
					Db.Exec("UPDATE License.User SET Password = ? WHERE ID = ?", user.ID)
				}
			}

			user.Password = f["password"].(string)
			session.Values["User"] = user
			if err := session.Save(r, w); err != nil {
				fmt.Printf("Can't save session data (%s)\n", err.Error())
			}
			fmt.Fprint(w, types.JSONResponse{"success": "Password Updated"})

		} else {
			fmt.Fprint(w, types.JSONResponse{"error": "Old Password incorrect"})
		}
	},
	"ActionSettings": func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		session, _ := store.Get(r, "data")
		var user types.User = session.Values["User"].(types.User)
		var settings types.Settings = session.Values["Settings"].(types.Settings)

		var f map[string]interface{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&f)
		if err != nil {
			log.Fatal(err)
		}

		Db.Exec("BEGIN EXCLUSIVE TRANSACTION")
		Db.Exec(`UPDATE License.Settings SET MapAPI = ?, Interpolate = ?, SnaptoRoad = ?, CameraPanTrigger = ? WHERE UserID = ? `,
			f["MapAPI"], f["Interpolate"], f["SnaptoRoad"], f["CameraPanTrigger"], user.ID)

		//If the user is an admin then allow update of admin level fields

		if user.Accesslevel == 10 {

			Db.Exec(`UPDATE License.CompanySettings SET RadioCommunication = ?, DataCommunication = ?,
                             SecurityRemoteAdmin = ?, SecurityConsoleAccess = ?, SecurityAdminPasswordReset = ?,
                             MobileSmartPhoneAccess = ?, MobileShowBusLocation = ?, MinZoom = ?, MaxZoom = ?, ClubBoundaryKM = ? WHERE CompanyID = (SELECT CompanyID FROM License.User WHERE ID = ?)`,
				f["RadioCommunication"], f["DataCommunication"],
				f["SecurityRemoteAdmin"], f["SecurityConsoleAccess"], f["SecurityAdminPasswordReset"],
				f["MobileSmartPhoneAccess"], f["MobileShowBusLocation"], f["MinZoom"], f["MaxZoom"], f["ClubBoundaryKM"], user.ID)

		}

		Db.Exec("COMMIT TRANSACTION")

		//Update the cookie too
		settings.MapAPI = f["MapAPI"].(string)

		//TODO see if I can improve this verbose crappy code
		if f["Interpolate"].(bool) {
			settings.Interpolate = 1
		} else {
			settings.Interpolate = 0
		}

		if f["SnaptoRoad"].(bool) {
			settings.SnaptoRoad = 1
		} else {
			settings.SnaptoRoad = 0
		}

		settings.CameraPanTrigger = int(f["CameraPanTrigger"].(float64))
		settings.MinZoom = int(f["MinZoom"].(float64))
		settings.MaxZoom = int(f["MaxZoom"].(float64))
		settings.ClubBoundaryKM = int(f["ClubBoundaryKM"].(float64))

		if f["RadioCommunication"].(bool) {
			settings.RadioCommunication = 1
		} else {
			settings.RadioCommunication = 0
		}

		if f["DataCommunication"].(bool) {
			settings.DataCommunication = 1
		} else {
			settings.DataCommunication = 0
		}

		if f["SecurityRemoteAdmin"].(bool) {
			settings.SecurityRemoteAdmin = 1
		} else {
			settings.SecurityRemoteAdmin = 0
		}

		if f["SecurityConsoleAccess"].(bool) {
			settings.SecurityConsoleAccess = 1
		} else {
			settings.SecurityConsoleAccess = 0
		}

		if f["SecurityAdminPasswordReset"].(bool) {
			settings.SecurityAdminPasswordReset = 1
		} else {
			settings.SecurityAdminPasswordReset = 0
		}

		if f["MobileSmartPhoneAccess"].(bool) {
			settings.MobileSmartPhoneAccess = 1
		} else {
			settings.MobileSmartPhoneAccess = 0
		}

		if f["MobileShowBusLocation"].(bool) {
			settings.MobileShowBusLocation = 1
		} else {
			settings.MobileShowBusLocation = 0
		}

		session.Values["Settings"] = settings
		if err := session.Save(r, w); err != nil {
			fmt.Printf("Can't save session data (%s)\n", err.Error())
		}
		fmt.Fprint(w, types.JSONResponse{"success": true})

	},
	"ActionHistorialRoute": func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "application/json")
		//map with string key and slice of string slices
		var Route = make(map[string][][]string)

		dateFrom := r.FormValue("dateFrom")
		dateTo := r.FormValue("dateTo")

		fmt.Printf("DateFrom is %s, DateTo is %s", dateFrom, dateTo)

		rows, err := Db.Query("SELECT BusID, Latitude, Longitude, Speed, Heading, Fix, DateTime FROM GPSRecords WHERE datetime >=? AND datetime <=? AND Fix AND SPEED > 10 GROUP BY id ORDER BY datetime asc", dateFrom, dateTo)
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

		fmt.Fprint(w, types.JSONResponse{"success": true, "data": Route})

	},
}

//Note - Template caching needs to be implemented http://golang.org/doc/articles/wiki/ There is an inefficiency in this code: renderTemplate calls ParseFiles every time a page is rendered.
var views = map[string]interface{}{

	"ViewInvalid": func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid view", 403)
	},

	"ViewLogin": func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		session, _ := store.Get(r, "data")
		if session == nil {
			http.Error(w, "Unauthorized", 401)
		} else {
			var user types.User
			var company types.Company
			var settings types.Settings
			user = session.Values["User"].(types.User)

			company = session.Values["Company"].(types.Company)
			settings = session.Values["Settings"].(types.Settings)
			fmt.Fprint(w, types.JSONResponse{"success": true, "message": "Login OK", "user": user, "company": company, "settings": settings})
		}

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
		//TODO redirect to root if not logged in - USE middleware for this later
		w.Header().Add("Content-Type", "application/json")

		//session, _ := store.Get(r, "session")

		var percentAvailable int = random.Intn(75)
		availability := [...]int{percentAvailable, 100 - percentAvailable}

		//TODO restrict these reports to a range of dates
		//dateFrom := r.FormValue("dateFrom")
		//dateTo := r.FormValue("dateTo")

		//23 3e

		var distance float64
		var weekday int

		//init all days to 0
		var kmPerDay [7]float64
		for i := 0; i < 7; i++ {
			kmPerDay[i] = 0
		}

		rows, err := Db.Query(`
                        SELECT strftime('%w', datetime(GPSR1.DateTime, 'localtime')) AS Weekday,
			SUM((strftime('%s',datetime(GPSR2.DateTime, "localtime")) - strftime('%s',datetime(GPSR1.DateTime, "localtime"))) *
			( (GPSR1.Speed + GPSR2.Speed) /2 )  / 3600) as Distance
			FROM GPSRecords GPSR1, GPSRecords GPSR2
			WHERE GPSR1.ID = GPSR2.ID -1
			AND GPSR1.Fix = 1
			GROUP BY Weekday`)

		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			if err := rows.Scan(&weekday, &distance); err != nil {
				log.Fatal(err)
			}
			kmPerDay[weekday] = distance

		}

		fmt.Fprint(w, types.JSONResponse{"Availability": availability, "KMPerDay": kmPerDay})

	},

	"ViewMap": func(w http.ResponseWriter, r *http.Request) {
		//TODO redirect to root if not logged in - USE middleware for this later
		w.Header().Add("Content-Type", "text/html")

		session, _ := store.Get(r, "session")
		fmt.Printf("Session is %s", types.JSONResponse{"session": session})

		var err error
		t := template.New("Map")
		t, err = template.ParseFiles("templates/map.html")
		if err != nil {
			log.Fatal("Failed to read the template file for map. Fix it")
		}
		t.Execute(w, session.Values)
	},

	"ViewLicense": func(w http.ResponseWriter, r *http.Request) {
		//TODO redirect to root if not logged in - USE middleware for this later
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
		//TODO redirect to root if not logged in - USE middleware for this later
		w.Header().Add("Content-Type", "text/html")

		session, _ := store.Get(r, "session")

		var mapAPI string
		var interpolate, snaptoroad bool
		var user types.User = session.Values["User"].(types.User)

		result := Db.QueryRow(`
                        SELECT S.MapAPI, S.Interpolate, S.SnaptoRoad
                        FROM License.Settings S, License.User U
			WHERE S.UserID = U.ID 
			AND U.ID = ?`, user.ID).Scan(&mapAPI, &interpolate, &snaptoroad)

		switch {
		case result != nil:
			log.Fatal(result)
		default:
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

		t.Execute(w, session.Values)
	},
}

func socketInit(w http.ResponseWriter, r *http.Request) {
	socket.WebSocketInit(w, r, store)
}

func handleHTTP() {
	Router := mux.NewRouter()

	viewRouter := Router.Methods("GET").Subrouter()
	actionRouter := Router.Methods("POST").Subrouter()

	//Handle web socket traffic specially
	Router.HandleFunc("/ws", socketInit)

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
	actionRouter.HandleFunc("/system/logout", actions["ActionLogout"].(func(http.ResponseWriter, *http.Request)))
	actionRouter.HandleFunc("/system/settings", actions["ActionSettings"].(func(http.ResponseWriter, *http.Request)))
	actionRouter.HandleFunc("/system/settings/password", actions["ActionSettingsPassword"].(func(http.ResponseWriter, *http.Request)))
	actionRouter.HandleFunc("/system/historicalroute", actions["ActionHistorialRoute"].(func(http.ResponseWriter, *http.Request)))
	actionRouter.HandleFunc("/", actions["ActionInvalid"].(func(http.ResponseWriter, *http.Request)))

	//Use the router
	http.Handle("/", Router)

	fmt.Printf("\nListening for HTTP on %s\n", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		fmt.Printf("\nFailed to listen for http on %s", *addr)
		log.Fatal("\nError: ", err)
	}

}

func init() {
	gob.Register(types.User{})
	gob.Register(types.Company{})
	gob.Register(types.Settings{})
}

func main() {

	flag.Parse()

	var err error

	if _, err := os.Stat("backend.db"); err != nil {
		log.Fatal("\nError: ", err)
	}

	if _, err := os.Stat("license.key"); err != nil {
		log.Fatal("\nError: ", err)
	}

	Db, err = sql.Open("sqlite3", "backend.db")
	if err != nil {
		fmt.Printf("Cannot open database backend.db . Exiting\n")
		os.Exit(1)
	}

	//try to attach the license.key if it opens then close it and attach
	LDb, err := sql.Open("sqlite3", "license.key")
	if err != nil {
		fmt.Printf("\nCannot open database license.key . Exiting\n")
		os.Exit(1)
	} else {
		fmt.Printf("\nLicense.key opened correctly")
	}

	LDb.Close()

	defer Db.Close()

	//handle web requests in a seperate go-routine
	go handleHTTP()

	//Socket related channels
	WSDataChannel := make(chan types.Record, 100)
	WSCommandChannel := make(chan int32)

	//
	NetworkChannel := make(chan int32, 1)

	go socket.Monitor(WSDataChannel, WSCommandChannel)
	go connectionManager(NetworkChannel, WSCommandChannel, WSDataChannel)

	<-NetworkChannel

}

func connectionManager(NetworkChannel chan int32, WSCommandChannel chan<- int32, WSDataChannel chan<- types.Record) {

	//select from first available channel ipc - note this blocks until there is data in one of the channels
	select {
	//keep slurping records from the bufered channel and farm them out to UpdateClient as a goroutine
	case msg := <-NetworkChannel:
		switch msg {
		case types.COMMAND_RECONNECT:
			WSCommandChannel <- types.COMMAND_QUIT

			//lnk.Close()

			lnk, err := net.Listen("tcp", *service)
			if err != nil {
				fmt.Printf("\nFailed to get tcp listener - %s", err.Error())
				os.Exit(1)
			}
			for {
				//this blocks until there is a connection
				tcpcon, err := lnk.Accept()

				fmt.Printf("\nLink Accepted - Receiving packets from Vehicle\n")
				if err != nil {
					fmt.Printf("\nFailed to create tcp connection - %s", err)
					os.Exit(1)
				}
				fmt.Printf("\nListening on TCP Port %s", *service)

				go handleClient(tcpcon.(*net.TCPConn), WSDataChannel, NetworkChannel)

			}

		}

	}
}

func logEntry(entry *types.GPSRecord, diagnostic *types.DiagnosticRecord) {

	_, err := Db.Exec("BEGIN TRANSACTION")
	_, err = Db.Exec("INSERT INTO GPSRecords (Message, Latitude, Longitude, Speed, Heading, Fix, DateTime, BusID) VALUES ( ? , ?, ? , ? , ? ,? ,? , ?)",
		entry.Message,
		entry.Latitude,
		entry.Longitude,
		entry.Speed,
		entry.Heading,
		entry.Fix,
		entry.Date,
		entry.ID)

	_, err = Db.Exec("INSERT INTO DiagnosticRecords (CPUTemperature, CPUVoltage, CPUFrequency, MemoryFree) VALUES (?, ?, ?, ?)",
		diagnostic.CPUTemp,
		diagnostic.CPUVolt,
		diagnostic.CPUFreq,
		diagnostic.MemFree)

	Db.Exec("COMMIT TRANSACTION")
	if err != nil {
		fmt.Printf("Failed to insert row %s", err)
	}

}

func handleClient(conn *net.TCPConn, WSDataChannel chan<- types.Record, NetworkChannel chan<- int32) {

	//defer anonymous func to handle panics - most likely panicking from garbage that was to be parsed.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from a panic \n", r)
		}
	}()

	var buff = make([]byte, 512)
	var incomingpacket types.Packet

	var R types.Record
	R.GPS = new(types.GPSRecord)
	R.Diagnostic = new(types.DiagnosticRecord)

	var n int
	var err error
	for {
		n, err = conn.Read(buff)

		if err != nil {
			fmt.Printf("\nError occured - %s", err.Error())
			fmt.Printf("\nError reading from TCP - Will recreate the connection \n")
			fmt.Printf("About to send COMMAND_RECONNECT")
			NetworkChannel <- types.COMMAND_RECONNECT
			fmt.Printf("Sent COMMAND_RECONNECT")
			return
		}

		//lets unmarshal those JSON bytes into the map https://groups.google.com/forum/#!topic/golang-nuts/77HJlZhWXpk  note to slice properly otherwise it chockes on trying to decode the full buffer
		err := json.Unmarshal(buff[:n], &incomingpacket)
		if err != nil {
			fmt.Printf("Failed to decode the JSON bytes -%s\n", err.Error())
		}

		//fmt.Printf("\nSentence was %s", string(incomingpacket["sentence"]))
		//fmt.Printf("\nDiagnostic data was %s", string(incomingpacket["diagnostics"]))

		diagnosticfields := strings.Split(string(incomingpacket["diagnostics"]), ",")
		if len(diagnosticfields) != 4 {
			fmt.Printf("Error. Diagnostic fields length is incorrect. Is %d should be %d", len(diagnosticfields), 4)
			fmt.Printf("The source string was %s\n", string(incomingpacket["diagnostics"]))
		}

		gpsfields := strings.Split(string(incomingpacket["sentence"]), ",")

		if len(gpsfields) != 7 {
			fmt.Printf("Error. GPS fields length is incorrect. Is %d should be %d\n", len(gpsfields), 7)
			fmt.Printf("The source string was %s\n", string(incomingpacket["sentence"]))
			continue
		}

		R.Diagnostic.CPUTemp, _ = strconv.ParseFloat(diagnosticfields[0][2:], 32)
		R.Diagnostic.CPUVolt, _ = strconv.ParseFloat(diagnosticfields[1][2:], 32)
		R.Diagnostic.CPUFreq, _ = strconv.ParseFloat(diagnosticfields[2][2:], 32)
		R.Diagnostic.MemFree, _ = strconv.ParseUint(diagnosticfields[3][2:], 10, 64)

		R.GPS.Message = gpsfields[0][1:]
		R.GPS.Latitude = gpsfields[0][2:]
		R.GPS.Longitude = gpsfields[1]
		R.GPS.Speed, _ = strconv.ParseFloat(gpsfields[2][1:], 32)
		R.GPS.Heading, _ = strconv.ParseFloat(gpsfields[3][1:], 32)
		R.GPS.Date, _ = time.Parse(time.RFC3339, gpsfields[4][1:])
		R.GPS.Fix = gpsfields[5][1:] == "true"
		R.GPS.ID = gpsfields[6][1:]

		/*
			fmt.Printf("Temp %d, Voltage %d, Frequency %d, MemoryFree %d",
				R.Diagnostic.CPUTemp,
				R.Diagnostic.CPUVolt,
				R.Diagnostic.CPUFreq,
				R.Diagnostic.MemFree)

			fmt.Printf("Message %s Lat %s Long %s speed %f heading %f fix %t date %s id %s\n",
				R.GPS.Message,
				R.GPS.Latitude,
				R.GPS.Longitude,
				R.GPS.Speed,
				R.GPS.Heading,
				R.GPS.Fix,
				R.GPS.Date,
				R.GPS.ID)
		*/

		if string(incomingpacket["sentence"][0:1]) != "T" {
			go logEntry(R.GPS, R.Diagnostic)
		}

		WSDataChannel <- R
		conn.Write([]byte("OK\n"))
	}
	return
}

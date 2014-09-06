package http

import (
	"../dao"
	"../socket"
	"../types"
	"../utility"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func HttpRouter(BindIP *string) {

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

	fmt.Printf("\nListening for HTTP on %s", *BindIP)
	err := http.ListenAndServe(*BindIP, nil)
	if err != nil {
		fmt.Printf("\nFailed to listen for http on %s", *BindIP)
		log.Fatal("\nError: ", err)
	}

}

func socketInit(w http.ResponseWriter, r *http.Request) {
	socket.WebSocketInit(w, r, store)
}

//Session information
//get random bytes
var store = sessions.NewCookieStore([]byte("emtec789")) //this needs to be randomized something from /dev/random

var actions = map[string]interface{}{
	"ActionInvalid": func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid Action", 403)
	},
	"ActionLogout": func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "application/json")

		if db == nil {
			log.Fatal(db)
		}

		session, _ := store.Get(r, "data")

		var user types.User = session.Values["User"].(types.User)

		//Update DB
		dao.LogOutUser(user.ID)

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

		name := r.FormValue("name")
		password := r.FormValue("password")

		user, company, settings, errors := dao.Login(name, password)
		
		if(len(errors == 0) {
			fmt.Fprint(w, types.JSONResponse{"success": true, "message": "Login ok", "user": user, "company": company, "settings": settings})
			session, _ := store.Get(r, "data")
			
			//TODO if this user is currently logged in then log them out
			//TODO log out old users who have been logged in more than 24 hours 
			
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

		password := dao.GetPassword(user.ID)
	
		if password == f["passwordold"] {
			//If only Allow admins to reset password is NOT set then update the users password
			if settings.SecurityAdminPasswordReset == 0 {
				dao.SetPassword(user.ID, f["password"])
			} else {
				if user.Accesslevel == 10 {
					dao.SetPassword(user.ID, f["password"])
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

		var settings map[string]interface{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&settings)
		if err != nil {
			log.Fatal(err)
		}

		dao.SetSettings(&settings)
		
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

		rows, err := db.Query("SELECT BusID, Latitude, Longitude, Speed, Heading, Fix, DateTime FROM GPSRecords WHERE datetime >=? AND datetime <=? AND Fix AND SPEED > 10 GROUP BY id ORDER BY datetime asc", dateFrom, dateTo)
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

		var random *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano())) //new random with unix time nano seconds as seed
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

		rows, err := db.Query(`
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

		result := db.QueryRow(`
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

package http

import (
	"../dao"
	"../socket"
	"../types"
	"../utility"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

/*
	There are a lot of work remaining on this

	There is next to no security - This is due to just a testing/prototyping phase. Security would need to be implemented before production.

	Previously view endpoints were hit and Go templates were executed. Now it is all done client side with AngularJS and partials.
	These endpoints need to be audited and cleaned up.

	Currently all the client side JS resides in separate controllers etc. The Go binary needs to concat them all and minify.
	Also the particular map API vendor needs to be added conditionally not all of them. Nginx config needs to make sure JS requests go through backend.

	CSS needs to be minimized too and audited to see what can be removed.
*/

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

		user, company, settings, errors := dao.LoginUser(name, password)
		//TODO IV + salt hash the password and compare against exiting hashed password - This is just testing code so plaintext is OK for now.


		if len(errors) == 0 {

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
				fmt.Fprint(w, types.JSONResponse{"success": true, "message": "Login OK", "user": user, "company": company, "settings": settings})

		} else {
				fmt.Fprint(w, types.JSONResponse{"success": false, "message": "Login Failed", "errors": Errors})
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

		password := dao.GetPassword(user.ID)

		//TODO hash and check against existing hashed password
		if password == f["passwordold"] {
			//If only Allow admins to reset password is NOT set then update the users password
			if settings.SecurityAdminPasswordReset == 0 {
				dao.SetPassword(user.ID, f["password"].(string))
			} else {
				if user.Accesslevel == 10 {
					dao.SetPassword(user.ID, f["password"].(string))
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
		err := decoder.Decode(&settings)
		if err != nil {
			log.Fatal(err)
		}

		dao.SetSettings(&user, &f)

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

		dateFrom := r.FormValue("dateFrom")
		dateTo := r.FormValue("dateTo")
		Route := dao.GetHistoricalRoute(dateFrom, dateTo)

		//This was very problematic. Where the packets were being recorded every second this would repond back with way too much data in the body and client side rendering whould choke.
		//There needs to be better heuristics to produce an optimized dataset that is will render quickly.
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

		kmPerDay := dao.GetKMReport()

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

		var user types.User = session.Values["User"].(types.User)
		mapAPI, interpolate, snaptoroad := dao.GetSettings(&user)

		session.Values["Settings"] = map[string]interface{}{
			"MapAPI":      mapAPI,
			"Interpolate": interpolate,
			"SnaptoRoad":  snaptoroad,
		}

		session.Save(r, w)

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

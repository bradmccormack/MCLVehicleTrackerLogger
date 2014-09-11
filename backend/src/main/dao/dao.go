package dao

import (
	"../types"
	"database/sql"
	"fmt"
	//_ "github.com/mattn/go-sqlite3" underscore must be scope/visibility of the package
	"github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var db *sql.DB

func init() {
	sql.Register("sqlite3_with_extensions",
		&sqlite3.SQLiteDriver{
			Extensions: []string{
				"./sqlite3_mod_distance.so",
			},
		})

}

func Open() {
	var err error

	db, err = sql.Open("sqlite3_with_extensions", "backend.db")
	if err != nil {
		fmt.Printf("Cannot open database backend.db . Exiting\n")
		os.Exit(1)
	}

	//try to attach the license.key if it opens then close it and attach
	LDb, err := sql.Open("sqlite3", "license.key")
	if err != nil {
		fmt.Printf("\nCannot open database license.key . Exiting\n")
		os.Exit(1)
	}
	LDb.Close()

}

func SavePacket(entry *types.GPSRecord, diagnostic *types.DiagnosticRecord) {

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

	_, err := db.Exec("BEGIN TRANSACTION")
	_, err = db.Exec("INSERT INTO GPSRecords (Message, Latitude, Longitude, Speed, Heading, Fix, DateTime, BusID) VALUES ( ? , ?, ? , ? , ? ,? ,? , ?)",
		entry.Message,
		entry.Latitude,
		entry.Longitude,
		entry.Speed,
		entry.Heading,
		entry.Fix,
		entry.Date,
		entry.ID)

	_, err = db.Exec("INSERT INTO DiagnosticRecords (CPUTemperature, CPUVoltage, CPUFrequency, MemoryFree) VALUES (?, ?, ?, ?)",
		diagnostic.CPUTemp,
		diagnostic.CPUVolt,
		diagnostic.CPUFreq,
		diagnostic.MemFree)

	db.Exec("COMMIT TRANSACTION")
	if err != nil {
		fmt.Printf("Failed to insert row %s", err)
	}

}

func LogOutUser(UserId int) {
	db.Exec("UPDATE ApplicationLogin SET LoggedOut = CURRENT_TIMESTAMP WHERE UserID = ? AND LoggedOut IS NULL", UserId)
}

func LoginUser(name, password string) (User types.User, Company types.Company, Settings types.Settings, errors []string) {

	db.Exec("ATTACH DATABASE 'license.key' AS L")

	var user types.User
	var company types.Company
	var settings types.Settings
	var Errors []string

	resulterr := db.QueryRow(`	
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
	case resulterr == sql.ErrNoRows:
		Errors = append(Errors, "Incorrect User/Password specified")

	case resulterr != nil:
		log.Fatal(resulterr)

	default:
		var Errors []string

		var LoggedInCount int

		var result error

		resulterr = db.QueryRow("SELECT COUNT(1) FROM L.ApplicationLogin WHERE LoggedOut IS NULL AND UserID = ?", user.ID).Scan(&LoggedInCount)
		if resulterr != nil {
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
			db.Exec("INSERT INTO L.ApplicationLogin (UserID) VALUES ( ?)", user.ID)
		}
	}

	fmt.Printf("Logging in baby")
	return user, company, settings, Errors

}

func GetPassword(UserId int) string {
	var password string
	_ = db.QueryRow("SELECT Password FROM License.User WHERE ID = ?", UserId).Scan(&password)
	return password
}

func SetPassword(UserID int, Password string) {
	db.Exec("UPDATE License.User SET Password = ? WHERE ID = ?", Password, UserID)
}

func GetSettings(User *types.User) (mapAPI string, interpolate, snaptoroad bool) {
	//var mapAPI string
	//var interpolate, snaptoroad bool

	//get settings
	result := db.QueryRow(`
                        SELECT S.MapAPI, S.Interpolate, S.SnaptoRoad
                        FROM License.Settings S, License.User U
			WHERE S.UserID = U.ID 
			AND U.ID = ?`, User.ID).Scan(&mapAPI, &interpolate, &snaptoroad)

	if result != nil {
		log.Fatal(result)
	}
	return mapAPI, interpolate, snaptoroad
}

func SetSettings(User *types.User, settings *map[string]interface{}) {
	db.Exec("BEGIN EXCLUSIVE TRANSACTION")

	//dao/dao.go:165: invalid operation: settings["MapAPI"] (index of type *map[string]interface {})

	db.Exec(`UPDATE License.Settings SET MapAPI = ?, Interpolate = ?, SnaptoRoad = ?, CameraPanTrigger = ? WHERE UserID = ? `,
		(*settings)["MapAPI"], (*settings)["Interpolate"], (*settings)["SnaptoRoad"], (*settings)["CameraPanTrigger"], User.ID)

	//If the user is an admin then allow update of admin level fields

	if User.Accesslevel == 10 {

		db.Exec(`UPDATE License.CompanySettings SET RadioCommunication = ?, DataCommunication = ?,
                             SecurityRemoteAdmin = ?, SecurityConsoleAccess = ?, SecurityAdminPasswordReset = ?,
                             MobileSmartPhoneAccess = ?, MobileShowBusLocation = ?, MinZoom = ?, MaxZoom = ?, ClubBoundaryKM = ? WHERE CompanyID = (SELECT CompanyID FROM License.User WHERE ID = ?)`,
			(*settings)["RadioCommunication"], (*settings)["DataCommunication"],
			(*settings)["SecurityRemoteAdmin"], (*settings)["SecurityConsoleAccess"], (*settings)["SecurityAdminPasswordReset"],
			(*settings)["MobileSmartPhoneAccess"], (*settings)["MobileShowBusLocation"], (*settings)["MinZoom"], (*settings)["MaxZoom"], (*settings)["ClubBoundaryKM"], User.ID)

	}

	db.Exec("COMMIT TRANSACTION")
}

func GetHistoricalRoute(DateFrom, DateTo string) map[string][][]string {
	fmt.Printf("DateFrom is %s, DateTo is %s", DateFrom, DateTo)

	var Route = make(map[string][][]string)
	rows, err := db.Query("SELECT BusID, Latitude, Longitude, Speed, Heading, Fix, DateTime FROM GPSRecords WHERE datetime >=? AND datetime <=? AND Fix AND SPEED > 10 GROUP BY id ORDER BY datetime asc", DateFrom, DateTo)
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
	return Route

}

func GetKMReport() [7]float64 {

	var kmPerDay [7]float64
	var distance float64
	var weekday int

	for i := 0; i < 7; i++ {
		kmPerDay[i] = 0
	}
	//TODO restrict these reports to a range of dates
	//dateFrom := r.FormValue("dateFrom")
	//dateTo := r.FormValue("dateTo")

	//23 3e

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
	return kmPerDay
}

func GetStreetName(Latitude, Longitude string) string {

	//TODO move this later
	db.Exec("ATTACH DATABASE 'geodata.db' AS Geo")

	//SELECT * FROM Locations ORDER BY distance(Latitude, Longitude, 51.503357, -0.1199)
	var Name, Lat, Long string
	var Distance string

	_ = db.QueryRow(`SELECT P.Name, L.Lat,L.Long, distance(L.Lat, L.Long, ?, ?) AS Distance
						 FROM Geo.LatLong AS L
						 JOIN Geo.POI AS P ON P.Id = L.POIID
						 WHERE Distance < 0.01
						 ORDER BY Distance
						 LIMIT 1`, Latitude, Longitude).Scan(&Name, &Lat, &Long, &Distance)

	/*
		if err != nil {
			fmt.Printf("Error happened %s\n", err)
		} else {
			fmt.Printf("\nName = %s, Lat = %s, Long = %s, dist = %s", Name, Lat, Long, Distance)
		}
	*/
	return Name
}

func Close() {
	db.Close()
}

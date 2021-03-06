package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strconv"
)


type DbFunc func(Db *sql.DB)

//TODO add indexes
func main() {

	BackendDbSchema := []DbFunc{
		func(Db *sql.DB) {
			Db.Exec(`CREATE TABLE GPSRecords (
					ID INTEGER PRIMARY KEY AUTOINCREMENT,
					Message TEXT,
					Latitude TEXT NOT NULL,
					Longitude TEXT NOT NULL,
					Speed INTEGER NOT NULL,
					Heading float NOT NULL,
					Fix INTEGER NOT NULL,
					DateTime date NOT NULL DEFAULT current_timestamp,
					BusID TEXT NOT NULL);`)

			Db.Exec(`CREATE TABLE DiagnosticRecords (
					ID INTEGER PRIMARY KEY AUTOINCREMENT,
					CPUTemperature REAL,
					CPUVoltage REAL,
					CPUFrequency REAL,
					MemoryFree INTEGER,
					Date DateTime DEFAULT CURRENT_TIMESTAMP);`)

			Db.Exec(`CREATE TABLE Support (
					SupportID INTEGER PRIMARY KEY AUTOINCREMENT,
					UserID INTEGER NOT NULL,
					Subject TEXT NOT NULL,
					Body TEXT NOT NULL,
					DateTime date NOT NULL DEFAULT current_timestamp,
					FOREIGN KEY (UserID) REFERENCES User(ID));`)

			Db.Exec(`CREATE TABLE Errors (
					ID INTEGER PRIMARY KEY AUTOINCREMENT,
					GPSRecordID INTEGER NOT NULL,
					Error TEXT,
					DateTime date NOT NULL DEFAULT current_timestamp,
					FOREIGN KEY (GPSRecordID) REFERENCES GPSrecords(id));`)

			Db.Exec(`CREATE TABLE Network (
					ID INTEGER PRIMARY KEY AUTOINCREMENT,
					GPSRecordID INTEGER NOT NULL,
					Acknowledge INTEGER NOT NULL DEFAULT 0,
					FOREIGN KEY (GPSRecordID) REFERENCES GPSRecords(id));`)

			Db.Exec("PRAGMA foreign_keys=ON;")
			Db.Exec("PRAGMA journal_mode=WAL;")

		},
	}

	LicenseDbSchema := []DbFunc{
		func(Db *sql.DB) {
			Db.Exec(`CREATE TABLE User (
					ID INTEGER PRIMARY KEY AUTOINCREMENT,
					FirstName TEXT NOT NULL,
					LastName TEXT NOT NULL,
					CompanyID INTEGER NOT NULL,
					Password TEXT NOT NULL,
					AccessLevel INTEGER NOT NULL DEFAULT 0,
					Email TEXT NOT NULL,
					FOREIGN KEY (CompanyID) REFERENCES Company(ID));`)

			Db.Exec(`CREATE TABLE Settings (
					ID INTEGER PRIMARY KEY AUTOINCREMENT,
					UserID INTEGER NOT NULL,
					MapAPI TEXT NOT NULL DEFAULT 'GoogleMaps',
					Interpolate INTEGER NOT NULL DEFAULT 1,
					SnaptoRoad INTEGER NOT NULL DEFAULT 1,
					CameraPanTrigger INTEGER NOT NULL DEFAULT 10,
					FOREIGN KEY (UserID) REFERENCES User(ID));`)

			Db.Exec(`CREATE TABLE ApplicationLogin (
					UserID INTEGER,
					LoggedIn date NOT NULL DEFAULT current_timestamp,
					LoggedOut date,
					PRIMARY KEY(UserID, LoggedIN));)`)

			Db.Exec(`CREATE TABLE Company (
					 ID INTEGER PRIMARY KEY AUTOINCREMENT,
					 Name TEXT NOT NULL,
					 Expiry date NOT NULL DEFAULT current_timestamp,
					 MaxUsers INTEGER NOT NULL DEFAULT 0,
					 LogoPath TEXT NOT NULL DEFAULT '');`)

			Db.Exec(`CREATE TABLE CompanySettings (
					ID integer primary key autoincrement,
					CompanyID integer not null,
					RadioCommunication integer not null default 1,
					DataCommunication integer not null default 1,
					SecurityRemoteAdmin integer not null default 0,
					SecurityConsoleAccess integer not null default 0,
					SecurityAdminPasswordReset integer not null default 0,
					MobileSmartPhoneAccess integer not null default 0,
					MobileShowBusLocation integer not null default 0,
					MinZoom integer not null default 10,
					Maxzoom integer not null default 2,
					HistoricalmapsKmMin integer not null default 10,
					ClubBoundaryKM integer not null default 100,
					FOREIGN KEY (CompanyID) REFERENCES Company(ID));`)

			//there needs to be a command that grabs current git master sha and updates this table with Version and VersionDate
			Db.Exec(`CREATE TABLE Version (
					ID TEXT PRIMARY KEY,
					SHA1 TEXT NOT NULL UNIQUE,
					ReleaseDate date NOT NULL DEFAULT current_timestamp);`)

			Db.Exec("PRAGMA foreign_keys=ON;")
		},
		func(Db *sql.DB) {

			//create a few users - note when I do salted hashing another function to update the passwords will be required LOL this is fucked
			Db.Exec(`INSERT INTO User (FirstName, LastName, CompanyID, Password, AccessLevel, Email)
					 VALUES ('guest','user', 1, 'guest', 0, 'guest@myclublink.com.au');`)
			Db.Exec(`INSERT INTO User (FirstName, LastName, CompanyID, Password, AccessLevel, Email)
					 VALUES ('Joe', 'Blow', 2, 'joe', 10, 'joe@blow.com');"`)
			Db.Exec(`INSERT INTO User (FirstName, LastName, CompanyID, Password, AccessLevel, Email)
					 VALUES ('Brad' , 'McCormack', 2, 'brad', 9, 'bradmccormack100@gmail.com');`)
 			//Set up some default settings
			Db.Exec("INSERT INTO Settings (UserID, MapAPI, Interpolate, SnaptoRoad, CameraPanTrigger) VALUES (1, 'Google Maps', 0, 0, 1);")
			Db.Exec("INSERT INTO Settings (UserID, MapAPI, Interpolate, SnaptoRoad, CameraPanTrigger) VALUES (2, 'Google Maps', 0, 0, 1);")
			Db.Exec("INSERT INTO Settings (UserID, MapAPI, Interpolate, SnaptoRoad, CameraPanTrigger) VALUES (3, 'Google Maps', 0, 0, 1);")
			Db.Exec("INSERT INTO Settings (UserID, MapAPI, Interpolate, SnaptoRoad, CameraPanTrigger) VALUES (4, 'Google Maps', 0, 0, 1);")

			//Set up some default companies and path to logo
			Db.Exec("INSERT INTO Company (Name, MaxUsers, Expiry, LogoPath) VALUES ('myClubLink' , 1, '2100-01-20 12:00:00', 'img/mcl_logo.png');")
			Db.Exec("INSERT INTO Company (Name, MaxUsers, Expiry, LogoPath) VALUES ('Sussex Inlet RSL Group', 5, '2015-06-6 12:00:00', 'img/sussex_logo.PNG');")

			//Note a company must have a company settings record
			Db.Exec(`INSERT INTO CompanySettings (CompanyID, RadioCommunication, DataCommunication, SecurityRemoteAdmin,
					SecurityConsoleAccess, SecurityAdminPasswordReset, MobileSmartPhoneAccess, MinZoom, MaxZoom, HistoricalmapsKmMin, ClubBoundaryKM)
					VALUES(1, 1, 1, 0, 0, 0, 0, 1, 10, 10, 100);`)

			Db.Exec(`INSERT INTO CompanySettings (CompanyID, RadioCommunication, DataCommunication, SecurityRemoteAdmin,
					SecurityConsoleAccess, SecurityAdminPasswordReset, MobileSmartPhoneAccess, MinZoom, MaxZoom, HistoricalmapsKmMin, ClubBoundaryKM)
					VALUES(2, 1, 1, 0, 0, 0, 0, 1, 10, 10, 100);`)
		},
	}

	GeoDbSchema := []DbFunc{
		func(Db *sql.DB) {

			//point of interest such as street
			Db.Exec(`CREATE TABLE POI (
					ID INTEGER NOT NULL PRIMARY KEY,
					Name TEXT NOT NULL);`)

			Db.Exec(`CREATE TABLE LatLong (
					 ID INTEGER NOT NULL PRIMARY KEY,
					 POIID INTEGER NOT NULL,
					 LAT TEXT NOT NULL,
					 LONG TEXT NOT NULL,
					 Name TEXT);`)

			Db.Exec("PRAGMA foreign_keys=ON;")
		},
	}

	DatabasesChanges := map[string][]DbFunc{
		"backend.db":  BackendDbSchema,
		"license.key": LicenseDbSchema,
		"geodata.db":  GeoDbSchema,
	}
	for k, DataBaseChanges := range DatabasesChanges {

		DbHandle, err := sql.Open("sqlite3", k)
		if err != nil {
			fmt.Printf("%s didn't exist. will be created", k)
			_, err := os.Create("./backend.db")
			if err != nil {
				log.Fatal("Cannot create %s!\n", k)
			}
		}

		var user_version int
		result := DbHandle.QueryRow("PRAGMA USER_VERSION").Scan(&user_version)
		if result != nil {
			log.Fatal("Cannot get user version\n")
		}

		for i := user_version; i < len(DataBaseChanges); i++ {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Schema changes failed for version %d on database %s\n", i, k)
					DbHandle.Exec("ROLLBACK TRANSACTION")
				}
			}()

			DbHandle.Exec("BEGIN EXCLUSIVE TRANSACTION")
			fmt.Printf("\n%s - executing schema version %d\n", k, i)
			Fn := DataBaseChanges[i]
			Fn(DbHandle)
			DbHandle.Exec("COMMIT TRANSACTION")
			user_version++
		}

		DbHandle.Exec("PRAGMA USER_VERSION=" + strconv.Itoa(user_version))
		DbHandle.Close()

	}
}

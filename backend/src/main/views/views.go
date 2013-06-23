package views

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
)

var Db *sql.DB

func ViewInvalid(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Invalid view", 403)
}

func ViewLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("In ViewLogin")
	var err error
	t := template.New("Login")
	t, err = template.ParseFiles("templates/login.html")
	if err != nil {
		fmt.Printf("Failed to parse the template file!\n")
		return
	}
	t.Execute(w, nil)
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

	t.Execute(w, settings) //second param is the data interface. It's the equiv of Bondi's AddDataSet I believe'

}

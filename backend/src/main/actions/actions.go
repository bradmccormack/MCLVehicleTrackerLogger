package actions

import (
	"fmt"
	"net/http"
)

func ActionInvalid(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Invalid Action", 403)
}

func ActionLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	fmt.Printf("In Action Login")
}

func ActionSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	fmt.Printf("In Action Settings")
}

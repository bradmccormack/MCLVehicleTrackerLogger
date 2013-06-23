package views

import (
	"fmt"
	"net/http"
)

func ActionLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	fmt.Printf("In Action Login")
}

func ActionSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	fmt.Printf("In Action Settings")
}

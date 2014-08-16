package socket

import (
	"../types"
	"../utility"
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"time"
)

//the string key will be a hash of the username and ip
var connections map[[32]byte]*types.ClientSocket

//exported
var VehicleChannel = make(chan types.Record, 100)

func init() {
	connections = make(map[[32]byte]*types.ClientSocket)
}

func WebSocketInit(w http.ResponseWriter, r *http.Request, cookiejar *sessions.CookieStore) {

	fmt.Printf("\n In Handlewebsocketinit \n")
	session, _ := cookiejar.Get(r, "data")

	var user types.User = session.Values["User"].(types.User)

	//fmt.Printf("Username is %s %s\n", user.Firstname, user.Lastname)
	//fmt.Printf("Web socket requested from %s\n", utility.GetIpAddress(r))

	if r.Method != "GET" {
		fmt.Printf("GET method request for socket. Not allowed\n")
		http.Error(w, "Method not allowed", 405)
		return
	}

	/* TODO fix this
		if r.Header.Get("Origin") + *addr != "http://" + r.Host {
	    		http.Error(w, "Origin not allowed", 403)
	    		return
	    	}
	*/

	var ip string = utility.GetIpAddress(r)

	//hash the incoming ip and username
	var buffer bytes.Buffer
	buffer.WriteString(ip)
	buffer.WriteString(user.FirstName)
	buffer.WriteString(user.LastName)
	fmt.Printf("WebSocket -> the ip is %s the user is %s\n", ip, user.FirstName)

	var hash = sha256.Sum256(buffer.Bytes())
	//fmt.Printf("The hash in web socket is %b\n", hash)

	if _, exists := connections[hash]; exists {
		fmt.Printf("Connection existed .. closing \n")
		connections[hash].Websocket.Close()
	} else {
		fmt.Printf("New connection created")
	}

	connection, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		fmt.Printf("Not a websocket handshake \n")
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		fmt.Printf("Something bad happened - %s", err)
		log.Println(err)
		return
	}

	//create new connection ready to go
	connections[hash] = new(types.ClientSocket)

	fmt.Printf("About to set the connection\n")
	connections[hash].Websocket = connection
	fmt.Printf("Amount of web socket connections is %d\n", len(connections))
}
func WebSocketClose(hash [32]byte) {
	if connections[hash] != nil {
		connections[hash].Websocket.Close()
	}

}

//this should have a buffered channel that will block the sender when it is full every 1 second it will read from the channels and send off shit to the webservers
//when it sends shit off it should do so using goroutines so they don't block
func Monitor() {
	fmt.Printf("\nin Monitor")

	for {
		fmt.Printf("\n Sleeping in Monitor")

		starttime := time.Now()

		for time.Since(starttime) < time.Second {
			if len(connections) == 0 {
				fmt.Printf("\nNo webclients listening via websocket.. not reporting")
			} else {
				R := <-VehicleChannel
				go UpdateClient(R.GPS, R.Diagnostic)
			}
			//read on the channel
		}

	}
}

func UpdateClient(entry *types.GPSRecord, diagnostic *types.DiagnosticRecord) {

	//fmt.Printf("Responding to %d listening clients\n", len(connections))
	for _, client := range connections {
		//get a websocket writer

		wswriter, _ := client.Websocket.NextWriter(websocket.TextMessage)

		if wswriter != nil {
			io.WriteString(wswriter, types.JSONResponse{"Entry": entry, "Diagnostic": diagnostic}.String())
		} else {
			fmt.Printf("No ws writer available\n") //this web socket was abruptly closed so we need to close that client
			client.Websocket.Close()
		}

	}
}

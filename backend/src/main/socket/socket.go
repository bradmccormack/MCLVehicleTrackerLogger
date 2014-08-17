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

func init() {
	connections = make(map[[32]byte]*types.ClientSocket)
}

//sends websocket pings to client so it doesn't close the session.. the built in ponghandler is left in place. I don't care
func heartBeat(c *websocket.Conn) {

	ticker := time.NewTicker(5 * time.Second)
	defer func() {
		ticker.Stop()
		c.Close()
	}()
	for {
		select {

		case <-ticker.C:
			_ = c.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(5*time.Second))

		}

	}
}

func WebSocketInit(w http.ResponseWriter, r *http.Request, cookiejar *sessions.CookieStore) {

	session, _ := cookiejar.Get(r, "data")

	var user types.User = session.Values["User"].(types.User)

	if r.Method != "GET" {
		fmt.Printf("GET method request for socket. Not allowed\n")
		http.Error(w, "Method not allowed", 405)
		return
	}

	var ip string = utility.GetIpAddress(r)

	//hash the incoming ip and username
	var buffer bytes.Buffer
	buffer.WriteString(ip)
	buffer.WriteString(user.FirstName)
	buffer.WriteString(user.LastName)

	var hash = sha256.Sum256(buffer.Bytes())

	if _, exists := connections[hash]; exists {
		connections[hash].Websocket.Close()
	}

	connection, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		fmt.Printf("Not a websocket handshake \n")
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}

	//create new connection ready to go
	connections[hash] = new(types.ClientSocket)

	connections[hash].Websocket = connection
	go heartBeat(connection)

	fmt.Printf("\nAmount of web socket connections is %d\n", len(connections))
}

func WebSocketClose(hash [32]byte) {
	if connections[hash] != nil {
		connections[hash].Websocket.Close()
	}
	delete(connections, hash)

}

//read only channels,
func Monitor(DataChannel <-chan types.Record, CommandChanel <-chan int) {

	for {

		starttime := time.Now()

		for time.Since(starttime) < time.Second {

			//select from first available channel ipc
			select {
			//keep slurping records from the bufered channel and farm them out to UpdateClient as a goroutine
			case data := <-DataChannel:
				go UpdateClient(data.GPS, data.Diagnostic)
			case command := <-CommandChanel:
				switch command {
				case (types.Command_Quit):
					return

				}
			}

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

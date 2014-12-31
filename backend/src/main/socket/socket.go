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

//The string key will be a hash of the username and ip
var connections map[[32]byte]*types.ClientSocket

func init() {
	connections = make(map[[32]byte]*types.ClientSocket)
}

//Sends websocket pings to client so it doesn't close the session.. the built in pong handler is left in place.
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
		fmt.Printf("\nGET method request for socket. Not allowed\n")
		http.Error(w, "Method not allowed", 405)
		return
	}

	var ip string = utility.GetIpAddress(r)

	//Hash the incoming ip and username
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
		fmt.Printf("\nNot a websocket handshake \n")
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}

	//Create new connection ready to go
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

//Read only channels,
func Monitor(DataChannel <-chan types.Record, CommandChannel <-chan int32) {

	for {
		starttime := time.Now()
		for time.Since(starttime) < time.Second {

			//select from first available channel ipc - note this blocks until there is data in one of the channels
			select {
			//keep slurping records from the bufered channel and farm them out to UpdateClient as a goroutine
			case data := <-DataChannel:
				go UpdateClient(data.GPS, data.Diagnostic)
			case command := <-CommandChannel:
				switch command {
				case (types.COMMAND_QUIT):
					fmt.Printf("\nQuit command this monitor should be exiting")
					return

				}
			}

		}

	}
}

func UpdateClient(entry *types.GPSRecord, diagnostic *types.DiagnosticRecord) {

	for _, client := range connections {
		//get a websocket writer

		wswriter, _ := client.Websocket.NextWriter(websocket.TextMessage)

		if wswriter != nil {
			io.WriteString(wswriter, types.JSONResponse{"Entry": entry, "Diagnostic": diagnostic}.String())
		} else {
			client.Websocket.Close()
		}

	}
}

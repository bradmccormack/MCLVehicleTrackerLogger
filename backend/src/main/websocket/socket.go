package websocket

import (
	"net/http"
	"strings"
)

/* Need to intermittently send a ping message http://tools.ietf.org/html/rfc6455#section-5.5.2 to keep the connection alive
func pingWebSockets() {

}

func handleWebSocketInit(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("\n In Handlewebsocketinit \n")
	session, _ := store.Get(r, "data")
	var user User = session.Values["User"].(User)

	//fmt.Printf("Username is %s %s\n", user.Firstname, user.Lastname)
	//fmt.Printf("Web socket requested from %s\n", utility.GetIpAddress(r))


	if r.Method != "GET" {
		fmt.Printf("GET method request for socket. Not allowed\n")
		http.Error(w, "Method not allowed", 405)
		return
	}

	/*
	if r.Header.Get("Origin") + *addr != "http://" + r.Host {
    		http.Error(w, "Origin not allowed", 403)
    		return
    	}
	*/


	var ip string = utility.GetIpAddress(r)

	//hash the incoming ip and username
	var buffer bytes.Buffer
	buffer.WriteString(ip)
	buffer.WriteString(user.Firstname)
	buffer.WriteString(user.Lastname)
	fmt.Printf("WebSocket -> the ip is %s the user is %s %s\n", ip, user.Firstname, user.Lastname)

	var hash = sha256.Sum256(buffer.Bytes())
	//fmt.Printf("The hash in web socket is %b\n", hash)

	/* //this seems to be superfluous
	if _, exists := connections[hash]; exists {
		fmt.Printf("Connection existed .. closing \n")
		connections[hash].websocket.Close()
        } else {
		fmt.Printf("New connection created");
	}
	*/

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
	connections[hash] = new(ClientSocket)

	connections[hash].websocket = connection
	fmt.Printf("Amount of web socket connections is %d\n", len(connections))

}
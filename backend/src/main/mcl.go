package main

import (
	"./dao"
	"./http"
	"./socket"
	"./types"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"

	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

//set the domain based upon the path the executable was run from
var domain string = "dev.myclublink.com.au"
var service = flag.String("service", ":6969", "tcp port to bind to")

var lat = flag.String("lat", "", "latitude to test")
var lon = flag.String("long", "", "longitude to test")

func init() {
	gob.Register(types.User{})
	gob.Register(types.Company{})
	gob.Register(types.Settings{})
}

func main() {

	var addr = flag.String("addr", ":8080", "http(s) service address")
	flag.Parse()

	if _, err := os.Stat("backend.db"); err != nil {
		log.Fatal("\nError: ", err)
	}

	if _, err := os.Stat("license.key"); err != nil {
		log.Fatal("\nError: ", err)
	}

	dao.Open()

	//TODO remove this shit.. its just testing street name stuff
	if *lat != "" && *lon != "" {
		//func ParseFloat(s string, bitSize int) (f float64, err error)
		latf, _ := strconv.ParseFloat(*lat, 64)
		longf, _ := strconv.ParseFloat(*lon, 64)
		street := dao.GetStreetName(latf, longf)
		fmt.Printf("\nStreet is %s\n", street)
		os.Exit(0)
	}

	defer dao.Close()

	//Socket related channels
	WSDataChannel := make(chan types.Record, 100) //buffered
	WSCommandChannel := make(chan int32)

	NetworkChannel := make(chan int32, 1) //buffered

	go socket.Monitor(WSDataChannel, WSCommandChannel)                    //start a websocket dude to arbitrate websockets
	go connectionManager(NetworkChannel, WSCommandChannel, WSDataChannel) //connection manager to handle (re)connects
	go http.HttpRouter(addr)

	//kick off initial connection - send COMMAND_RECONNECT on the Network Channel
	NetworkChannel <- types.COMMAND_RECONNECT

	//wait forever nicely
	select {}

}

func connectionManager(NetworkChannel chan int32, WSCommandChannel chan<- int32, WSDataChannel chan<- types.Record) {

	//select from first available channel ipc - note this blocks until there is data in one of the channels
	select {
	//keep slurping records from the bufered channel and farm them out to UpdateClient as a goroutine
	case msg := <-NetworkChannel:
		switch msg {
		case types.COMMAND_RECONNECT:

			lnk, err := net.Listen("tcp", *service)
			if err != nil {
				fmt.Printf("\nFailed to get tcp listener - %s", err.Error())
				os.Exit(1)
			}
			fmt.Printf("\nListening on TCP Port %s", *service)

			for {
				//this blocks until there is a connection
				tcpcon, err := lnk.Accept()

				fmt.Printf("\nLink Accepted - Receiving packets from Vehicle")
				if err != nil {
					fmt.Printf("\nFailed to create tcp connection - %s", err)
					os.Exit(1)
				}
				go handleClient(tcpcon.(*net.TCPConn), WSDataChannel, NetworkChannel)
			}
		}
	}
}

func handleClient(conn *net.TCPConn, WSDataChannel chan<- types.Record, NetworkChannel chan<- int32) {

	//defer anonymous func to handle panics - most likely panicking from garbage that was to be parsed.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from a panic \n", r)
		}
	}()

	var buff = make([]byte, 512)
	var incomingpacket types.Packet

	var R types.Record
	R.GPS = new(types.GPSRecord)
	R.Diagnostic = new(types.DiagnosticRecord)

	var n int
	var err error
	for {
		n, err = conn.Read(buff)

		if err != nil {
			fmt.Printf("\nError occured - %s, will recreate the connection.", err.Error())
			NetworkChannel <- types.COMMAND_RECONNECT
			return
		}

		//lets unmarshal those JSON bytes into the map https://groups.google.com/forum/#!topic/golang-nuts/77HJlZhWXpk  note to slice properly otherwise it chockes on trying to decode the full buffer
		err := json.Unmarshal(buff[:n], &incomingpacket)
		if err != nil {
			fmt.Printf("Failed to decode the JSON bytes -%s\n", err.Error())
		}

		diagnosticfields := strings.Split(string(incomingpacket["diagnostics"]), ",")
		if len(diagnosticfields) != 4 {
			fmt.Printf("Error. Diagnostic fields length is incorrect. Is %d should be %d", len(diagnosticfields), 4)
			fmt.Printf("The source string was %s\n", string(incomingpacket["diagnostics"]))
		}

		gpsfields := strings.Split(string(incomingpacket["sentence"]), ",")

		if len(gpsfields) != 7 {
			fmt.Printf("Error. GPS fields length is incorrect. Is %d should be %d\n", len(gpsfields), 7)
			fmt.Printf("The source string was %s\n", string(incomingpacket["sentence"]))
			continue
		}

		R.Diagnostic.CPUTemp, _ = strconv.ParseFloat(diagnosticfields[0][2:], 32)
		R.Diagnostic.CPUVolt, _ = strconv.ParseFloat(diagnosticfields[1][2:], 32)
		R.Diagnostic.CPUFreq, _ = strconv.ParseFloat(diagnosticfields[2][2:], 32)
		R.Diagnostic.MemFree, _ = strconv.ParseUint(diagnosticfields[3][2:], 10, 64)

		R.GPS.Message = gpsfields[0][1:]
		R.GPS.Latitude = gpsfields[0][2:]
		R.GPS.Longitude = gpsfields[1]
		R.GPS.Speed, _ = strconv.ParseFloat(gpsfields[2][1:], 32)
		R.GPS.Heading, _ = strconv.ParseFloat(gpsfields[3][1:], 32)
		R.GPS.Date, _ = time.Parse(time.RFC3339, gpsfields[4][1:])
		R.GPS.Fix = gpsfields[5][1:] == "true"
		R.GPS.ID = gpsfields[6][1:]

		if string(incomingpacket["sentence"][0:1]) != "T" {
			go func() {
				dao.SavePacket(R.GPS, R.Diagnostic)
			}()
		}

		WSDataChannel <- R
		conn.Write([]byte("OK\n"))
	}
	return
}

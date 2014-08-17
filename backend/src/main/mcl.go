package main

import (
	"./http"
	"./socket"
	"./types"
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
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
var Db *sql.DB

func init() {
	gob.Register(types.User{})
	gob.Register(types.Company{})
	gob.Register(types.Settings{})
}

func main() {

	var addr = flag.String("addr", ":8080", "http(s) service address")
	flag.Parse()

	var err error

	if _, err := os.Stat("backend.db"); err != nil {
		log.Fatal("\nError: ", err)
	}

	if _, err := os.Stat("license.key"); err != nil {
		log.Fatal("\nError: ", err)
	}

	Db, err = sql.Open("sqlite3", "backend.db")
	if err != nil {
		fmt.Printf("Cannot open database backend.db . Exiting\n")
		os.Exit(1)
	}

	//try to attach the license.key if it opens then close it and attach
	LDb, err := sql.Open("sqlite3", "license.key")
	if err != nil {
		fmt.Printf("\nCannot open database license.key . Exiting\n")
		os.Exit(1)
	} else {
		fmt.Printf("\nLicense.key opened correctly")
	}
	LDb.Close()

	defer Db.Close()

	//Socket related channels
	WSDataChannel := make(chan types.Record, 100) //buffered
	WSCommandChannel := make(chan int32)

	NetworkChannel := make(chan int32, 1) //buffered

	go socket.Monitor(WSDataChannel, WSCommandChannel)                    //start a websocket dude to arbitrate websockets
	go connectionManager(NetworkChannel, WSCommandChannel, WSDataChannel) //connection manager to handle (re)connects
	go http.HttpRouter(addr, Db)

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

func logEntry(entry *types.GPSRecord, diagnostic *types.DiagnosticRecord) {

	_, err := Db.Exec("BEGIN TRANSACTION")
	_, err = Db.Exec("INSERT INTO GPSRecords (Message, Latitude, Longitude, Speed, Heading, Fix, DateTime, BusID) VALUES ( ? , ?, ? , ? , ? ,? ,? , ?)",
		entry.Message,
		entry.Latitude,
		entry.Longitude,
		entry.Speed,
		entry.Heading,
		entry.Fix,
		entry.Date,
		entry.ID)

	_, err = Db.Exec("INSERT INTO DiagnosticRecords (CPUTemperature, CPUVoltage, CPUFrequency, MemoryFree) VALUES (?, ?, ?, ?)",
		diagnostic.CPUTemp,
		diagnostic.CPUVolt,
		diagnostic.CPUFreq,
		diagnostic.MemFree)

	Db.Exec("COMMIT TRANSACTION")
	if err != nil {
		fmt.Printf("Failed to insert row %s", err)
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

		//fmt.Printf("\nSentence was %s", string(incomingpacket["sentence"]))
		//fmt.Printf("\nDiagnostic data was %s", string(incomingpacket["diagnostics"]))

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

		if string(incomingpacket["sentence"][0:1]) != "T" {
			go logEntry(R.GPS, R.Diagnostic)
		}

		WSDataChannel <- R
		conn.Write([]byte("OK\n"))
	}
	return
}

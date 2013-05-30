package main

import (
    "fmt"
    "net"
    "os"
    "time"
)


func handleConnection() {
}


func main() {

    service := ":215";
    udpAddr, err := net.ResolveUDPAddr("udp4", service);
    if err != nil {
        fmt.Printf("Failed to resolve UDP address")
        os.Exit(1)
    }
    fmt.Printf("Listening on UDP Port %s", service)


    con, err := net.ListenUDP("udp",udpAddr)
    if(err != nil) {
        fmt.Printf("Failed to create udp connection " + err.Error())
        os.Exit(1)
    }
    defer con.Close();

    //palm off reading and writing to a go routine
    go func(conn *net.UDPConn) {
            var buff [512]byte

            /*GPS input will look like
             P L-34.50108,150.81094,S0.00,H147.2,D2013-05-26,T11:00:18,Ftrue,ISUS01
            P= This is a position message, L is lattitude and longitude seperated by comma, S is Speed in KM/H H is heading in degrees
            F is boolean fix or not D will be date and T for time  and I will be ID of logger
            */

            _, addr, err := conn.ReadFromUDP(buff[:])
            if err != nil {
                fmt.Printf("Error reading from UDP")
            }

            daytime := time.Now().String()
            conn.WriteToUDP([]byte(daytime), addr)
            fmt.Printf("Responded")

    }(con)



}

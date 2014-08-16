package types

import (
	"github.com/gorilla/websocket"
	"time"
)

type Packet map[string]string
type Response map[string]interface{}

type GPSRecord struct {
	Latitude  string
	Longitude string
	Message   string
	Speed     float64
	Heading   float64
	Fix       bool
	Date      time.Time
	ID        string
}

type DiagnosticRecord struct {
	CPUTemp float64
	CPUVolt float64
	CPUFreq float64
	MemFree uint64
}

type Company struct {
	Name     string
	Maxusers int
	Expiry   string
	LogoPath string
}

type User struct {
	ID          int
	Firstname   string
	Lastname    string
	Password    string
	Accesslevel int
	Email       string
}

type Settings struct {
	MapAPI                     string
	Interpolate                int
	SnaptoRoad                 int
	CameraPanTrigger           int
	RadioCommunication         int
	DataCommunication          int
	SecurityRemoteAdmin        int
	SecurityConsoleAccess      int
	SecurityAdminPasswordReset int
	MobileSmartPhoneAccess     int
	MobileShowBusLocation      int
	MinZoom                    int
	MaxZoom                    int
	HistoricalmapsKmMin        int
	ClubBoundaryKM             int
}

type ClientSocket struct {
	Websocket    *websocket.Conn
	Ip, Username string
}

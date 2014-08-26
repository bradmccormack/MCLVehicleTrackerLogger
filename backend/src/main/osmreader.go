package main

import (
	"database/sql"
	"encoding/xml"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type Node struct {
	Id   string `xml:"id,attr"`
	Lat  string `xml:"lat,attr"`
	Lon  string `xml:"lon,attr"`
	Name string
}

type Way struct {
	Id      string `xml : "id,attr"`
	Name    string
	NodeIds []string `xml: "nd>ref"`
}

const (
	NODE int32 = 0
	WAY  int32 = 0
)

func main() {

	nodeMap := make(map[string]Node)
	wayMap := make(map[string]Way)

	dbname := flag.String("database", "geodata.db", "Name of database to insert records into")
	file := flag.String("input", "data.osm", "Name of osm(xml) file to import rows from")
	flag.Parse()

	var Db *sql.DB
	var err error
	Db, err = sql.Open("sqlite3", *dbname)
	if err != nil {
		fmt.Printf("Cannot open database backend.db . Exiting\n")
		os.Exit(1)
	}

	fmt.Printf("\nReading %s ...\n", *file)

	xmlFile, err := os.Open(*file)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	var elm string
	var ParentID string
	var ParentType int32

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			elm = se.Name.Local

			switch elm {
			case "node":

				var n Node
				decoder.DecodeElement(&n, &se)

				ParentID = n.Id
				ParentType = NODE
				nodeMap[n.Id] = n

			case "way":

				w := Way{
					Id:      se.Attr[0].Value,
					NodeIds: make([]string, 0), //slice of strings
				}

				ParentID = w.Id
				ParentType = WAY
				wayMap[w.Id] = w

			case "nd":
				//find all node refs and add them to NodeIds
				nodeRef := se.Attr[0].Value

				WayParent := wayMap[ParentID]
				WayParent.NodeIds = append(WayParent.NodeIds, nodeRef)
				wayMap[ParentID] = WayParent
				//wayMap[ParentID].NodeIds = append(wayMap[ParentID].NodeIds, nodeRef) //is illegal https://code.google.com/p/go/issues/detail?id=3117

			case "tag":
				//this needs checking... might need to use decodeElement
				if se.Attr[0].Value == "name" {
					if ParentType == NODE {
						NodeParent := nodeMap[ParentID]
						NodeParent.Name = se.Attr[1].Value
						nodeMap[ParentID] = NodeParent //rare but I can see intstances such as   <tag k="name" v="Nowra Community Hospital"/>
					}
					if ParentType == WAY {

						WayParent := wayMap[ParentID]
						WayParent.Name = se.Attr[1].Value
						wayMap[ParentID] = WayParent //common looks like   <tag k="name" v="Kinghorne Street"/>
					}
				}
			}

		default:
		}

	}

	fmt.Printf("\nNumber of Places of Interest (streets etc) %d", len(wayMap))
	fmt.Printf("\nNumber of Nodes %d\n", len(nodeMap))

	Db.Exec("DROP INDEX IDX_latlng")

	Db.Exec("BEGIN TRANSACTION")
	for _, poi := range wayMap {
		Db.Exec("INSERT INTO POI (ID, Name) VALUES(?, ?)", poi.Id, poi.Name)
		for _, nodeID := range poi.NodeIds {
			node := nodeMap[nodeID]

			Db.Exec("INSERT INTO LatLong (ID, POIID, Name, Lat, Long) VALUES(?, ?, ?, ?, ?)", node.Id, poi.Id, node.Name, node.Lat, node.Lon)
		}
	}
	Db.Exec("CREATE INDEX IDX_latlng ON LatLong(Lat,Long, POID")
	Db.Exec("COMMIT TRANSACTION")

}

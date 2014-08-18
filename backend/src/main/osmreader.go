package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
)

type Dictionary struct {
	XMLName xml.Name `xml:  "osm"`
	Nodes   []node   `xml: "osm>node`
}

type node struct {
	Id        int    `xml : "id,attr"`
	Visible   string `xml : "visible,attr"`
	Version   int    `xml : "version,attr"`
	Changeset string `xml : "changeset,attr"`
	Timestamp string `xml : "timestamp, attr"`
	User      string `xml : "user,attr"`
	Uid       string `xml : "uid,attr"`
	Lat       string `xml : "lat,attr"`
	Lon       string `xml : "lon,attr"`
}

/*
type way struct {
	id     string   `xml : "id,attr"`
	nodeid []string `xml: "nd>ref"`
}
*/

func main() {

	nodeMap := make(map[int]node)
	//wayMap := make(map[string]way)

	//db := flag.String("database", "address.db", "Name of database to insert records into")

	file := flag.String("input", "data.osm", "Name of osm(xml) file to import rows from")
	flag.Parse()

	fmt.Printf("\nReading %s ...\n", *file)

	xmlFile, err := os.Open(*file)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	var inElement string
	var i int = 0
	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			// If we just read a StartElement token
			inElement = se.Name.Local

			switch inElement {
			case "node":
				var n node
				decoder.DecodeElement(&n, &se)
				if i == 0 {
					fmt.Printf("ID is %d", n.Id)
					i++
				}

				/*
					if i == 0 {
						for k, v := range se.Attr {
							fmt.Printf("\nAttr is %s value is %s", k, v)
						}
						i++
					}
				*/

				nodeMap[n.Id] = n
				/*
					case "way":
						var w way
						decoder.DecodeElement(&w, &se)
						wayMap[w.id] = w
				*/
			}

		default:
		}

	}
	//fmt.Printf("\nNumber of waynodes is %d", len(wayMap))
	fmt.Printf("\nNumber of nodes is %d", len(nodeMap))
}

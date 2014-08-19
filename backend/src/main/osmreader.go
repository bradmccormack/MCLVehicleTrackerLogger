package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
)

type Node struct {
	Id   string `xml : "attr"`
	Lat  string `xml : "attr"`
	Lon  string `xml : "attr"`
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
				//no need to use DecodeElement for now.. this is likely probably faster !
				n := Node{
					Id:  se.Attr[0].Value,
					Lat: se.Attr[7].Value,
					Lon: se.Attr[8].Value,
				}
				ParentID = n.Id
				ParentType = NODE
				nodeMap[n.Id] = n

			case "way":
				w := Way{
					Id:      se.Attr[0].Value,
					NodeIds: make([]string, 100), //slice of strings .. not sure about 100 here..
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
				if se.Attr[0].Value == "name" {
					if ParentType == NODE {
						NodeParent := nodeMap[ParentID]
						NodeParent.Name = se.Attr[1].Value
						nodeMap[ParentID] = NodeParent //rare but I can see intstances such as   <tag k="name" v="Nowra Community Hospital"/>
					} else if ParentType == WAY {
						fmt.Printf("\n Looking at a tag element, value is %s, wayparentid is %s\n", se.Attr[1].Value, ParentID)

						WayParent := wayMap[ParentID]
						WayParent.Name = se.Attr[1].Value
						wayMap[ParentID] = WayParent //common looks like   <tag k="name" v="Kinghorne Street"/>
					}

				}
			}

		default:
		}

	}
	fmt.Printf("\nNumber of nodes is %d", len(nodeMap))
	WayTest := wayMap["283933425"]
	fmt.Printf("\nId is %s", WayTest.Id)
	fmt.Printf("\nName is %s", WayTest.Name) //failed

	fmt.Printf("Nodes belonging to this are\n")
	for _, v := range WayTest.NodeIds {
		fmt.Printf("\nID is %s", v)
	}
	fmt.Printf("\n")

}

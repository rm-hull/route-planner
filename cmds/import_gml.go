package cmds

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"

	"github.com/rm-hull/route-planner/models"
)

func ImportGmlData(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)

	for {
		token, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatalf("error reading token: %v", err)
		}

		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Local == "featureMember" {
				var feature models.FeatureMember
				err := decoder.DecodeElement(&feature, &se)
				if err != nil {
					log.Fatalf("error decoding element: %v", err)
				}

				if feature.RoadLink != nil {
					fmt.Printf("RoadLink ID: %s\n", feature.RoadLink.ID)
				} else if feature.RoadNode != nil {
					fmt.Printf("RoadNode ID: %s\n", feature.RoadNode.ID)
				} else if feature.MotorwayJunction != nil {
					fmt.Printf("MotorwayJunction ID: %s\n", feature.MotorwayJunction.ID)
				} else {
					fmt.Printf("Unhandled: %v\n", se)
				}
			}
		}
	}
}

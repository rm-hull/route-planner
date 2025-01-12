package models

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// Root element representing the collection of feature members
type FeatureCollection struct {
	XMLName        xml.Name        `xml:"FeatureCollection"`
	FeatureMembers []FeatureMember `xml:"featureMember"`
}

// Represents a single feature member (RoadLink, RoadNode, Motorway Junction)
type FeatureMember struct {
	XMLName          xml.Name          `xml:"featureMember"`
	RoadLink         *RoadLink         `xml:"RoadLink,omitempty"`
	RoadNode         *RoadNode         `xml:"RoadNode,omitempty"`
	MotorwayJunction *MotorwayJunction `xml:"MotorwayJunction,omitempty"`
}

// RoadLink struct
type RoadLink struct {
	XMLName                  xml.Name  `xml:"RoadLink"`
	ID                       string    `xml:"id,attr"`
	CentrelineGeometry       Geometry  `xml:"centrelineGeometry>LineString"`
	EndNode                  NodeRef   `xml:"endNode"`
	StartNode                NodeRef   `xml:"startNode"`
	RoadClassification       CodeValue `xml:"roadClassification"`
	RoadFunction             CodeValue `xml:"roadFunction"`
	FormOfWay                CodeValue `xml:"formOfWay"`
	RoadClassificationNumber *string   `xml:"roadClassificationNumber,omitempty"`
	Name1                    *string   `xml:"name1,omitempty"`
	Length                   Length    `xml:"length"`
	Loop                     bool      `xml:"loop"`
	PrimaryRoute             bool      `xml:"primaryRoute"`
	TrunkRoad                bool      `xml:"trunkRoad"`
}

// RoadNode struct
type RoadNode struct {
	XMLName        xml.Name  `xml:"RoadNode"`
	ID             string    `xml:"id,attr"`
	Geometry       Point     `xml:"geometry>Point"`
	FormOfRoadNode CodeValue `xml:"formOfRoadNode"`
}

type MotorwayJunction struct {
	XMLName        xml.Name `xml:"MotorwayJunction"`
	ID             string   `xml:"id,attr"`
	Geometry       Point    `xml:"geometry>Point"`
	JunctionNumber string   `xml:"junctionNumber"`
}

// Geometry representation for LineString
type Geometry struct {
	ID           string `xml:"id,attr"`
	SRSName      string `xml:"srsName,attr"`
	SRSDimension int    `xml:"srsDimension,attr"`
	PosList      string `xml:"posList"`
}

func (g Geometry) AsLineString() string {
	values := strings.Fields(g.PosList)
	if len(values)%g.SRSDimension != 0 {
		panic(fmt.Errorf("coordinates length (%d) is not divisible by srsDimension (%d)", len(values), g.SRSDimension))
	}

	var coordinates []string
	for i := 0; i < len(values); i += g.SRSDimension {
		var tuple []string
		for j := 0; j < g.SRSDimension; j++ {
			tuple = append(tuple, values[i+j])
		}
		coordinates = append(coordinates, strings.Join(tuple, " "))
	}

	return "LINESTRING(" + strings.Join(coordinates, ",") + ")"
}

// Represents a point geometry
type Point struct {
	ID           string `xml:"gml:id,attr"`
	SRSName      string `xml:"srsName,attr"`
	SRSDimension int    `xml:"srsDimension,attr"`
	Position     string `xml:"pos"`
}

func (p Point) AsPoint() any {
	var sb strings.Builder
	sb.WriteString("POINT(")
	sb.WriteString(p.Position)
	sb.WriteString(")")
	return sb.String()
}

// Reference to a start or end node
type NodeRef struct {
	Href string `xml:"href,attr"`
}

func (node *NodeRef) Ref() string {
	return strings.TrimPrefix(node.Href, "#")
}

// Represents a classification or function code with a codeSpace attribute
type CodeValue struct {
	CodeSpace string `xml:"codeSpace,attr"`
	Value     string `xml:",chardata"`
}

// Represents a length with unit of measure
type Length struct {
	Unit  string  `xml:"uom,attr"`
	Value float64 `xml:",chardata"`
}

func (l Length) ConvertTo(newUnit string) float64 {
	if l.Unit == newUnit {
		return l.Value
	}
	panic("unimplemented")
}

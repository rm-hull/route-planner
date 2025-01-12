package models

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

// Root element representing the collection of feature members
type FeatureCollection struct {
	XMLName        xml.Name        `xml:"FeatureCollection"`
	FeatureMembers []FeatureMember `xml:"featureMember"`
}

// Represents a single feature member (RoadLink or RoadNode)
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
	RoadNameTOID             *string   `xml:"roadNameTOID,omitempty"`
	RoadNumberTOID           *string   `xml:"roadNumberTOID,omitempty"`
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

type CoordinateProvider interface {
	GetPosList() string
	GetSRSDimension() int
}

func ParseCoordinates(provider CoordinateProvider) ([][]float64, error) {
	// Split the space-separated string
	parts := strings.Fields(provider.GetPosList())

	// Validate that the length of parts is divisible by srsDimension
	if len(parts)%provider.GetSRSDimension() != 0 {
		return nil, fmt.Errorf("coordinates length (%d) is not divisible by srsDimension (%d)", len(parts), provider.GetSRSDimension())
	}

	// Convert to [][]float64
	var coordinates [][]float64
	for i := 0; i < len(parts); i += provider.GetSRSDimension() {
		var tuple []float64
		for j := 0; j < provider.GetSRSDimension(); j++ {
			val, err := strconv.ParseFloat(parts[i+j], 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse coordinates: %w", err)
			}
			tuple = append(tuple, val)
		}
		coordinates = append(coordinates, tuple)
	}

	return coordinates, nil
}

// Geometry representation for LineString
type Geometry struct {
	ID           string `xml:"id,attr"`
	SRSName      string `xml:"srsName,attr"`
	SRSDimension int    `xml:"srsDimension,attr"`
	PosList      string `xml:"posList"`
}

func (g Geometry) GetPosList() string {
	return g.PosList
}

func (g Geometry) GetSRSDimension() int {
	return g.SRSDimension
}

// Represents a point geometry
type Point struct {
	ID           string `xml:"gml:id,attr"`
	SRSName      string `xml:"srsName,attr"`
	SRSDimension int    `xml:"srsDimension,attr"`
	Position     string `xml:"pos"`
}

func (p Point) GetPosList() string {
	return p.Position
}

func (p Point) GetSRSDimension() int {
	return p.SRSDimension
}

// Reference to a start or end node
type NodeRef struct {
	Href string `xml:"href,attr"`
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

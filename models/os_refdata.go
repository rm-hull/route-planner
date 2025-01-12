package models

import (
	"encoding/xml"
)

type Dictionary struct {
	XMLName     xml.Name          `xml:"Dictionary"`
	Description string            `xml:"description"`
	Identifier  Identifier        `xml:"identifier"`
	Entries     []DictionaryEntry `xml:"dictionaryEntry"`
}

type DictionaryEntry struct {
	Definition Definition `xml:"Definition"`
}

type Definition struct {
	ID          string     `xml:"id,attr"`
	Description string     `xml:"description,omitempty"`
	Identifier  Identifier `xml:"identifier"`
}

type Identifier struct {
	CodeSpace string `xml:"codeSpace,attr"`
	Value     string `xml:",chardata"`
}

package sana

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type OuterXML struct {
	XMLName xml.Name `xml:"string"`
	Content string   `xml:",chardata"` // This captures the escaped inner XML string
}

// TbAnavin represents the root of the inner XML structure
type TbAnavin struct {
	XMLName xml.Name     `xml:"tb_Anavin"`
	Items   []AnavinItem `xml:"Item"`
}

type AnavinItem struct {
	ID   string `xml:"ID"`
	Name string `xml:"Name"`
}

type AnavinResponse struct {
	XMLName xml.Name     `xml:"tb_Anavin"`
	Items   []AnavinItem `xml:"Item"`
}

type PayamItem struct {
	ID   string `xml:"ID"`
	Name string `xml:"Name"`
}

type PayamResponse struct {
	Items []PayamItem `xml:"payam"`
}

type SabtChidemanResponse struct {
	Message string
	IsOK    bool
}

type MosafereDataType int

const (
	DataTypeRoomReception     MosafereDataType = 1
	DataTypeNewCompanionEnter MosafereDataType = 11
	DataTypeRoomSettlement    MosafereDataType = 2
	DataTypeCompanionExit     MosafereDataType = 21
	DataTypeRoomReservation   MosafereDataType = 3
	DataTypeReservationCancel MosafereDataType = 31
)

type SabtMosaferinResponse struct {
	ShomareOtagh    string
	ShomarePaziresh string
	RecordMosafer   int
	CodeError       int
	IsOK            bool
}

type ErrorCodeItem struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
}

type ErrorCodeResponse struct {
	Items []ErrorCodeItem `json:"tb_Anavin"`
}

type SabtChidemanInput struct {
	Sakhteman  int
	TedadOtagh int
	Floors     []Floor
}

type Floor struct {
	Number int
	Rooms  []string
}

func BuildRac(input SabtChidemanInput) string {
	var result strings.Builder
	fmt.Fprintf(&result, "sakhteman:%d|tedadotagh:%d", input.Sakhteman, input.TedadOtagh)
	for _, floor := range input.Floors {
		for _, room := range floor.Rooms {
			fmt.Fprintf(&result, "|%d:%s", floor.Number, room)
		}
	}
	result.WriteString("|")
	return result.String()
}

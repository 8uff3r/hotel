package sana

import "fmt"

type AnavinItem struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
}

type AnavinResponse struct {
	Items []AnavinItem `json:"tb_Anavin"`
}

type PayamItem struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
}

type PayamResponse struct {
	Items []PayamItem `json:"payam"`
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
	result := fmt.Sprintf("sakhteman:%d|tedadotagh:%d", input.Sakhteman, input.TedadOtagh)
	for _, floor := range input.Floors {
		for _, room := range floor.Rooms {
			result += fmt.Sprintf("|%d:%s", floor.Number, room)
		}
	}
	result += "|"
	return result
}

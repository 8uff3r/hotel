package sana

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

type RoomRack struct {
	Sakhteman  int
	TedadAtagh int
	OtaghList  []RoomRackFloor
}

type RoomRackFloor struct {
	ShomareTabaghe string
	OtaghList      []string
}

func BuildRac(roomRack RoomRack) string {
	result := "sakhteman:" + string(rune('0'+roomRack.Sakhteman)) + "|tedadotagh:" + string(rune('0'+roomRack.TedadAtagh))
	for _, floor := range roomRack.OtaghList {
		result += "|" + floor.ShomareTabaghe + ":"
		for i, otagh := range floor.OtaghList {
			if i > 0 {
				result += ","
			}
			result += otagh
		}
	}
	return result
}

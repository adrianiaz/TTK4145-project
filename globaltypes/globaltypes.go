package globaltypes

import "encoding/json"

const (
	N_FLOORS  = 4
	N_BUTTONS = 3
)

type ElevatorBehaviour int

const (
	EB_Idle ElevatorBehaviour = iota
	EB_DoorOpen
	EB_Moving
)

type orders2D [N_FLOORS][N_BUTTONS]bool

type AllOrders map[string]orders2D //all the order matrices for all the elevators

type TravelDir int

type ElevatorState struct {
	Floor           int
	Behaviour       ElevatorBehaviour
	TravelDirection TravelDir
	ElevatorID      string
	Requests        orders2D
}

const (
	TravelDown TravelDir = iota - 1
	TravelStop
	TravelUp
)

type Ledger struct {
	BackupMaster []string `json:"backupMaster"` //maby int instead og string if use of IP?
	Alive        []string `json:"alive"`
	Orders       []string `json:"orders"`
}

func Serialize(ledger Ledger) (string, error) {
	//serialize to json
	ledgerJson, err := json.Marshal(ledger)
	if err != nil {
		return "", err
	}
	return string(ledgerJson), nil
}

func Deserialize(ledgerJson string) (*Ledger, error) {
	var ledger Ledger
	err := json.Unmarshal([]byte(ledgerJson), &ledger)
	if err != nil {
		return nil, err
	}
	return &ledger, nil
}

type ButtonType int

type NewOrder struct {
	Floor      int
	BtnType    ButtonType
	ElevatorID int
}

type CompletedOrder struct {
	ElevatorID int
	Floor      int
}

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

type TravelDir int

type ElevatorState struct {
	Floor           int
	Behaviour       ElevatorBehaviour
	TravelDirection TravelDir
	ElevatorID      string
	Requests        [N_FLOORS][N_BUTTONS]bool
}

const (
	TravellingUp TravelDir = iota
	TravellingDown
)

type ButtonType int

const (
	BT_HallUp   ButtonType = 0
	BT_HallDown ButtonType = 1
	BT_Cab      ButtonType = 2
)

//order structs and Ledger struct and member functions

type NewOrder struct {
	ElevatorID int
	Floor      int
	BtnType    ButtonType
}

type CompletedOrder struct {
	ElevatorID int
	Floor      int
	OrderID    int
}
type ActiveOrder struct {
	ElevatorID string
	OrderID    int
	FromFloor  int
	ToFloor    int
}

type Ledger struct {
	//create a map where elevatorID is the key and the value is a slice of ActiveOrders
	ActiveOrders    map[string][]ActiveOrder `json:"activeOrders"`
	ElevatorStates  []ElevatorState          `json:"elevatorStates"`
	BackupMasterlst []string                 `json:"backupMaster"`
	Alive           []bool                   `json:"alive"`
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

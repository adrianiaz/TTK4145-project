package globaldefinitions

import (
	"encoding/json"
)

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

type ElevatorConfig struct {
	ClearRequestVariant ClearRequestVariant
}

type ClearRequestVariant int

type Orders2D [N_FLOORS][N_BUTTONS]bool

type TravelDir int

type ElevatorState struct {
	ElevatorID      string            `json:"elevatorID"`
	Floor           int               `json:"floor"`
	Behaviour       ElevatorBehaviour `json:"behaviour"`
	TravelDirection TravelDir         `json:"travelDirection"`
	Requests        Orders2D          `json:"requests"`
	Config          ElevatorConfig    `json:"config"`
}

const (
	ClearRequests_All ClearRequestVariant = iota
	ClearRequests_InMotorDir
)

const (
	TravelDown TravelDir = iota - 1
	TravelStop
	TravelUp
)

type ButtonType int

const (
	BT_HallUp   ButtonType = 0
	BT_HallDown ButtonType = 1
	BT_Cab      ButtonType = 2
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

//order structs and Ledger struct and member functions

type Order struct {
	NewOrder   bool       `json:"newOrder"` //true if newOrder false if completed
	ElevatorID string     `json:"elevatorID"`
	Floor      int        `json:"floor"`
	BtnType    ButtonType `json:"btnType"`
}

type AllOrders map[string]Orders2D //all the order matrices for all the elevators
type AllElevatorStates map[string]ElevatorState

type Ledger struct {
	//create a map where elevatorID is the key and the value is a slice of ActiveOrders
	ActiveOrders   AllOrders         `json:"activeOrders"`
	ElevatorStates AllElevatorStates `json:"elevatorStates"`
	NodeHierarchy  []string          `json:"backupMaster"` //first element is master, second is backupMaster, rest are slaves
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

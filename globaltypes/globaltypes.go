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

type elevatorState struct {
	Floor           int
	Behaviour       ElevatorBehaviour
	TravelDirection TravelDir
	ElevatorID      string
	Requests        [N_FLOORS][N_BUTTONS]bool
}

const (
	travellingUp TravelDir = iota
	travellingDown
)

type Ledger struct {
	BackupMaster []string `json:"backupMaster"` //maby int instead og string if use of IP?
	Alive        []string `json:"alive"`
	Orders       []string `json:"orders"`
}

func serialize(ledger Ledger) (string, error) {
	//serialize to json
	ledgerJson, err := json.Marshal(ledger)
	if err != nil {
		return "", err
	}
	return string(ledgerJson), nil
}

func deserialize(ledgerJson string) (*Ledger, error) {
	var ledger Ledger
	err := json.Unmarshal([]byte(ledgerJson), &ledger)
	if err != nil {
		return nil, err
	}
	return &ledger, nil
}

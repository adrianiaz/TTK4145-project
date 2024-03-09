package ledger

import (
	"encoding/json"
)

type Ledger struct {
	BackupMaster []string `json:"backupMaster"` //maby int instead og string if use of IP?
	Alive        []string `json:"alive"`
	Orders       []string `json:"orders"`
}

func serialize(ledger Ledger) (string, error){
	//serialize to json
	ledgerJson, err := json.Marshal(ledger)
	if err != nil{
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


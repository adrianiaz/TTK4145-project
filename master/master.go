package master

import (
	gd "github.com/adrianiaz/TTK4145-project/globaldefinitions"
)

//uses ordersToMasterch for both newOrders and completedOrders

// type Ledger struct {
// 	//create a map where elevatorID is the key and the value is a slice of ActiveOrders
// 	ActiveOrders    AllOrders         `json:"activeOrders"`
// 	ElevatorStates  AllElevatorStates `json:"elevatorStates"`
// 	BackupMasterlst []string          `json:"backupMaster"`
// 	Alive           []bool            `json:"alive"`
// }

func Master(ordersToMasterCh <-chan gd.Order) {
	//initialize a ledger with default values
	ledger := gd.Ledger{
		ActiveOrders:    make(gd.AllOrders),
		ElevatorStates:  make(gd.AllElevatorStates),
		BackupMasterlst: make([]string, 0),
		Alive:           make([]bool, 0),
	}
	for {
		select {
		case order := <-ordersToMasterCh:
			switch order.NewOrder {
			case true:
				if order.BtnType == gd.BT_Cab {
					ledger.ActiveOrders[order.ElevatorID] = updateLedger(ledger, order, true)
				} else { //hallcall up or down
					//find the elevator with the least distance to the order
					//append a new active order to the ActiveOrders slice with elevatorID as key
					//send the new active order to the elevator
					return
				}
			case false:
				ledger.ActiveOrders[order.ElevatorID] = updateLedger(ledger, order, false)
			}

		}
	}
}

func updateLedger(ledger gd.Ledger, order gd.Order, orderChange bool) gd.Orders2D {
	orderToChange := ledger.ActiveOrders[order.ElevatorID]
	orderToChange[order.Floor][int(order.BtnType)] = orderChange
	return orderToChange
}

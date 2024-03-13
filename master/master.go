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
					allHallrequests := findAllHallRequests(ledger.ActiveOrders) // Need to include new order in this
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

// should put hall requests for a floor in a
func findAllHallRequests(allorders gd.AllOrders) [gd.N_FLOORS][2]bool {
	var allHallRequests [gd.N_FLOORS][2]bool
	//loop through the gd.AllOrders map
	for _, elevator := range allorders {
		for floor := 0; floor < gd.N_FLOORS; floor++ {
			for btnType := 0; btnType < 2; btnType++ {
				if elevator[floor][btnType] {
					allHallRequests[floor][btnType] = true
				}
			}

		}
	}
	return allHallRequests
}

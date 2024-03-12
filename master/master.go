package master

import (
	"github.com/adrianiaz/TTK4145-project/globaltypes"
)

// placeholder types

//uses ordersToMasterch

func master(ordersToMasterCh <-chan interface{}) {
	//initialize a ledger with default values
	ledger := globaltypes.Ledger{}
	for {
		select {
		case order := <-ordersToMasterCh:
			switch order := order.(type) {
			case globaltypes.NewOrder:
				if order.BtnType == globaltypes.BT_Cab {
					//append a new active order to the ActiveOrders slice with elevatorID as key
					newActiveOrder := globaltypes.ActiveOrder{
						ElevatorID: order.ElevatorID,
						OrderID:    len(ledger.ActiveOrders[order.ElevatorID]), //orderID is the index of the new order in the slice
						ToFloor:    order.Floor,
						FromFloor:  0, //placeholder, need to get the fromFloor from the elevatorState
					}
					ledger.ActiveOrders[order.ElevatorID] = append(ledger.ActiveOrders[order.ElevatorID], newActiveOrder)
				} else { //hallcall up or down
					//find the elevator with the least distance to the order
					//append a new active order to the ActiveOrders slice with elevatorID as key
					//send the new active order to the elevator
					return
				}
			case CompletedOrder:
				// handle completed order
				// ...
			}
		}
	}
}

package master

import (
	"fmt"

	gd "github.com/adrianiaz/TTK4145-project/globaldefinitions"
	"github.com/adrianiaz/TTK4145-project/globaltypes"
)

//uses ordersToMasterch for both newOrders and completedOrders

// type Ledger struct {
// 	//create a map where elevatorID is the key and the value is a slice of ActiveOrders
// 	ActiveOrders    AllOrders         `json:"activeOrders"`
// 	ElevatorStates  AllElevatorStates `json:"elevatorStates"`
// 	BackupMasterlst []string          `json:"backupMaster"`
// 	Alive           []bool            `json:"alive"`
// }

func Master(ordersToMasterCh <-chan interface{}) {
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
			switch order := order.(type) {
			case gd.NewOrder:
				if order.BtnType == gd.BT_Cab {
					//ledger.ActiveOrders[order.ElevatorID] = gd.Orders2D{}
					activeOrders := ledger.ActiveOrders[order.ElevatorID]
					activeOrders[order.Floor][int(order.BtnType)] = true
					ledger.ActiveOrders[order.ElevatorID] = activeOrders
					//send the new active order to the elevator

				} else { //hallcall up or down
					//find the elevator with the least distance to the order
					//append a new active order to the ActiveOrders slice with elevatorID as key
					//send the new active order to the elevator
					return
				}
			case globaltypes.CompletedOrder:
				//remove the completed order from ActiveOrders and rewrites to ledger
				ledger.ActiveOrders[order.ElevatorID] = newActiveOrderlst(ledger.ActiveOrders[order.ElevatorID], order.OrderID)
			}

		}
	}
}

// ElevatorID is local to each elevator, so two elevators can have and order with the same OrderID
func assignOrderID(elevatorActiveOrders []globaltypes.ActiveOrder) int {
	//find the highest orderID in the slice and return orderID+1
	highestOrderID := 0
	for _, order := range elevatorActiveOrders {
		if order.OrderID > highestOrderID {
			highestOrderID = order.OrderID
		}
	}
	return highestOrderID + 1
}

func findOrderIndex(activeOrders []globaltypes.ActiveOrder, orderID int) int {
	for i, order := range activeOrders {
		if order.OrderID == orderID {
			return i
		}
	}
	fmt.Println("Order not found")
	return -1
}

func newActiveOrderlst(elevatorActiveOrders []globaltypes.ActiveOrder, orderID int) []globaltypes.ActiveOrder {
	index := findOrderIndex(elevatorActiveOrders, orderID)
	if index == -1 {
		fmt.Println("No new order, returning original list")
		return elevatorActiveOrders
	}
	elevatorActiveOrders[index] = elevatorActiveOrders[len(elevatorActiveOrders)-1]
	return elevatorActiveOrders[:len(elevatorActiveOrders)-1]
}

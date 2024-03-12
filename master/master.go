package master

import (
	"fmt"

	"github.com/adrianiaz/TTK4145-project/globaltypes"
)

//uses ordersToMasterch for both newOrders and completedOrders

func Master(ordersToMasterCh <-chan interface{}) {
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
						OrderID:    assignOrderID(ledger.ActiveOrders[order.ElevatorID]),
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

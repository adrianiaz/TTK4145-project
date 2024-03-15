package master

import (
	"fmt"
	"sort"
	"strconv"

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

func Master(
	ordersToMasterCh <-chan gd.Order,
	isMaster <-chan bool,
	alive_fromWatchDog <-chan []string,
) {
	//initialize a ledger with default values
	ledger := gd.Ledger{
		ActiveOrders:    make(gd.AllOrders),
		ElevatorStates:  make(gd.AllElevatorStates),
		BackupMasterlst: make([]string, 0),
		Alive:           make([]bool, 0),
	}

slaveLoop:
	for {
		select {
		case <-isMaster:
			//break out of for-loop and start mastermode
			break slaveLoop
		}
	}

	//masterLoop
	for {
		select {
		case order := <-ordersToMasterCh:
			switch order.NewOrder {
			case true:
				if order.BtnType == gd.BT_Cab {
					ledger.ActiveOrders[order.ElevatorID] = updateSingleOrder(ledger, order, true)
				} else { //hall request

					OptimalHallReqAssignment := extractOptimalHallRequests(ledger, order)
					for elevatorID := range ledger.ActiveOrders {
						elevatorOrders := ledger.ActiveOrders[elevatorID]
						for floor := 0; floor < gd.N_FLOORS; floor++ {
							for btnType := 0; btnType < 2; btnType++ {
								elevatorOrders[floor][btnType] = OptimalHallReqAssignment[elevatorID][floor][btnType]
							}
						}
					}
				}
			case false:
				ledger.ActiveOrders[order.ElevatorID] = updateSingleOrder(ledger, order, false)

			}
		case alivePeers := <-alive_fromWatchDog:
			ledger.BackupMasterlst = sortAndRemoveOwnID(alivePeers, "master")
		}
	}
}

func updateSingleOrder(ledger gd.Ledger, order gd.Order, orderChange bool) gd.Orders2D {
	orderToChange := ledger.ActiveOrders[order.ElevatorID]
	orderToChange[order.Floor][int(order.BtnType)] = orderChange
	return orderToChange
}

func sortAndRemoveOwnID(alivePeers []string, ownId string) []string {
	var intPeers []int
	for _, peer := range alivePeers {
		if peer != ownId {
			intPeer, err := strconv.Atoi(peer)
			if err != nil {
				fmt.Printf("Error converting string to int: %v\n", err)
				return alivePeers
			}
			intPeers = append(intPeers, intPeer)
		}
	}
	sort.Ints(intPeers)
	//convert back to string in same order as before
	var sortedPeers []string
	for _, peer := range intPeers {
		sortedPeers = append(sortedPeers, strconv.Itoa(peer))
	}

	return sortedPeers
}

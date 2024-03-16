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
	isMaster <-chan gd.Ledger,
	alive_fromWatchDog <-chan []string,

	ledger_toOrderHandler chan<- gd.Ledger,

	id string,
) {
	//initialize a ledger with default values
	ledger := gd.Ledger{
		ActiveOrders:   make(gd.AllOrders),
		ElevatorStates: make(gd.AllElevatorStates),
		NodeHierarchy:  make([]string, 0),
	}

slaveLoop:
	for {
		select {
		case lastObservedLedger := <-isMaster:
			ledger = lastObservedLedger
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
					ledger_toOrderHandler <- ledger
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
					ledger_toOrderHandler <- ledger
				}
			case false:
				ledger.ActiveOrders[order.ElevatorID] = updateSingleOrder(ledger, order, false)
				ledger_toOrderHandler <- ledger
			}
		case alivePeers := <-alive_fromWatchDog:
			ledger.NodeHierarchy = sortHierarchy(alivePeers, id)
		}
	}
}

func updateSingleOrder(ledger gd.Ledger, order gd.Order, orderChange bool) gd.Orders2D {
	orderToChange := ledger.ActiveOrders[order.ElevatorID]
	orderToChange[order.Floor][int(order.BtnType)] = orderChange
	return orderToChange
}

func sortHierarchy(alivePeers []string, ownId string) []string {
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
	sortedPeers = append(sortedPeers, ownId)
	for _, peer := range intPeers {
		sortedPeers = append(sortedPeers, strconv.Itoa(peer))
	}

	return sortedPeers
}

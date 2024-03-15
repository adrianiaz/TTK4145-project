package network

import (
	gd "github.com/adrianiaz/TTK4145-project/globaldefinitions"
	"github.com/adrianiaz/TTK4145-project/network/peers"
)

func StartNetworkModule(
	ledgerRx chan gd.Ledger,
	ledgerTx chan gd.Ledger,
	singleOrderRx chan gd.Order,
	singleOrderTx chan gd.Order,
	elevatorStateRx chan gd.ElevatorState,
	elevatorStateTx chan gd.ElevatorState,

	peerUpdateCh chan peers.PeerUpdate,
	peerTxEnable chan bool,

	ledger_toOrderHandler chan<- gd.Ledger,
	ledger_fromMaster <-chan gd.Ledger,
	order_fromOrderHandler <-chan gd.Order,
	order_toMaster chan<- gd.Order,
	elevatorState_toMaster chan<- gd.ElevatorState,
	elevatorState_fromElevatorController <-chan gd.ElevatorState,
) {

	id := "placeholderID"

	for {
		select {
		case ledger := <-ledgerRx:
			if id == "master" {
				break
			}
			ledger_toOrderHandler <- ledger

		case ledger := <-ledger_fromMaster: //only happens if current node is master
			ledgerTx <- ledger
			ledger_toOrderHandler <- ledger //send to order module

		case order := <-order_fromOrderHandler:
			if id == "master" {
				order_toMaster <- order
			} else {
				singleOrderTx <- order
			}

		case order := <-singleOrderRx:
			if id == "master" {
				order_toMaster <- order
			}

		case state := <-elevatorStateRx:
			if id == "master" {
				elevatorState_toMaster <- state
			}

		case state := <-elevatorState_fromElevatorController:
			if id == "master" {
				elevatorState_toMaster <- state
			} else {
				elevatorStateTx <- state
			}

		}
	}
}

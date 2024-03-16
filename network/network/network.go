package network

import (
	gd "github.com/adrianiaz/TTK4145-project/globaldefinitions"
)

func NetworkMessageForwarder(
	ledgerRx chan gd.Ledger,
	ledgerTx chan gd.Ledger,
	singleOrderRx chan gd.Order,
	singleOrderTx chan gd.Order,
	elevatorStateRx chan gd.ElevatorState,
	elevatorStateTx chan gd.ElevatorState,
	ledger_toWatchDog chan<- gd.Ledger,
	ledger_toOrderHandler chan<- gd.Ledger,
	ledger_fromMaster <-chan gd.Ledger,
	order_fromOrderHandler <-chan gd.Order,
	order_toMaster chan<- gd.Order,
	elevatorState_toMaster chan<- gd.ElevatorState,
	elevatorState_fromElevatorController <-chan gd.ElevatorState,

	id string,
) {
	var localLedger gd.Ledger

	for {
		select {
		case ledger := <-ledgerRx:
			if id == localLedger.NodeHierarchy[0] {
				break
			}
			ledger_toOrderHandler <- ledger
			ledger_toWatchDog <- ledger
			localLedger = ledger

		case ledger := <-ledger_fromMaster: //only happens if current node is master
			ledgerTx <- ledger
			ledger_toOrderHandler <- ledger
			ledger_toWatchDog <- ledger
			localLedger = ledger

		case order := <-order_fromOrderHandler:
			if id == localLedger.NodeHierarchy[0] {
				order_toMaster <- order
			} else {
				singleOrderTx <- order
			}

		case order := <-singleOrderRx:
			if id == localLedger.NodeHierarchy[0] {
				order_toMaster <- order
			}

		case state := <-elevatorStateRx:
			if id == localLedger.NodeHierarchy[0] {
				elevatorState_toMaster <- state
			}

		case state := <-elevatorState_fromElevatorController:
			if id == localLedger.NodeHierarchy[0] {
				elevatorState_toMaster <- state
			} else {
				elevatorStateTx <- state
			}

		}
	}
}

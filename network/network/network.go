package network

import (
	"fmt"

	gd "github.com/adrianiaz/TTK4145-project/globaldefinitions"
	"github.com/adrianiaz/TTK4145-project/network/bcast"
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

	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	go bcast.Transmitter(16568, elevatorStateTx)
	go bcast.Receiver(16568, elevatorStateRx)

	go bcast.Transmitter(16569, ledgerTx)
	go bcast.Receiver(16569, ledgerRx)

	go bcast.Transmitter(16570, singleOrderTx)
	go bcast.Receiver(16570, singleOrderRx)

	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		//Send new peerlist to master, so that it can write to ledger who is alive
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

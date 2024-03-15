package main

import (
	"flag"

	gd "github.com/adrianiaz/TTK4145-project/globaldefinitions"
	"github.com/adrianiaz/TTK4145-project/network/bcast"
	"github.com/adrianiaz/TTK4145-project/network/peers"
)

func main() {

	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	//Network channels & transmitters/receivers
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	elevatorStateTx := make(chan gd.ElevatorState)
	elevatorStateRx := make(chan gd.ElevatorState)
	ledgerTx := make(chan gd.Ledger)
	ledgerRx := make(chan gd.Ledger)
	singleOrderTx := make(chan gd.Order)
	singleOrderRx := make(chan gd.Order)

	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	go bcast.Transmitter(16568, elevatorStateTx)
	go bcast.Receiver(16568, elevatorStateRx)

	go bcast.Transmitter(16569, ledgerTx)
	go bcast.Receiver(16569, ledgerRx)

	go bcast.Transmitter(16570, singleOrderTx)
	go bcast.Receiver(16570, singleOrderRx)

	select {}
}

package main

import (
	"flag"

	gd "github.com/adrianiaz/TTK4145-project/globaldefinitions"
	"github.com/adrianiaz/TTK4145-project/master"
	"github.com/adrianiaz/TTK4145-project/network/bcast"
	"github.com/adrianiaz/TTK4145-project/network/network"
	"github.com/adrianiaz/TTK4145-project/network/peers"
)

func main() {
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	// transmitters/receivers channels
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

	// master channels
	order_toMaster := make(chan gd.Order)
	elevatorState_toMaster := make(chan gd.ElevatorState)
	isMaster := make(chan gd.Ledger)
	alive_fromWatchDog := make(chan []string)
	ledger_toNetwork := make(chan gd.Ledger)

	// network channels
	ledger_toWatchDog := make(chan gd.Ledger)
	ledger_toOrderHandler := make(chan gd.Ledger)
	order_fromOrderHandler := make(chan gd.Order)
	elevatorState_fromElevatorController := make(chan gd.ElevatorState)

	go master.Master(
		order_toMaster,
		isMaster,
		alive_fromWatchDog,
		elevatorState_toMaster,
		ledger_toNetwork,
		id,
	)

	go network.NetworkMessageForwarder(
		ledgerRx,
		ledgerTx,
		singleOrderRx,
		singleOrderTx,
		elevatorStateRx,
		elevatorStateTx,
		ledger_toWatchDog,
		ledger_toOrderHandler,
		ledger_toNetwork,
		order_fromOrderHandler,
		order_toMaster,
		elevatorState_toMaster,
		elevatorState_fromElevatorController,
		id,
	)

	go network.WatchDog(peerUpdateCh,
		isMaster,
		alive_fromWatchDog,
		ledger_toWatchDog,
		id)

	select {}
}

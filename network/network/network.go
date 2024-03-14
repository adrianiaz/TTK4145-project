package network

import (
	"fmt"

	gd "github.com/adrianiaz/TTK4145-project/globaldefinitions"
	"github.com/adrianiaz/TTK4145-project/network/bcast"
	"github.com/adrianiaz/TTK4145-project/network/peers"
)

// placeholder struct
type Ledger struct {
	orderRequests int
}

func startNetworkModule(
	ledgerRx chan gd.Ledger,
	ledgerTx chan gd.Ledger,
	singleOrderRx chan gd.Order,
	singleOrderTx chan gd.Order,

	peerUpdateCh chan peers.PeerUpdate,
	peerTxEnable chan bool,

	ledger_toOrderHandler chan<- gd.Ledger,
	ledgerMasterCh <-chan gd.Ledger,
	order_fromOrderHandler <-chan gd.Order,
	order_toMaster chan<- gd.Order,
) {

	id := "placeholderID"

	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

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

		case ledger := <-ledgerMasterCh: //only happens if current node is master
			ledgerTx <- ledger
			ledger_toOrderHandler <- ledger //send to order module

		case order := <-order_fromOrderHandler:
			if id == "master" {
				order_toMaster <- order
			} else {
				singleOrderTx <- order
			}
		case order := <-singleOrderRx:
			order_toMaster <- order
		}
	}
}

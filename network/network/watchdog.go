package network

import (
	gd "github.com/adrianiaz/TTK4145-project/globaldefinitions"
	"github.com/adrianiaz/TTK4145-project/network/peers"
)

func WatchDog(
	peerUpdateCh chan peers.PeerUpdate,
	peerTxEnable chan bool,
	isMaster chan<- gd.Ledger,
	alive_toMaster chan<- []string,

	ledger_fromNetwork <-chan gd.Ledger,
	id string,
) {

	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	var localLedger gd.Ledger

	for {
		select {
		case p := <-peerUpdateCh:
			if localLedger.NodeHierarchy[1] == id { //if this node is backupMaster
				for _, peer := range p.Lost {
					if peer == localLedger.NodeHierarchy[0] {
						isMaster <- localLedger
					}
				}
			}
			alive_toMaster <- p.Peers
		case ledger := <-ledger_fromNetwork:
			localLedger = ledger
		}
	}
}

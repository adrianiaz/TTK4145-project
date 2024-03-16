package network

import (
	gd "github.com/adrianiaz/TTK4145-project/globaldefinitions"
	"github.com/adrianiaz/TTK4145-project/network/peers"
)

func WatchDog(
	peerUpdateCh chan peers.PeerUpdate,
	isMaster chan<- gd.Ledger,
	alive_toMaster chan<- []string,

	ledger_fromNetwork <-chan gd.Ledger,
	id string,
) {
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

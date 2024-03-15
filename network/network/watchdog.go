package network

import (
	gd "github.com/adrianiaz/TTK4145-project/globaldefinitions"
	"github.com/adrianiaz/TTK4145-project/network/peers"
)

func WatchDog(
	peerUpdateCh chan peers.PeerUpdate,
	peerTxEnable chan bool,
	isMaster chan<- bool,
	alive_toMaster chan<- []string,

	ledger_toWatchdog <-chan gd.Ledger,
) {
	id := "placeholderID"

	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	var localLedger gd.Ledger

	for {
		select {
		case p := <-peerUpdateCh:
			if localLedger.BackupMasterlst[0] == id {
				for _, peer := range p.Lost {
					if peer == "master" {
						isMaster <- true
					}
				}
			}
			alive_toMaster <- p.Peers
		case ledger := <-ledger_toWatchdog:
			localLedger = ledger
		}
	}
}

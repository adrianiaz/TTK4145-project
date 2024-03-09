package network

import "fmt"

// placeholder struct
type Ledger struct {
	orderRequests int
}

func startNetworkModule(downstreamLedgerCh chan Ledger, newOrderCh chan []int, servicedOrderCh chan []int) {

	//placeholder information
	ledgerplaceholder := Ledger{orderRequests: 0}

	for {
		select {
		case ledger := <-downstreamLedgerCh:
			// Placeholder: Handle incoming ledger data
			fmt.Println("Received downstream ledger:", ledger)
			fmt.Println(ledgerplaceholder.orderRequests)
		case newOrder := <-newOrderCh:
			// Placeholder: Handle incoming new order
			fmt.Println("Received new order:", newOrder)
		case servicedOrder := <-servicedOrderCh:
			// Placeholder: Handle incoming serviced order
			fmt.Println("Received serviced order:", servicedOrder)
		}

	}
}

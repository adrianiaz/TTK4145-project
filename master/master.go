package master

//"github.com/adrianiaz/TTK4145-project/ledger"

// placeholder types
type ButtonType int

const (
	BT_HallUp   ButtonType = 0
	BT_HallDown            = 1
	BT_Cab                 = 2
)

type NewOrder struct {
	Floor      int
	BtnType    ButtonType
	ElevatorID int
}

type CompletedOrder struct {
	ElevatorID int
	Floor      int
}

//uses ordersToMasterch

func master(ordersToMasterCh <-chan interface{}) {
	//instantiate new ledger with default values
	//ledger := ledger.Ledger{BackupMaster: []string{" "}, Alive: true, Orders: []string{" "}}

	for {
		select {
		case order := <-ordersToMasterCh:
			switch order := order.(type) {
			case NewOrder:
				order.ElevatorID = 0 // placegholder for stopping warnings
				// handle new order
				// ...
			case CompletedOrder:
				// handle completed order
				// ...
			}
		}
	}
}

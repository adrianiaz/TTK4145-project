package orderHandler

//import(elevio)

type ButtonType int

const (
	BT_HallUp   ButtonType = 0
	BT_HallDown            = 1
	BT_Cab                 = 2
)

type ButtonEvent struct{
	Floor  int
	Button ButtonType
}

type NewOrder struct {
	Floor   int
	BtnType ButtonType
	ElevatorID int
}

func orderHandler(ButtonPressCh chan ButtonEvent, NewOrderCh chan NewOrder) {

	for {
		select {

		case button := <-ButtonPressCh:
			//Do something with the button event
			//buttons := []ButtonEvent{button} //A slice of ButtonEvents
			newOrder := NewOrder{
				Floor: button.Floor,
				BtnType: button.Button,
				ElevatorID: 0, //placeholder
			}
			NewOrderCh <- newOrder
		}

	}
}

package orderHandler

//import(elevio)

// type ButtonEvent struct{
// 	Floor  int
// 	Button ButtonType
// }

type NewOrder struct {
	Floor   int
	BtnType ButtonType
}

func orderHandler(ButtonPressChan chan ButtonEvent, CompletedOrderChan chan OrderEvent) {

	for {
		select {

		case button := <-ButtonPressChan:
			//Do something with the button event
			//buttons := []ButtonEvent{button} //A slice of ButtonEvents

		//orders complete
		case completed := <-CompletedOrderChan:
			//Do something with the completed order. Send to network.

		}

	}
}

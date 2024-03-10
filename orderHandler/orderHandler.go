package orderHandler

//import(elevio)

func orderHandler(ButtonPressCh chan ButtonEvent, NewOrderCh chan NewOrder) {

	for {
		select {

		case button := <-ButtonPressCh:
			//Do something with the button event
			//buttons := []ButtonEvent{button} //A slice of ButtonEvents
			newOrder := NewOrder{
				Floor:      button.Floor,
				BtnType:    button.Button,
				ElevatorID: 0, //placeholder
			}
			NewOrderCh <- newOrder
		}

	}
}

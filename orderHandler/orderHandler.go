package orderHandler

import (
	gd "github.com/adrianiaz/TTK4145-project/globaldefinitions"
)

func orderHandler(ButtonPressCh chan gd.ButtonEvent, NewOrderCh chan gd.NewOrder) {

	for {
		select {

		case button := <-ButtonPressCh:
			//Do something with the button event
			//buttons := []ButtonEvent{button} //A slice of ButtonEvents
			newOrder := gd.NewOrder{
				Floor:      button.Floor,
				BtnType:    button.Button,
				ElevatorID: "0", //placeholder
			}
			NewOrderCh <- newOrder
		}

	}
}

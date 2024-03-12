package orderHandler

import (
	"github.com/adrianiaz/TTK4145-project/globaltypes"
)

func orderHandler(ButtonPressCh chan ButtonEvent, NewOrderCh chan globaltypes.NewOrder) {

	for {
		select {

		case button := <-ButtonPressCh:
			//Do something with the button event
			//buttons := []ButtonEvent{button} //A slice of ButtonEvents
			newOrder := globaltypes.NewOrder{
				Floor:      button.Floor,
				BtnType:    button.Button,
				ElevatorID: 0, //placeholder
			}
			NewOrderCh <- newOrder
		}

	}
}

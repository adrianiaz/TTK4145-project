package orderHandler

import (
	gd "github.com/adrianiaz/TTK4145-project/globaldefinitions"
)

func orderHandlerModule(ButtonPressCh chan gd.ButtonEvent, 
						NewOrderCh chan gd.NewOrder, 
						ID string, 
						FromMasterCh chan gd.Ledger, 
						ToElevatorControllerCh chan gd.Orders2D,
						) {

	for {
		select {

		case button := <-ButtonPressCh: //(ButtonEvent struct)
			newOrder := gd.NewOrder{
				Floor:      button.Floor,
				BtnType:    button.Button,
				ElevatorID: "0", //placeholder
			}
			NewOrderCh <- newOrder //to master

		case ledgerFromMaster := <-FromMasterCh: //(distributed ledger(struct))
			newLocalOrder2D := gd.Orders2D{
				ledgerFromMaster.AllOrders[ID] //newOrder2D is a matrix[bool]
			}
			ToElevatorControllerCh <- newLocalOrder2D

			//set lights
			lightMatrix := newLocalOrder2D
			for _, matrix := range ledgerFromMaster.ActiveOrders {
				for floor:=0; floor < gd.N_FLOORS; floor++ {
					for btn:=0; btn < gd.N_BUTTONS-1; btn++ {
						lightMatrix[floor][btn] = matrix[floor][btn] || newLocalOrder2D[floor][btn]
					}
			lightCh <- lightMatrix
		}
	}
}

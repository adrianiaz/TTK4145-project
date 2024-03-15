package orderHandler

import (
	gd "github.com/adrianiaz/TTK4145-project/globaldefinitions"
)


func orderHandlerModule(
	ID string,
	ButtonPressCh          <-chan gd.ButtonEvent, 
	CompletedOrderCh       <-chan gd.ButtonEvent, //buttonEvent? or order or matrix?
	LedgerFromMasterCh     <-chan gd.Ledger,      
	OrderToMasterCh        chan<- gd.Order,       
	OrdersToElevatorCtrlCh chan<- gd.Orders2D,    
	LightCh                chan   gd.Orders2D,      
) {

	for {
		select {

		case button := <-ButtonPressCh:
			newOrder := gd.Order{
				NewOrder:   true,
				Floor:      button.Floor,
				BtnType:    button.Button,
				ElevatorID: ID,
			}
			OrderToMasterCh <- newOrder

		case ledgerFromMaster := <-LedgerFromMasterCh:
			newLocalOrder2D := ledgerFromMaster.ActiveOrders[ID]
			OrdersToElevatorCtrlCh <- newLocalOrder2D

			//set lights
			lightMatrix := newLocalOrder2D
			for _, matrix := range ledgerFromMaster.ActiveOrders {
				for floor := 0; floor < gd.N_FLOORS; floor++ {
					for btn := 0; btn < gd.N_BUTTONS-1; btn++ {
						lightMatrix[floor][btn] = matrix[floor][btn] || newLocalOrder2D[floor][btn]
					}
				}
			}
			LightCh <- lightMatrix
			


		case completedOrder := <-CompletedOrderCh: //(ButtonEvent struct) from elevatorcontroller?
			newCompletedOrder := gd.Order{
				NewOrder:   false,
				ElevatorID: ID,
				Floor:      completedOrder.Floor,
				BtnType:    completedOrder.Button,
			}
			OrderToMasterCh <- newCompletedOrder
		}
	}
}

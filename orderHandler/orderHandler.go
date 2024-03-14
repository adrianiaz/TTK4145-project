package orderHandler

import (
	gd "github.com/adrianiaz/TTK4145-project/globaldefinitions"
)

type OrderHandlerChannels struct {
	ButtonPressCh          <-chan gd.ButtonEvent // Receive button presses
	NewOrderCh             chan<- gd.Order       // Send new orders
	LedgerFromMasterCh     <-chan gd.Ledger      // Receive ledger from master
	OrderToMasterCh        chan<- gd.Order       // Send orders to master
	OrdersToElevatorCtrlCh chan<- gd.Orders2D    // Send orders to elevator controller
	CompletedOrderCh       <-chan gd.ButtonEvent // Receive completed orders
	LightCh                chan gd.Orders2D      // Bi-directional channel for updating elevator lights
}

func orderHandlerModule(
	ID string,
	CH OrderHandlerChannels,
) {

	for {
		select {

		case button := <-CH.ButtonPressCh: //(ButtonEvent struct)
			newOrder := gd.Order{
				NewOrder:   true,
				Floor:      button.Floor,
				BtnType:    button.Button,
				ElevatorID: "0", //placeholder
			}
			CH.OrderToMasterCh <- newOrder

		case ledgerFromMaster := <-CH.LedgerFromMasterCh:
			newLocalOrder2D := ledgerFromMaster.ActiveOrders[ID]
			CH.OrdersToElevatorCtrlCh <- newLocalOrder2D

			//set lights
			lightMatrix := newLocalOrder2D
			for _, matrix := range ledgerFromMaster.ActiveOrders {
				for floor := 0; floor < gd.N_FLOORS; floor++ {
					for btn := 0; btn < gd.N_BUTTONS-1; btn++ {
						lightMatrix[floor][btn] = matrix[floor][btn] || newLocalOrder2D[floor][btn]
					}
				}
			}
			CH.LightCh <- lightMatrix

		case completedOrder := <-CH.CompletedOrderCh: //(ButtonEvent struct) from elevatorcontroller
			newCompletedOrder := gd.Order{
				NewOrder:   false,
				ElevatorID: ID,
				Floor:      completedOrder.Floor,
				BtnType:    completedOrder.Button,
			}
			CH.OrderToMasterCh <- newCompletedOrder //to master

			//update lights matrix
			lightMatrixUpdate := <-CH.LightCh
			lightMatrixUpdate[completedOrder.Floor][completedOrder.Button] = false
			CH.LightCh <- lightMatrixUpdate
		}
	}
}

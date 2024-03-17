package orderHandler

import (
	"github.com/adrianiaz/TTK4145-project/elevio"
	gd "github.com/adrianiaz/TTK4145-project/globaldefinitions"
)

func OrderHandler(
	ID string,
	ButtonPressCh chan gd.ButtonEvent,
	CompletedOrderCh <-chan gd.ButtonEvent,
	Ledger_fromMasterCh <-chan gd.Ledger,
	Order_toMasterCh chan<- gd.Order,
	Orders_toElevatorCtrlCh chan<- gd.Orders2D,
	LightCh chan gd.Orders2D,
) {

	go elevio.PollButtons(ButtonPressCh)

	for {
		select {

		case button := <-ButtonPressCh:
			newOrder := gd.Order{
				NewOrder:   true,
				Floor:      button.Floor,
				BtnType:    button.Button,
				ElevatorID: ID,
			}
			Order_toMasterCh <- newOrder

		case ledgerFromMaster := <-Ledger_fromMasterCh:
			newLocalOrder2D := ledgerFromMaster.ActiveOrders[ID]
			Orders_toElevatorCtrlCh <- newLocalOrder2D

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

		case completedOrder := <-CompletedOrderCh:
			newCompletedOrder := gd.Order{
				NewOrder:   false,
				ElevatorID: ID,
				Floor:      completedOrder.Floor,
				BtnType:    completedOrder.Button,
			}
			Order_toMasterCh <- newCompletedOrder
		}
	}
}

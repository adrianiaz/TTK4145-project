package elevatorcontroller

import (
	"github.com/adrianiaz/TTK4145-project/elevio"
	gd "github.com/adrianiaz/TTK4145-project/globaldefinitions"
)

func (elev Elevator) fsm_onInitBetweenFloors() {
	elevio.SetMotorDirection(elevio.MD_Down)
	elev.State.Behaviour = gd.EB_Moving
	elev.State.TravelDirection = gd.TravelDown
}

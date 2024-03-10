package elevatorcontroller

import (
	"github.com/adrianiaz/TTK4145-project/elevio"
	"github.com/adrianiaz/TTK4145-project/globaltypes"
)

type Elevator struct {
	Floor     int
	Behaviour globaltypes.ElevatorBehaviour
	Direction elevio.MotorDirection
	Requests  [][]bool
}

/* type ClearRequestVariant struct {
} */

func (elev Elevator) RespondToOrder() bool {
	switch elev.Behaviour {
	}
	return false
}

/* func NewElevator(numFloors int) *Elevator {
	requests := make([][]bool, numFloors)
	for i := range requests {
		requests[i] = make([]bool, 3)
	}
	return &Elevator{
		Floor:     0,
		Behaviour: EB_Idle,
		Direction: elevio.MD_Stop,
		Requests:  requests,
	}
} */

package elevatorcontroller

import (
	"time"

	. "github.com/adrianiaz/TTK4145-project/elevio"
	. "github.com/adrianiaz/TTK4145-project/globaltypes"
)

type Elevator struct {
	State  ElevatorState
	lights orders2D
	orders orders2D
}

type ElevatorChannels struct {
	OrderCh        chan bool
	CurrentFloorCh chan int
	ObstructionCh  chan bool
	StopCh         chan bool
}

func InitiateElevator(ElevatorID string) Elevator {
	elev := Elevator{
		State: ElevatorState{
			Floor:           0,
			Behaviour:       EB_Idle,
			TravelDirection: TravelStop,
			ElevatorID:      ElevatorID,
		},
		orders{},
		lights{},
	}
	return elev
}

func StartElevatorController(ElevatorID string, orderch OrderCh) {

	Println("ElevatorController started")
	elev := InitiateElevator(ElevatorID)

	doorOpenCh := make(chan bool, 100)
	doorOpenLamp := time.Newtimer(time.Second * 3)
	doorOpenLamp.Stop()

}

/* func (elev Elevator) UpdateState(floor int, behaviour ElevatorBehaviour, travelDir TravelDir) {
	elev.State.Floor = floor
	elev.State.Behaviour = behaviour
	elev.State.TravelDirection = travelDir
} */

func (elev Elevator) orderRegisteredAbove() bool {
	for floorAbove := elev.State.Floor + 1; floorAbove < N_FLOORS; floorAbove++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if elev.orders[floorAbove][btn] {
				return true
			}
		}
	}
	return false
}

func (elev Elevator) orderRegisteredBelow() bool {
	for floorBelow := elev.State.Floor - 1; floorBelow >= 0; floorBelow-- {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if elev.orders[floorBelow][btn] {
				return true
			}
		}
	}
	return false
}

func (elev Elevator) chooseDirection() TravelDir {
	orderAbove := elev.orderRegisteredAbove()
	orderBelow := elev.orderRegisteredBelow()

	switch elev.State.TravelDirection {
	case TravelUp:
		if orderAbove {
			return TravelUp
		} else if orderBelow {
			return TravelDown
		} else {
			return TravelDown
		}

	case TravelDown:
		if orderBelow {
			return TravelDown
		} else if orderAbove {
			return TravelUp
		} else {
			return TravelUp
		}

	default:
		return TravelStop
	}
}

/* func (elev Elevator) oppositeDirection() TravelDir {
	switch elev.State.TravelDirection {
	case TravelUp:
		return TravelDown
	case TravelDown:
		return TravelUp
	default:
		Println("Error in oppositeDirection")
		return TravelUp
} */

func (elev Elevator) orderRegisteredAtCurrentFloor() bool {
	for btn := 0; btn < N_BUTTONS; btn++ {
		if elev.orders[elev.State.Floor][btn] {
			return true
		}
	}
	return false
}

func (elev Elevator) elevatorShouldStop() bool {
	orderAbove := elev.orderRegisteredAbove()
	orderBelow := elev.orderRegisteredBelow()
	currentFloor := elev.State.Floor

	switch elev.State.TravelDirection {
	case TravelUp:
		return elev.orders[currentFloor][BT_Cab] || elev.orders[currentFloor][BT_HallUp] || (!orderAbove)

	case TravelDown:
		return elev.orders[currentFloor][BT_Cab] || elev.orders[currentFloor][BT_HallDown] || (!orderBelow)

	case TravelStop:
		return false
	default:
		return false
	}
}

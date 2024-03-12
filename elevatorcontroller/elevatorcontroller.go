package elevatorcontroller

import (
	"fmt"
	"time"

	. "github.com/adrianiaz/TTK4145-project/elevio"
	gd "github.com/adrianiaz/TTK4145-project/globaldefinitions"
)

type Elevator struct {
	State  gd.ElevatorState
	lights gd.Orders2D
	orders gd.Orders2D
}

type ElevatorChannels struct {
	OrderCh        chan bool
	CurrentFloorCh chan int
	ObstructionCh  chan bool
	StopCh         chan bool
}

func InitiateElevator(ElevatorID string, addr string, numFloors int) Elevator {
	Init(addr, numFloors)
	elev := Elevator{
		State: gd.ElevatorState{
			Floor:           0,
			Behaviour:       gd.EB_Idle,
			TravelDirection: gd.TravelStop,
			ElevatorID:      ElevatorID,
		},
		orders: gd.Orders2D{},
		lights: gd.Orders2D{},
	}
	return elev
}

func StartElevatorController(ElevatorID string, addr string, numFloors int, orderch bool) {

	fmt.Println("ElevatorController started")
	elev := InitiateElevator(ElevatorID, addr, numFloors)

	doorOpenCh := make(chan bool, 100)
	doorOpenLamp := time.NewTimer(time.Second * 3)
	doorOpenLamp.Stop()

}

/* func (elev Elevator) UpdateState(floor int, behaviour ElevatorBehaviour, travelDir TravelDir) {
	elev.State.Floor = floor
	elev.State.Behaviour = behaviour
	elev.State.TravelDirection = travelDir
} */

func (elev Elevator) orderRegisteredAbove() bool {
	for floorAbove := elev.State.Floor + 1; floorAbove < gd.N_FLOORS; floorAbove++ {
		for btn := 0; btn < gd.N_BUTTONS; btn++ {
			if elev.orders[floorAbove][btn] {
				return true
			}
		}
	}
	return false
}

func (elev Elevator) orderRegisteredBelow() bool {
	for floorBelow := elev.State.Floor - 1; floorBelow >= 0; floorBelow-- {
		for btn := 0; btn < gd.N_BUTTONS; btn++ {
			if elev.orders[floorBelow][btn] {
				return true
			}
		}
	}
	return false
}

func (elev Elevator) chooseDirection() gd.TravelDir {
	orderAbove := elev.orderRegisteredAbove()
	orderBelow := elev.orderRegisteredBelow()

	switch elev.State.TravelDirection {
	case gd.TravelUp:
		if orderAbove {
			return gd.TravelUp
		} else if orderBelow {
			return gd.TravelDown
		} else {
			return gd.TravelDown //implemented to prevent unexpected behaviour or getting stuck
		}

	case gd.TravelDown:
		if orderBelow {
			return gd.TravelDown
		} else if orderAbove {
			return gd.TravelUp
		} else {
			return gd.TravelUp //implemented to prevent unexpected behaviour or getting stuck
		}

	default:
		return gd.TravelStop
	}
}

// Just used in the function requestsChooseDirection
type DirBehaviourPair struct {
	Dir       gd.TravelDir
	Behaviour gd.ElevatorBehaviour
}

// alternative to chooseDirection
func (elev Elevator) requestsChooseDirection() DirBehaviourPair {
	orderAbove := elev.orderRegisteredAbove()
	orderBelow := elev.orderRegisteredBelow()
	orderHere := elev.orderRegisteredAtCurrentFloor()

	switch elev.State.TravelDirection {
	case gd.TravelUp:
		if orderAbove {
			return DirBehaviourPair{gd.TravelUp, gd.EB_Moving}
		} else if orderHere {
			return DirBehaviourPair{gd.TravelStop, gd.EB_DoorOpen}
		} else if orderBelow {
			return DirBehaviourPair{gd.TravelDown, gd.EB_Moving}
		} else {
			return DirBehaviourPair{gd.TravelStop, gd.EB_Idle} //consider implementing an alternative that handles unexpected behaviour
		} //e.g. changing to {gd.TravelDown, gd.EB_Moving} if no orders are above
	case gd.TravelDown:
		if orderBelow {
			return DirBehaviourPair{gd.TravelDown, gd.EB_Moving}
		} else if orderHere {
			return DirBehaviourPair{gd.TravelStop, gd.EB_DoorOpen}
		} else if orderAbove {
			return DirBehaviourPair{gd.TravelUp, gd.EB_Moving}
		} else {
			return DirBehaviourPair{gd.TravelStop, gd.EB_Idle}
		}
	case gd.TravelStop:
		if orderHere {
			return DirBehaviourPair{gd.TravelStop, gd.EB_DoorOpen}
		} else if orderAbove {
			return DirBehaviourPair{gd.TravelUp, gd.EB_Moving}
		} else if orderBelow {
			return DirBehaviourPair{gd.TravelDown, gd.EB_Moving}
		} else {
			return DirBehaviourPair{gd.TravelStop, gd.EB_Idle}
		}
	default:
		return DirBehaviourPair{gd.TravelStop, gd.EB_Idle}
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
	for btn := 0; btn < gd.N_BUTTONS; btn++ {
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
	case gd.TravelUp:
		return elev.orders[currentFloor][gd.BT_Cab] || elev.orders[currentFloor][gd.BT_HallUp] || (!orderAbove)

	case gd.TravelDown:
		return elev.orders[currentFloor][gd.BT_Cab] || elev.orders[currentFloor][gd.BT_HallDown] || (!orderBelow)

	case gd.TravelStop:
		return false
	default:
		return false
	}
}

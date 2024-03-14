package elevatorcontroller

import (
	"fmt"
	"time"

	"github.com/adrianiaz/TTK4145-project/elevio"
	//. "github.com/adrianiaz/TTK4145-project/elevio"
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
	elevio.Init(addr, numFloors)
	elev := Elevator{
		State: gd.ElevatorState{
			Floor:           0,
			Behaviour:       gd.EB_Idle,
			TravelDirection: gd.TravelStop,
			ElevatorID:      ElevatorID,
			Config: gd.ElevatorConfig{
				ClearRequestVariant: gd.ClearRequests_InMotorDir,
			},
		},
		orders: gd.Orders2D{},
		lights: gd.Orders2D{},
	}
	return elev
}

func StartElevatorController(ElevatorID string, addr string, numFloors int, ch ElevatorChannels) {

	fmt.Println("ElevatorController started")
	elev := InitiateElevator(ElevatorID, addr, numFloors)

	//doorOpenCh := make(chan bool, 100)
	doorOpenTimer := time.NewTimer(3 * time.Second)
	doorOpenTimer.Stop()
	elevio.SetDoorOpenLamp(false)

	for {
		select {
		case ArrivingAtFloor := <-ch.CurrentFloorCh:
			elev.State.Floor = ArrivingAtFloor
			elevio.SetFloorIndicator(elev.State.Floor)

			switch elev.State.Behaviour {
			case gd.EB_Moving:
				if elev.shouldStop() {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elev.State.Behaviour = gd.EB_DoorOpen
					//elev.State.Config.ClearRequestVariant = gd.ClearRequests_InMotorDir //not sure if this is necessary here considering it is already set in the InitiateElevator function
					elevio.SetDoorOpenLamp(true)
					doorOpenTimer.Reset(3 * time.Second)
				} else {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elev.State.Behaviour = gd.EB_Idle
				}

				//fmt.Println("Current floor: ", elev.State.Floor, "\n Behaviour: ", elev.State.Behaviour, "\n TravelDirection: ", elev.State.TravelDirection)
				fmt.Println("Arriving at floor: ", elev.State.Floor, "\n Behaviour: ", elev.chooseDirection().Behaviour, "\n TravelDirection: ", elev.chooseDirection().Dir)
				elevio.SetMotorDirection(elevio.MD_Stop)
				elev.State.Behaviour = gd.EB_Idle
				elev.State.Config.ClearRequestVariant = gd.ClearRequests_InMotorDir //not sure if this is necessary here considering it is already set in the InitiateElevator function
				doorOpenTimer.Reset(3 * time.Second)
			}
		case <-doorOpenTimer.C:
			elevio.SetDoorOpenLamp(false)

		default:
			elev.State.Floor = -1
			elevio.SetMotorDirection(elevio.MD_Down)
			elev.State.Behaviour = gd.EB_Moving
			elev.State.Config.ClearRequestVariant = gd.ClearRequests_InMotorDir //not sure if this is necessary here considering it is already set in the InitiateElevator function
			doorOpenTimer.Reset(3 * time.Second)
		}
	}
}

/* func (elev Elevator) UpdateState(floor int, behaviour ElevatorBehaviour, travelDir TravelDir) {
	elev.State.Floor = floor
	elev.State.Behaviour = behaviour
	elev.State.TravelDirection = travelDir
} */

func (elev Elevator) orderAbove() bool {
	for floorAbove := elev.State.Floor + 1; floorAbove < gd.N_FLOORS; floorAbove++ {
		for btn := 0; btn < gd.N_BUTTONS; btn++ {
			if elev.orders[floorAbove][btn] {
				return true
			}
		}
	}
	return false
}

func (elev Elevator) orderBelow() bool {
	for floorBelow := elev.State.Floor - 1; floorBelow >= 0; floorBelow-- {
		for btn := 0; btn < gd.N_BUTTONS; btn++ {
			if elev.orders[floorBelow][btn] {
				return true
			}
		}
	}
	return false
}

func (elev Elevator) orderAtCurrentFloor() bool {
	for btn := 0; btn < gd.N_BUTTONS; btn++ {
		if elev.orders[elev.State.Floor][btn] {
			return true
		}
	}
	return false
}

/* func (elev Elevator) chooseDirection2() gd.TravelDir {
	orderAbove := elev.orderAbove()
	orderBelow := elev.orderBelow()

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
} */

// Just used in the function requestsChooseDirection
type DirBehaviourPair struct {
	Dir       gd.TravelDir
	Behaviour gd.ElevatorBehaviour
}

// alternative to chooseDirection
func (elev Elevator) chooseDirection() DirBehaviourPair {
	//orderAbove := elev.orderAbove()
	//orderBelow := elev.orderBelow()
	//orderHere := elev.orderAtCurrentFloor()

	switch elev.State.TravelDirection {
	case gd.TravelUp:
		if elev.orderAbove() {
			return DirBehaviourPair{gd.TravelUp, gd.EB_Moving}
		} else if elev.orderAtCurrentFloor() {
			return DirBehaviourPair{gd.TravelDown, gd.EB_DoorOpen}
		} else if elev.orderBelow() {
			return DirBehaviourPair{gd.TravelDown, gd.EB_Moving}
		} else {
			return DirBehaviourPair{gd.TravelStop, gd.EB_Idle} //consider implementing an alternative that handles unexpected behaviour
		} //e.g. changing to {gd.TravelDown, gd.EB_Moving} if no orders are above
	case gd.TravelDown:
		if elev.orderBelow() {
			return DirBehaviourPair{gd.TravelDown, gd.EB_Moving}
		} else if elev.orderAtCurrentFloor() {
			return DirBehaviourPair{gd.TravelUp, gd.EB_DoorOpen}
		} else if elev.orderAbove() {
			return DirBehaviourPair{gd.TravelUp, gd.EB_Moving}
		} else {
			return DirBehaviourPair{gd.TravelStop, gd.EB_Idle}
		}
	case gd.TravelStop:
		if elev.orderAtCurrentFloor() {
			return DirBehaviourPair{gd.TravelStop, gd.EB_DoorOpen}
		} else if elev.orderAbove() {
			return DirBehaviourPair{gd.TravelUp, gd.EB_Moving}
		} else if elev.orderBelow() {
			return DirBehaviourPair{gd.TravelDown, gd.EB_Moving}
		} else {
			return DirBehaviourPair{gd.TravelStop, gd.EB_Idle}
		}
	default:
		return DirBehaviourPair{gd.TravelStop, gd.EB_Idle}
	}
}

func (elev Elevator) shouldStop() bool {
	//orderAbove := elev.orderAbove()
	//orderBelow := elev.orderBelow()
	currentFloor := elev.State.Floor

	switch elev.State.TravelDirection {
	case gd.TravelUp:
		return elev.orders[currentFloor][gd.BT_HallUp] || elev.orders[currentFloor][gd.BT_Cab] || (!elev.orderAbove())

	case gd.TravelDown:
		return elev.orders[currentFloor][gd.BT_HallDown] || elev.orders[currentFloor][gd.BT_Cab] || (!elev.orderBelow())

	case gd.TravelStop:
		return false
	default:
		return false
	}
}

func (elev Elevator) ClearOrdersAtCurrentFloor() {
	switch elev.State.Config.ClearRequestVariant {
	case gd.ClearRequests_All:
		for btn := 0; btn < gd.N_BUTTONS; btn++ {
			elev.orders[elev.State.Floor][btn] = false
		}

	case gd.ClearRequests_InMotorDir:
		elev.orders[elev.State.Floor][gd.BT_Cab] = false

		switch elev.State.TravelDirection {
		case gd.TravelUp:
			elev.orders[elev.State.Floor][gd.BT_HallUp] = false
			if !elev.orderAbove() {
				elev.orders[elev.State.Floor][gd.BT_HallDown] = false
			}
		case gd.TravelDown:
			elev.orders[elev.State.Floor][gd.BT_HallDown] = false
			if !elev.orderBelow() {
				elev.orders[elev.State.Floor][gd.BT_HallUp] = false
			}
		case gd.TravelStop:
		default:
			elev.orders[elev.State.Floor][gd.BT_HallUp] = false
			elev.orders[elev.State.Floor][gd.BT_HallDown] = false
		}
	}
}

// have to see if this will be necessary here
func (elev Elevator) clearLightsAtCurrentFloor() {
	for btn := 0; btn < gd.N_BUTTONS; btn++ {
		elev.lights[elev.State.Floor][btn] = false
	}
}

// setButtonLamp sets the light of the button at the given floor to the given value, might have to make some adjustments to this function
func (elev Elevator) setButtonLights() {
	elevio.SetFloorIndicator(elev.State.Floor)
	for floor := 0; floor < gd.N_FLOORS; floor++ {
		for btn := 0; btn < gd.N_BUTTONS; btn++ {
			elevio.SetButtonLamp(gd.ButtonType(btn), floor, elev.lights[floor][btn])
		}
	}
}

// sets the floor indicator and the cab lights
func (elev Elevator) setCabLights() {
	elevio.SetFloorIndicator(elev.State.Floor)
	for floor := 0; floor < gd.N_FLOORS; floor++ {
		elevio.SetButtonLamp(gd.BT_Cab, floor, elev.lights[floor][gd.BT_Cab])
	}
}

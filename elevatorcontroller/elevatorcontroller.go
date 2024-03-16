package elevatorcontroller

import (
	"fmt"
	"time"

	"github.com/adrianiaz/TTK4145-project/elevio"
	gd "github.com/adrianiaz/TTK4145-project/globaldefinitions"
)

type Elevator struct {
	State  gd.ElevatorState
	lights gd.Orders2D
	orders gd.Orders2D
}

type ElevatorChannels struct {
	OrderCh          chan bool //add arrows to indicate direction
	CurrentFloorCh   chan int
	ObstructionEvent chan bool
	StopCh           chan bool
	LocalLights2D    chan gd.Orders2D
	LocalOrder2D     <-chan gd.Orders2D
	ToMaster         chan<- gd.ElevatorState
}

const (
	timeUntilTimeout    time.Duration = 10 * time.Second
	timeUntilDoorCloses time.Duration = 3 * time.Second
)

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

	doorOpening := make(chan bool, 100)
	doorOpenDuration := time.NewTimer(timeUntilDoorCloses)
	doorOpenDuration.Stop()
	timeoutTimer := time.NewTimer(timeUntilTimeout) //increase this a bit in case of for example several obstruction events in a row
	elevio.SetDoorOpenLamp(false)
	obstruction := false

	sendElevatorState := time.NewTimer(1 * time.Second)

	for {
		select {
		case ArrivingAtFloor := <-ch.CurrentFloorCh:
			elev.State.Floor = ArrivingAtFloor
			fmt.Println("Arriving at floor: ", elev.State.Floor)
			elevio.SetFloorIndicator(elev.State.Floor)

			switch elev.State.Behaviour {
			case gd.EB_Moving:
				if elev.shouldStop() {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elev.clearOrdersAtCurrentFloor()
					doorOpening <- true
					elev.setButtonLights() //with channel?, doublecheck if this is correct, in the line above the setDoorOpenLamp func is used
					elev.State.Behaviour = gd.EB_DoorOpen
				}
				if (elev.orders == gd.Orders2D{}) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elev.State.Behaviour = gd.EB_Idle
					timeoutTimer.Stop()
				}
				timeoutTimer.Reset(timeUntilTimeout)

			case gd.EB_Idle, gd.EB_DoorOpen:
				elevio.SetMotorDirection(elevio.MD_Stop)
				timeoutTimer.Stop()
			}
			ch.ToMaster <- elev.State

		case <-doorOpening:
			fmt.Println("The door is open")
			elevio.SetDoorOpenLamp(true)
			doorOpenDuration.Reset(timeUntilDoorCloses)
			elevio.SetDoorOpenLamp(false)
			timeoutTimer.Stop()

		case doorObstructed := <-ch.ObstructionEvent:
			obstruction = doorObstructed
			if obstruction {
				fmt.Println("Obstruction detected")
				elevio.SetDoorOpenLamp(true)
				elev.State.Behaviour = gd.EB_DoorOpen
			} else {
				doorOpenDuration.Reset(timeUntilDoorCloses)
				elevio.SetDoorOpenLamp(false)
			}

		case <-doorOpenDuration.C: //rename?
			if obstruction {
				break
			}

			fmt.Println("The door is closed")

			switch elev.State.Behaviour {
			case gd.EB_DoorOpen:
				DirBehaviourPair := elev.chooseDirection()
				elev.State.Behaviour = DirBehaviourPair.Behaviour
				elev.State.TravelDirection = DirBehaviourPair.Dir
				motorDir := elevio.MotorDirection(elev.State.TravelDirection)

				switch elev.State.Behaviour {
				case gd.EB_DoorOpen:
					doorOpening <- true
					elev.clearOrdersAtCurrentFloor() //update this function to turn off the lights as well
					elev.setButtonLights()           //Doublecheck if this is correct
				case gd.EB_Moving:
					timeoutTimer.Reset(timeUntilTimeout)
				case gd.EB_Idle:
					elevio.SetDoorOpenLamp(false)
					elevio.SetMotorDirection(motorDir) //differ from the handed out C code, different argument, here we know that motordir = MD_Stop (TravelStop) from looking at chooseDirection()
					timeoutTimer.Stop()
				}
				//timeoutTimer.Reset(5 * time.Second)
			}

		case lightMatrix := <-ch.LocalLights2D:
			elev.lights = lightMatrix
			elev.setButtonLights()

		case orderMatrix := <-ch.LocalOrder2D: //navn endres
			//fmt.Println("Button pressed at floor: ", btn.Floor, " Button type: ", btn.Button)
			elev.orders = orderMatrix
			fmt.Println("Orders: ", elev.orders)

			switch elev.State.Behaviour {
			case gd.EB_DoorOpen:
				if elev.orderAtCurrentFloor() {
					elev.clearOrdersAtCurrentFloor()
					doorOpenDuration.Reset(3 * time.Second)
				} else {
					//elev.orders = orderMatrix
				}

			case gd.EB_Idle:
				if elev.orderAtCurrentFloor() {
					elev.State.Behaviour = gd.EB_DoorOpen
					doorOpening <- true
				} else {
					elev.orders = orderMatrix
					DirBehaviourPair := elev.chooseDirection()
					elev.State.Behaviour = DirBehaviourPair.Behaviour
					elev.State.TravelDirection = DirBehaviourPair.Dir
					motorDir := elevio.MotorDirection(elev.State.TravelDirection)

					switch elev.State.Behaviour {
					case gd.EB_DoorOpen:
						doorOpening <- true
						elev.clearOrdersAtCurrentFloor()
						elev.setButtonLights()
					case gd.EB_Moving:
						elevio.SetMotorDirection(motorDir)
						timeoutTimer.Reset(timeUntilTimeout)
					case gd.EB_Idle:
						timeoutTimer.Stop()
					}
				}
				DirBehaviourPair := elev.chooseDirection()
				elev.State.Behaviour = DirBehaviourPair.Behaviour
				elev.State.TravelDirection = DirBehaviourPair.Dir

			case gd.EB_Moving:

				// switch elev.State.Behaviour {
				// case gd.EB_DoorOpen:
				// 	if clearOrdersImmediately(elev, btn.Floor, btn.Button) {
				// 		doorOpenDuration.Reset(3 * time.Second)
				// 	} else {
				// 		elev.orders[btn.Floor][btn.Button] = true
				// 		elev.lights[btn.Floor][btn.Button] = true
				// 		//elev.setButtonLights()
				// 	}
				// 	timeoutTimer.Stop()
				// case gd.EB_Moving:
				// 	elev.orders[btn.Floor][btn.Button] = true
				// 	elev.lights[btn.Floor][btn.Button] = true

				// case gd.EB_Idle:
				// 	elev.orders[btn.Floor][btn.Button] = true
				// 	elev.lights[btn.Floor][btn.Button] = true

				// 	DirBehaviourPair := elev.chooseDirection()
				// 	elev.State.Behaviour = DirBehaviourPair.Behaviour
				// 	elev.State.TravelDirection = DirBehaviourPair.Dir
				// 	motordir := elevio.MotorDirection(elev.State.TravelDirection)

				// 	switch elev.State.Behaviour {
				// 	case gd.EB_DoorOpen:
				// 		doorOpened <- true
				// 		elev.clearOrdersAtCurrentFloor()
				// 	case gd.EB_Moving:
				// 		elevio.SetMotorDirection(motordir)
				// 		timeoutTimer.Reset(timeUntilTimeout)
				// 	case gd.EB_Idle:
				// 		timeoutTimer.Stop()
				// 	}
			}

		case <-sendElevatorState.C:
			ch.ToMaster <- elev.State
			sendElevatorState.Reset(1 * time.Second)

		case <-timeoutTimer.C:
			fmt.Println("Elevator has timed out")
			elev.State.Behaviour = gd.EB_Moving
			elevio.SetMotorDirection(elevio.MD_Down)
			elev.State.TravelDirection = gd.TravelDown
			//timeoutTimer.Reset(timeUntilTimeout)

			/* default:
			elev.State.Floor = -1
			elevio.SetMotorDirection(elevio.MD_Down)
			elev.State.TravelDirection = gd.TravelDown
			elev.State.Behaviour = gd.EB_Moving
			timeoutTimer.Reset(timeUntilTimeout) */
		}
	}
}

/* func (elev Elevator) fsm_onInitBetweenFloors() { //default case within startElevatorController
	elevio.SetMotorDirection(elevio.MD_Down)
	elev.State.TravelDirection = gd.TravelDown
	elev.State.Behaviour = gd.EB_Moving
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

func (elev Elevator) clearOrdersAtCurrentFloor() {
	switch elev.State.Config.ClearRequestVariant {
	case gd.ClearRequests_All:
		for btn := 0; btn < gd.N_BUTTONS; btn++ {
			elev.orders[elev.State.Floor][btn] = false
			elev.lights[elev.State.Floor][btn] = false
		}

	case gd.ClearRequests_InMotorDir:
		elev.orders[elev.State.Floor][gd.BT_Cab] = false
		elev.lights[elev.State.Floor][gd.BT_Cab] = false

		switch elev.State.TravelDirection {
		case gd.TravelUp:
			if !elev.orderAbove() && !elev.orders[elev.State.Floor][gd.BT_HallUp] {
				elev.orders[elev.State.Floor][gd.BT_HallDown] = false
				elev.lights[elev.State.Floor][gd.BT_HallDown] = false
			}
			elev.orders[elev.State.Floor][gd.BT_HallUp] = false
			elev.lights[elev.State.Floor][gd.BT_HallUp] = false
		case gd.TravelDown:
			if !elev.orderBelow() && !elev.orders[elev.State.Floor][gd.BT_HallDown] {
				elev.orders[elev.State.Floor][gd.BT_HallUp] = false
				elev.lights[elev.State.Floor][gd.BT_HallUp] = false
			}
			elev.orders[elev.State.Floor][gd.BT_HallDown] = false
			elev.lights[elev.State.Floor][gd.BT_HallDown] = false
		case gd.TravelStop:
		default:
			elev.orders[elev.State.Floor][gd.BT_HallUp] = false
			elev.orders[elev.State.Floor][gd.BT_HallDown] = false
			elev.lights[elev.State.Floor][gd.BT_HallUp] = false
			elev.lights[elev.State.Floor][gd.BT_HallDown] = false
		}
	}
}

func clearOrdersImmediately(elev Elevator, btn_floor int, btn_type gd.ButtonType) bool {
	switch elev.State.Config.ClearRequestVariant {
	case gd.ClearRequests_All:
		return elev.State.Floor == btn_floor
	case gd.ClearRequests_InMotorDir:
		return elev.State.Floor == btn_floor && ((elev.State.TravelDirection == gd.TravelUp && btn_type == gd.BT_HallUp) ||
			(elev.State.TravelDirection == gd.TravelDown && btn_type == gd.BT_HallDown) ||
			(elev.State.TravelDirection == gd.TravelStop) || (btn_type == gd.BT_Cab))
	default:
		return false
	}
}

// setButtonLamp sets the light of the button at the given floor to the given value, might have to make some adjustments to this function
func (elev Elevator) setButtonLights() {
	for floor := 0; floor < gd.N_FLOORS; floor++ {
		for btn := 0; btn < gd.N_BUTTONS; btn++ {
			elevio.SetButtonLamp(gd.ButtonType(btn), floor, elev.lights[floor][btn])
		}
	}
}

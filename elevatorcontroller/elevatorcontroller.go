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

func StartElevatorController(
	ElevatorID string,
	numFloors int,

	currentFloorCh chan int,
	obstructionEvent chan bool,
	stopCh chan bool,
	completedOrder_toOrderHandler chan<- gd.ButtonEvent,
	localLights2D <-chan gd.Orders2D,
	localOrder2D <-chan gd.Orders2D,
	toMaster chan<- gd.ElevatorState,
) {

	fmt.Println("ElevatorController started")

	elev := InitiateElevator(ElevatorID, "localhost:15657", numFloors)

	go elevio.PollFloorSensor(currentFloorCh)
	go elevio.PollObstructionSwitch(obstructionEvent)
	go elevio.PollStopButton(stopCh)

	doorOpening := make(chan bool, 100)

	doorOpenDuration := time.NewTimer(timeUntilDoorCloses)
	timeoutTimer := time.NewTimer(timeUntilTimeout)
	sendElevatorState := time.NewTimer(1 * time.Second)

	doorOpenDuration.Stop()
	elevio.SetDoorOpenLamp(false)
	obstruction := false

	for {
		select {
		case ArrivingAtFloor := <-currentFloorCh:
			elev.State.Floor = ArrivingAtFloor
			fmt.Println("Arriving at floor: ", elev.State.Floor)
			elevio.SetFloorIndicator(elev.State.Floor)

			switch elev.State.Behaviour {
			case gd.EB_Moving:
				if elev.shouldStop() {
					elevio.SetMotorDirection(elevio.MD_Stop)
					ordersCleared := elev.clearOrdersAtCurrentFloor()
					elev.setButtonLights()
					for _, order := range ordersCleared {
						completedOrder_toOrderHandler <- order
					}
					doorOpening <- true
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
			toMaster <- elev.State

		case <-doorOpening:
			fmt.Println("The door is open")
			elevio.SetDoorOpenLamp(true)
			doorOpenDuration.Reset(timeUntilDoorCloses)
			elevio.SetDoorOpenLamp(false)
			timeoutTimer.Stop()

		case doorObstructed := <-obstructionEvent:
			obstruction = doorObstructed
			if obstruction {
				fmt.Println("Obstruction detected")
				elevio.SetDoorOpenLamp(true)
				elev.State.Behaviour = gd.EB_DoorOpen
			} else {
				doorOpenDuration.Reset(timeUntilDoorCloses)
				elevio.SetDoorOpenLamp(false)
			}

		case <-doorOpenDuration.C:
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
					ordersCleared := elev.clearOrdersAtCurrentFloor()
					elev.setButtonLights()
					for _, order := range ordersCleared {
						completedOrder_toOrderHandler <- order
					}
				case gd.EB_Moving:
					timeoutTimer.Reset(timeUntilTimeout)
				case gd.EB_Idle:
					elevio.SetDoorOpenLamp(false)
					elevio.SetMotorDirection(motorDir)
					timeoutTimer.Stop()
				}
			}

		case lightMatrix := <-localLights2D:
			elev.lights = lightMatrix
			elev.setButtonLights()

		case orderMatrix := <-localOrder2D:
			elev.orders = orderMatrix
			fmt.Println("Orders: ", elev.orders)

			switch elev.State.Behaviour {
			case gd.EB_DoorOpen:
				if elev.shouldClearOrderImmediately(elev.orders) {
					ordersCleared := elev.returnClearedOrders(elev.orders)
					elev.setButtonLights()
					for _, order := range ordersCleared {
						completedOrder_toOrderHandler <- order
					}
					doorOpening <- true
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
						ordersCleared := elev.clearOrdersAtCurrentFloor()
						elev.setButtonLights()
						for _, order := range ordersCleared {
							completedOrder_toOrderHandler <- order
						}

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
			}

		case <-stopCh:
			elevio.SetMotorDirection(elevio.MD_Stop)

		case <-sendElevatorState.C:
			toMaster <- elev.State
			sendElevatorState.Reset(1 * time.Second)

		case <-timeoutTimer.C:
			fmt.Println("Elevator has timed out")
			elev.State.Behaviour = gd.EB_Moving
			elevio.SetMotorDirection(elevio.MD_Down)
			elev.State.TravelDirection = gd.TravelDown

		}
	}
}

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
			return DirBehaviourPair{gd.TravelStop, gd.EB_Idle}
		}
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

func (elev Elevator) clearOrdersAtCurrentFloor() []gd.ButtonEvent {
	var clearedOrders []gd.ButtonEvent
	switch elev.State.Config.ClearRequestVariant {
	case gd.ClearRequests_All:
		for btn := 0; btn < gd.N_BUTTONS; btn++ {
			elev.orders[elev.State.Floor][btn] = false
			elev.lights[elev.State.Floor][btn] = false
			clearedOrders = append(clearedOrders, gd.ButtonEvent{Floor: elev.State.Floor, Button: gd.ButtonType(btn)})
		}
	case gd.ClearRequests_InMotorDir:
		elev.orders[elev.State.Floor][gd.BT_Cab] = false
		elev.lights[elev.State.Floor][gd.BT_Cab] = false
		clearedOrders = append(clearedOrders, gd.ButtonEvent{Floor: elev.State.Floor, Button: gd.BT_Cab})

		switch elev.State.TravelDirection {
		case gd.TravelUp:
			if !elev.orderAbove() && !elev.orders[elev.State.Floor][gd.BT_HallUp] {
				elev.orders[elev.State.Floor][gd.BT_HallDown] = false
				elev.lights[elev.State.Floor][gd.BT_HallDown] = false
				clearedOrders = append(clearedOrders, gd.ButtonEvent{Floor: elev.State.Floor, Button: gd.ButtonType(gd.BT_HallDown)})
			}
			elev.orders[elev.State.Floor][gd.BT_HallUp] = false
			elev.lights[elev.State.Floor][gd.BT_HallUp] = false
			clearedOrders = append(clearedOrders, gd.ButtonEvent{Floor: elev.State.Floor, Button: gd.ButtonType(gd.BT_HallUp)})
		case gd.TravelDown:
			if !elev.orderBelow() && !elev.orders[elev.State.Floor][gd.BT_HallDown] {
				elev.orders[elev.State.Floor][gd.BT_HallUp] = false
				elev.lights[elev.State.Floor][gd.BT_HallUp] = false
				clearedOrders = append(clearedOrders, gd.ButtonEvent{Floor: elev.State.Floor, Button: gd.ButtonType(gd.BT_HallUp)})
			}
			elev.orders[elev.State.Floor][gd.BT_HallDown] = false
			elev.lights[elev.State.Floor][gd.BT_HallDown] = false
			clearedOrders = append(clearedOrders, gd.ButtonEvent{Floor: elev.State.Floor, Button: gd.ButtonType(gd.BT_HallDown)})
		case gd.TravelStop:
		default:
			elev.orders[elev.State.Floor][gd.BT_HallUp] = false
			elev.orders[elev.State.Floor][gd.BT_HallDown] = false
			elev.lights[elev.State.Floor][gd.BT_HallUp] = false
			elev.lights[elev.State.Floor][gd.BT_HallDown] = false
			clearedOrders = append(clearedOrders, gd.ButtonEvent{Floor: elev.State.Floor, Button: gd.ButtonType(gd.BT_HallUp)})
		}

	}
	return clearedOrders
}

func (elev Elevator) shouldClearOrderImmediately(Orders gd.Orders2D) bool {
	buttonEvents := OrderMatrixToButtonEvent(Orders)
	for _, event := range buttonEvents {
		floor := event.Floor
		btn_type := event.Button

		if elev.State.Floor == floor {
			switch elev.State.Config.ClearRequestVariant {
			case gd.ClearRequests_All:
				return true
			case gd.ClearRequests_InMotorDir:
				return (elev.State.TravelDirection == gd.TravelUp && btn_type == gd.BT_HallUp) ||
					(elev.State.TravelDirection == gd.TravelDown && btn_type == gd.BT_HallDown) ||
					(elev.State.TravelDirection == gd.TravelStop) || (btn_type == gd.BT_Cab)

			}
		}
	}
	return false
}

func (elev Elevator) returnClearedOrders(Orders gd.Orders2D) []gd.ButtonEvent {
	buttonEvents := OrderMatrixToButtonEvent(Orders)
	var clearedOrders []gd.ButtonEvent
	for _, event := range buttonEvents {
		floor := event.Floor
		btn_type := event.Button

		if elev.State.Floor == floor {
			switch elev.State.Config.ClearRequestVariant {
			case gd.ClearRequests_All:
				clearedOrders = append(clearedOrders, gd.ButtonEvent{Floor: floor, Button: gd.ButtonType(btn_type)})
			case gd.ClearRequests_InMotorDir:
				if (elev.State.TravelDirection == gd.TravelUp && btn_type == gd.BT_HallUp) ||
					(elev.State.TravelDirection == gd.TravelDown && btn_type == gd.BT_HallDown) ||
					(elev.State.TravelDirection == gd.TravelStop) || (btn_type == gd.BT_Cab) {
					clearedOrders = append(clearedOrders, gd.ButtonEvent{Floor: floor, Button: gd.ButtonType(btn_type)})
				}
			}

		}
	}
	return clearedOrders
}

func (elev Elevator) setButtonLights() {
	for floor := 0; floor < gd.N_FLOORS; floor++ {
		for btn := 0; btn < gd.N_BUTTONS; btn++ {
			elevio.SetButtonLamp(gd.ButtonType(btn), floor, elev.lights[floor][btn])
		}
	}
}

func OrderMatrixToButtonEvent(Orders gd.Orders2D) []gd.ButtonEvent {
	var NewButtonEvent []gd.ButtonEvent
	for floor := 0; floor < len(Orders); floor++ {
		for button := 0; button < len(Orders[0]); button++ {
			if Orders[floor][button] {
				buttonType := gd.ButtonType(button)
				NewButtonEvent = append(NewButtonEvent, gd.ButtonEvent{Floor: floor, Button: buttonType})
			}
		}
	}

	return NewButtonEvent
}

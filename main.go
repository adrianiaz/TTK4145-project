package main

import (
	"fmt"

	"github.com/adrianiaz/TTK4145-project/elevio"
	gd "github.com/adrianiaz/TTK4145-project/globaldefinitions"
)

func main() {

	numFloors := 4 //floors

	elevio.Init("localhost:15657", numFloors)

	var d elevio.MotorDirection = elevio.MD_Up
	//elevio.SetMotorDirection(d)
	currentFloor := elevio.GetFloor()

	drv_buttons := make(chan gd.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			if a.Button == gd.BT_Cab {
				switch {
				case currentFloor == 0:
					if a.Floor == currentFloor {
						elevio.SetMotorDirection(elevio.MD_Stop)
					} else {
						elevio.SetMotorDirection(elevio.MD_Up)
					}

				case currentFloor == 1:
					if a.Floor == currentFloor {
						elevio.SetMotorDirection(elevio.MD_Stop)
					} else if a.Floor > currentFloor {
						elevio.SetMotorDirection(elevio.MD_Up)
					} else {
						elevio.SetMotorDirection(elevio.MD_Down)
					}

				case currentFloor == 2:
					if a.Floor == currentFloor {
						elevio.SetMotorDirection(elevio.MD_Stop)
					} else if a.Floor > currentFloor {
						elevio.SetMotorDirection(elevio.MD_Up)
					} else {
						elevio.SetMotorDirection(elevio.MD_Down)
					}

				case currentFloor == 3:
					if a.Floor == currentFloor {
						elevio.SetMotorDirection(elevio.MD_Stop)

					} else {
						elevio.SetMotorDirection(elevio.MD_Down)
					}
				}
			}
			elevio.SetButtonLamp(a.Button, a.Floor, true)

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			if a == numFloors-1 {
				d = elevio.MD_Down
			} else if a == 0 {
				d = elevio.MD_Up
			}
			switch {
			case elevio.GetButton(gd.BT_Cab, currentFloor):
				d = elevio.MD_Stop
			}
			elevio.SetMotorDirection(d)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(d)
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := gd.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
		}
	}
}

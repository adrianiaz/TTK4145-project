package main

import (
	"github.com/adrianiaz/TTK4145-project/elevio"
)

func main() {
	numFloors := 4

	elevio.Init("localhost:15657", numFloors)

	//channel for changes
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	//input processes to state machine
	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	//FSM
	for {
		select {
		case c := <-drv_buttons:

		}
	}

}

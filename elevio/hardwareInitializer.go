package elevio

import (
	gd "github.com/adrianiaz/TTK4145-project/globaldefinitions"
)

func HardwareInitilizer(
	hw_floor chan<- int,
	hw_obstruction chan<- bool,
	hw_stop chan<- bool,
	hw_button chan<- gd.ButtonEvent,
	numFloors int,
) {
	Init("localhost:15657", numFloors)

	go PollButtons(hw_button)
	go PollFloorSensor(hw_floor)
	go PollObstructionSwitch(hw_obstruction)
	go PollStopButton(hw_stop)

}

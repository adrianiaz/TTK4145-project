package globaltypes

const (
	N_FLOORS  = 4
	N_BUTTONS = 3
)

type ElevatorBehaviour int

const (
	EB_Idle ElevatorBehaviour = iota //this should probably be moved to another file.
	EB_DoorOpen
	EB_Moving
)

const (
	travellingUp = iota
	travellingDown
)

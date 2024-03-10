package globaltypes

const (
	N_FLOORS  = 4
	N_BUTTONS = 3
)

type ElevatorBehaviour int

const (
	EB_Idle ElevatorBehaviour = iota
	EB_DoorOpen
	EB_Moving
)

type TravelDir int

type elevatorState struct {
	Floor           int
	Behaviour       ElevatorBehaviour
	TravelDirection TravelDir
	ElevatorID      string
	Requests        [N_FLOORS][N_BUTTONS]bool
}

const (
	travellingUp TravelDir = iota
	travellingDown
)

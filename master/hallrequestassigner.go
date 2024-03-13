package master

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"

	gd "github.com/adrianiaz/TTK4145-project/globaldefinitions"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

type HRAElevState struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [][2]bool               `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

// takes in the elevator state returns the new hallorders for the elevators
func HRA(allHallRequests [gd.N_FLOORS][2]bool, states gd.AllElevatorStates) {

	hraExecutable := ""
	switch runtime.GOOS {
	case "linux":
		hraExecutable = "hall_request_assigner"
	case "windows":
		hraExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	//instantiate a map of HRAElevState structs
	stateMap := make(map[string]HRAElevState)
	for elevatorID, state := range states {
		stateMap[elevatorID] = HRAElevState{
			Behavior:    string(state.Behaviour),
			Floor:       state.Floor,
			Direction:   string(state.TravelDirection),
			CabRequests: state.Requests[state.Floor],
		}
	}

	input := HRAInput{
		HallRequests: allHallRequests,
		States: map[string]HRAElevState{
			"one": HRAElevState{
				Behavior:    "moving",
				Floor:       2,
				Direction:   "up",
				CabRequests: []bool{false, false, false, true},
			},
			"two": HRAElevState{
				Behavior:    "idle",
				Floor:       0,
				Direction:   "stop",
				CabRequests: []bool{false, false, false, false},
			},
		},
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
		return
	}

	ret, err := exec.Command("../hall_request_assigner/"+hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
		return
	}

	output := new(map[string][][2]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
		return
	}

	fmt.Printf("output: \n")
	for k, v := range *output {
		fmt.Printf("%6v :  %+v\n", k, v)
	}
}

func orderSplitter() {

}

func extractCabRequests(allElevStates gd.AllElevatorStates) map[string][]bool {
	CabRequestMap := make(map[string][]bool) // Initialize CabRequestMap as an empty map
	for _, state := range allElevStates {
		CabRequestMap[state.ElevatorID] = []bool{} // Initialize each elevator's floor requests as an empty array
		for floor := range state.Requests {
			CabRequestMap[state.ElevatorID] = append(CabRequestMap[state.ElevatorID], state.Requests[floor][gd.BT_Cab])
		}
	}
	return CabRequestMap
}

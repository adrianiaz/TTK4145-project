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
func HRA(allHallRequests [][2]bool, states gd.AllElevatorStates) {

	hraExecutable := ""
	switch runtime.GOOS {
	case "linux":
		hraExecutable = "hall_request_assigner"
	case "windows":
		hraExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	//extract cab requests from allElevStates
	cabRequestMap := extractCabRequests(states)

	//instantiate a map of HRAElevState structs
	stateMap := make(map[string]HRAElevState)
	for elevatorID, state := range states {
		stateMap[elevatorID] = HRAElevState{
			Behavior:    string(state.Behaviour),
			Floor:       state.Floor,
			Direction:   string(state.TravelDirection),
			CabRequests: cabRequestMap[elevatorID],
		}
	}

	input := HRAInput{
		HallRequests: allHallRequests,
		States:       stateMap,
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

func extractCabRequests(allElevStates gd.AllElevatorStates) map[string][]bool {
	cabRequestMap := make(map[string][]bool) // Initialize CabRequestMap as an empty map
	for _, state := range allElevStates {
		cabRequestMap[state.ElevatorID] = []bool{} // Initialize each elevator's floor requests as an empty array
		for floor := range state.Requests {
			cabRequestMap[state.ElevatorID] = append(cabRequestMap[state.ElevatorID], state.Requests[floor][gd.BT_Cab])
		}
	}
	return cabRequestMap
}

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
func extractOptimalHallRequests(ledger gd.Ledger, newHallRequest gd.Order) map[string][][2]bool {

	hraExecutable := ""
	switch runtime.GOOS {
	case "linux":
		hraExecutable = "hall_request_assigner"
	case "windows":
		hraExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	states := ledger.ElevatorStates
	allHallRequests := findAllHallRequests(ledger.ActiveOrders, newHallRequest) // Need to include new order in this
	cabRequestMap := extractCabRequests(states)

	optimalHallRequests := new(map[string][][2]bool)

	//instantiate a map of HRAElevState structs
	stateMap := make(map[string]HRAElevState)
	for elevatorID, state := range states {
		stateMap[elevatorID] = HRAElevState{
			Behavior:    fmt.Sprint(state.Behaviour),
			Floor:       state.Floor,
			Direction:   fmt.Sprint(state.TravelDirection),
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
		return *optimalHallRequests
	}

	ret, err := exec.Command("../hall_request_assigner/"+hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
		return *optimalHallRequests
	}

	err = json.Unmarshal(ret, &optimalHallRequests)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
		return *optimalHallRequests
	}

	fmt.Printf("output: \n")
	for k, v := range *optimalHallRequests {
		fmt.Printf("%6v :  %+v\n", k, v)
	}
	return *optimalHallRequests
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

func findAllHallRequests(allorders gd.AllOrders, newHallRequest gd.Order) [][2]bool {
	allHallRequests := make([][2]bool, gd.N_FLOORS)
	for _, elevator := range allorders {
		for floor := 0; floor < gd.N_FLOORS; floor++ {
			for btnType := 0; btnType < 2; btnType++ {
				allHallRequests[floor][btnType] = elevator[floor][btnType]
			}
		}
	}
	allHallRequests[newHallRequest.Floor][int(newHallRequest.BtnType)] = true
	return allHallRequests
}

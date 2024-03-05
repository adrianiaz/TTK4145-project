package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"
	"time"

	"github.com/adrianiaz/TTK4145-project/elevio"
)

func main() {

	numFloors := 4

	elevio.Init("localhost:15657", numFloors) //connection for elevatorprogram

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
	// Primary-backup module. Running this function opens a termnal listening for a master node, if not found, makes itself the master.
	// The backup should recieve the ledgerlist where its orders are specified. The master sends this message to all the slaves including the backup.

	connectionVar := "127.0.0.1:27106" //localhost, should be changed to local network

	addr, err := net.ResolveUDPAddr("udp", connectionVar)
	if err != nil {
		log.Fatal(err)
	}

	listening, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	var message int //This should be the ledger

	//Backup mode loop, here slave code should go
	for {
		buf := make([]byte, 1024)
		err = listening.SetReadDeadline(time.Now().Add(3 * time.Second)) // give time for master to send message
		if err != nil {
			fmt.Println("Read time-out")
			break
		}
		n, err := listening.Read(buf)
		if err != nil {
			log.Fatal(err)
			break
		}
		message, _ = strconv.Atoi(string(buf[:n]))
	}

	//close this socket that we are reading from as it is old backup
	listening.Close()

	//This code executes only if there is no master on the network. Thus the program starts again and the current node takes over as master
	exec.Command("gnome-terminal", "--", "go", "run", "primback.go").Run() //run primback.go might not work, might have to make in run the main.go-file
	time.Sleep(1 * time.Second)                                            //pause to give time for connection to close before

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	//Primary loop, here master code should go
	for {
		message++
		_, err := conn.Write([]byte(strconv.Itoa(message)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(message)
		time.Sleep(1 * time.Second)
	}
}

//code to be put into master loop

// 	for {
// 		select {
// 		case c := <-drv_buttons:

// 		}
// 	}

// }

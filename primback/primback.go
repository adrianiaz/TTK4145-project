package primback

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"
	"time"
)

// Primary-backup module. Running this function opens a termnal listening for a master node, if not found, makes itself the master.
// The backup should recieve the ledgerlist where its orders are specified. The master sends this message to all the slaves including the backup.
func primback() {
	connectionVar := "127.0.0.1:27106"

	addr, err := net.ResolveUDPAddr("udp", connectionVar)
	if err != nil {
		log.Fatal(err)
	}

	listening, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	var message int //This should be the ledger

	//Backup mode loop
	for {
		buf := make([]byte, 1024)
		err = listening.SetReadDeadline(time.Now().Add(3 * time.Second))
		n, err := listening.Read(buf)
		if err != nil {
			break
		}
		message, _ = strconv.Atoi(string(buf[:n]))
	}

	//close this socket that we are reading from as it is old backup
	listening.Close()

	//This code executes only if there is no master on the network. Thus the program starts again and the current node takes over as master
	exec.Command("gnome-terminal", "--", "go", "run", "primback.go").Run() //run primback.go might not work, might have to make in run the main.go-file
	time.Sleep(1 * time.Second)

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	//Primary loop
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
